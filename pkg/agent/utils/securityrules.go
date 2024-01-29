// Copyright 2019 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"fmt"
	"strings"

	"yunion.io/x/pkg/util/regutils"
	"yunion.io/x/pkg/util/secrules"
)

type SecurityRule struct {
	r          *secrules.SecurityRule
	ovsMatches []string
}

func NewSecurityRule(s string) (*SecurityRule, error) {
	s = strings.TrimSpace(s)
	r, err := secrules.ParseSecurityRule(s)
	if err != nil {
		return nil, err
	}
	return &SecurityRule{r: r}, nil
}

func (sr *SecurityRule) Direction() secrules.TSecurityRuleDirection {
	return sr.r.Direction
}

func (sr *SecurityRule) OvsMatches() []string {
	if sr.ovsMatches != nil {
		return sr.ovsMatches
	}

	var nwProto string
	var nwField string
	var v6Field string
	var tpField string
	var protoMatch string
	var nwMatch string
	var tpMatch []string

	r := sr.r
	switch r.Direction {
	case secrules.DIR_IN:
		nwField = "ip,nw_src="
		v6Field = "ipv6,ipv6_src="
		tpField = "tp_dst="
	case secrules.DIR_OUT:
		nwField = "ip,nw_dst="
		v6Field = "ipv6,ipv6_dst="
		tpField = "tp_dst="
	}
	if r.IPNet != nil {
		netStr := r.IPNet.String()
		if regutils.MatchCIDR6(netStr) {
			// ipv6
			nwProto = "ipv6"
			if ones, bits := r.IPNet.Mask.Size(); ones == 128 && bits == 128 {
				nwMatch = v6Field + r.IPNet.IP.String()
			} else {
				nwMatch = v6Field + netStr
			}
		} else {
			// ipv4
			nwProto = "ip"
			if ones, bits := r.IPNet.Mask.Size(); ones == 32 && bits == 32 {
				nwMatch = nwField + r.IPNet.IP.String()
			} else {
				nwMatch = nwField + netStr
			}
		}
	}
	switch r.Protocol {
	case secrules.PROTO_ANY:
		if len(nwMatch) == 0 {
			tpMatch = append(tpMatch, "ipv6", "ip")
		}
		// protoMatch = "ip"
	case secrules.PROTO_TCP, secrules.PROTO_UDP:
		if len(r.Ports) > 0 {
			for _, p := range r.Ports {
				tpMatch = append(tpMatch, tpField+fmt.Sprintf("%d", p))
			}
		}
		if r.PortStart > 0 && r.PortStart <= r.PortEnd {
			ms := PortRangeToMasks(uint16(r.PortStart), uint16(r.PortEnd))
			for _, m := range ms {
				// NOTE both start and end should never be zero, the
				// check is here just in case
				if m[1] == 0xffff {
					tpMatch = append(tpMatch, fmt.Sprintf("%s%d", tpField, m[0]))
				} else if m[1] == 0 {
					break
				} else {
					var vs string
					if m[0] == 0 {
						vs = "0"
					} else {
						vs = fmt.Sprintf("0x%x", m[0])
					}
					tpMatch = append(tpMatch, fmt.Sprintf("%s%s/0x%x", tpField, vs, m[1]))
				}
			}
		}
		protoMatch = r.Protocol
	case secrules.PROTO_ICMP:
		if nwProto == "ipv6" {
			protoMatch = "icmp6"
		} else if nwProto == "ip" {
			protoMatch = "icmp"
		} else {
			tpMatch = append(tpMatch, "icmp", "icmp6")
		}
	default:
		protoMatch = r.Protocol
	}

	var m string
	if len(nwMatch) > 0 {
		m = nwMatch
	}
	if len(protoMatch) > 0 {
		if len(m) > 0 {
			m += ","
		}
		m += protoMatch
	}
	if len(tpMatch) == 0 {
		sr.ovsMatches = []string{m}
	} else {
		ms := []string{}
		for _, tpm := range tpMatch {
			if len(m) > 0 {
				tpm = m + "," + tpm
			}
			ms = append(ms, tpm)
		}
		sr.ovsMatches = ms
	}
	return sr.ovsMatches
}

func (sr *SecurityRule) OvsActionAllow() bool {
	return sr.r.Action == secrules.SecurityRuleAllow
}

func (sr *SecurityRule) IsWildMatch() bool {
	return sr.r.IsWildMatch()
}

// TODO squash neighbouring rules of the same direction
type SecurityRules struct {
	inRules       []*SecurityRule
	outRules      []*SecurityRule
	inOvsMatches  []string
	outOvsMatches []string
	outAllowAny   bool
}

func (sr *SecurityRules) rulesString(srs []*SecurityRule) string {
	v := []string{}
	for _, r := range srs {
		v = append(v, r.r.String())
	}
	return strings.Join(v, "; ")
}

func (sr *SecurityRules) InRulesString() string {
	return sr.rulesString(sr.inRules)
}

func (sr *SecurityRules) OutRulesString() string {
	return sr.rulesString(sr.outRules)
}

func NewSecurityRules(s string) (*SecurityRules, error) {
	inRules := []*SecurityRule{}
	outRules := []*SecurityRule{}
	srs := strings.Split(s, ";")
	for _, sr := range srs {
		sr = strings.TrimSpace(sr)
		if len(sr) == 0 {
			continue
		}
		r, err := NewSecurityRule(sr)
		if err != nil {
			// TODO err wrapper
			return nil, err
		}
		switch r.Direction() {
		case secrules.DIR_IN:
			inRules = append(inRules, r)
		case secrules.DIR_OUT:
			outRules = append(outRules, r)
		}
	}
	// In the case where no secgroup was assigned, default security_rules
	// "in:allow_any; out:allow_any" will be used by the caller
	if l := len(inRules); l == 0 || (l > 0 && !inRules[l-1].IsWildMatch()) {
		r, _ := NewSecurityRule("in:deny any")
		inRules = append(inRules, r)
	}
	if l := len(outRules); l == 0 || (l > 0 && !outRules[l-1].IsWildMatch()) {
		r, _ := NewSecurityRule("out:allow any")
		outRules = append(outRules, r)
	}
	return &SecurityRules{
		inRules:  inRules,
		outRules: outRules,
	}, nil
}

func PortRangeToMasks(s, e uint16) [][2]uint16 {
	r := [][2]uint16{}
	if s == e {
		r = append(r, [2]uint16{s, ^uint16(0)})
		return r
	}
	sp, ep := uint32(s), uint32(e)
	ep = ep + 1
	for sp < ep {
		b := uint32(1)
		for (sp+b) <= ep && (sp&(b-1)) == 0 {
			b <<= 1
		}
		b >>= 1
		r = append(r, [2]uint16{uint16(sp), uint16(^(b - 1))})
		sp = sp + b
	}
	return r
}
