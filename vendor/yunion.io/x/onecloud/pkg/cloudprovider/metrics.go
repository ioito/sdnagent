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

package cloudprovider

import (
	"strings"
	"time"
)

type TResourceType string
type TMetricType string

func (self TMetricType) Name() string {
	if !strings.Contains(string(self), ".") {
		return string(self)
	}
	return string(self)[0:strings.Index(string(self), ".")]
}

func (self TMetricType) Key() string {
	return func(key string) string {
		if len(key) == 0 {
			return ""
		}
		first, last := 0, len(key)
		if strings.Contains(key, ",") {
			last = strings.Index(key, ",")
		}
		if strings.Contains(key, ".") {
			first = strings.LastIndex(key, ".") + 1
		}
		return key[first:last]
	}(string(self))
}

const (
	METRIC_RESOURCE_TYPE_RDS    TResourceType = "RDS"
	METRIC_RESOURCE_TYPE_SERVER TResourceType = "SERVER"
	METRIC_RESOURCE_TYPE_REDIS  TResourceType = "REDIS"
	METRIC_RESOURCE_TYPE_LB     TResourceType = "LB"
	METRIC_RESOURCE_TYPE_BUCKET TResourceType = "BUCKET"
)

const (
	// RDS监控指标
	RDS_METRIC_TYPE_CPU_USAGE  TMetricType = "rds_cpu.usage_active"
	RDS_METRIC_TYPE_MEM_USAGE  TMetricType = "rds_mem.used_percent"
	RDS_METRIC_TYPE_NET_BPS_RX TMetricType = "rds_netio.bps_recv"
	RDS_METRIC_TYPE_NET_BPS_TX TMetricType = "rds_netio.bps_send"

	RDS_METRIC_TYPE_DISK_USAGE      TMetricType = "rds_disk.used_percent"
	RDS_METRIC_TYPE_DISK_READ_BPS   TMetricType = "rds_diskio.read_bps"
	RDS_METRIC_TYPE_DISK_WRITE_BPS  TMetricType = "rds_diskio.write_bps"
	RDS_METRIC_TYPE_DISK_IO_PERCENT TMetricType = "rds_diskio.used_percent"

	RDS_METRIC_TYPE_CONN_COUNT  TMetricType = "rds_conn.used_count"
	RDS_METRIC_TYPE_CONN_ACTIVE TMetricType = "rds_conn.active_count"
	RDS_METRIC_TYPE_CONN_USAGE  TMetricType = "rds_conn.used_percent"
	RDS_METRIC_TYPE_CONN_FAILED TMetricType = "rds_conn.failed_count"

	RDS_METRIC_TYPE_QPS              TMetricType = "rds_qps.query_qps"
	RDS_METRIC_TYPE_TPS              TMetricType = "rds_tps.trans_qps"
	RDS_METRIC_TYPE_INNODB_READ_BPS  TMetricType = "rds_innodb.read_bps"
	RDS_METRIC_TYPE_INNODB_WRITE_BPS TMetricType = "rds_innodb.write_bps"

	VM_METRIC_TYPE_CPU_USAGE  TMetricType = "vm_cpu.usage_active"
	VM_METRIC_TYPE_MEM_USAGE  TMetricType = "vm_mem.used_percent"
	VM_METRIC_TYPE_DISK_USAGE TMetricType = "vm_disk.used_percent"

	VM_METRIC_TYPE_DISK_IO_READ_BPS   TMetricType = "vm_diskio.read_bps"
	VM_METRIC_TYPE_DISK_IO_WRITE_BPS  TMetricType = "vm_diskio.write_bps"
	VM_METRIC_TYPE_DISK_IO_READ_IOPS  TMetricType = "vm_diskio.read_iops"
	VM_METRIC_TYPE_DISK_IO_WRITE_IOPS TMetricType = "vm_diskio.write_iops"

	VM_METRIC_TYPE_NET_BPS_RX TMetricType = "vm_netio.bps_recv"
	VM_METRIC_TYPE_NET_BPS_TX TMetricType = "vm_netio.bps_send"

	VM_METRIC_TYPE_EIP_BPS_IN  TMetricType = "vm_eipio.bps_in"
	VM_METRIC_TYPE_EIP_BPS_OUT TMetricType = "vm_eipio.bps_out"

	VM_METRIC_TYPE_EIP_PPS_IN  TMetricType = "vm_eipio.pps_in"
	VM_METRIC_TYPE_EIP_PPS_OUT TMetricType = "vm_eipio.pps_out"

	REDIS_METRIC_TYPE_CPU_USAGE      = "dcs_cpu.usage_active"
	REDIS_METRIC_TYPE_MEM_USAGE      = "dcs_mem.used_percent"
	REDIS_METRIC_TYPE_NET_BPS_RX     = "dcs_netio.bps_recv"
	REDIS_METRIC_TYPE_NET_BPS_TX     = "dcs_netio.bps_sent"
	REDIS_METRIC_TYPE_CONN_USAGE     = "dcs_conn.used_percent"
	REDIS_METRIC_TYPE_OPT_SES        = "dcs_instantopt.opt_sec"
	REDIS_METRIC_TYPE_CACHE_KEYS     = "dcs_cachekeys.key_count"
	REDIS_METRIC_TYPE_CACHE_EXP_KEYS = "dcs_cachekeys.key_count,exp=expire"
	REDIS_METRIC_TYPE_DATA_MEM_USAGE = "dcs_datamem.used_byte"
	REDIS_METRIC_TYPE_SERVER_LOAD    = "dcs_cpu.server_load"
	REDIS_METRIC_TYPE_CONN_ERRORS    = "dcs_conn.errors"

	LB_METRIC_TYPE_SNAT_PORT       = "haproxy.used_snat_port"
	LB_METRIC_TYPE_SNAT_CONN_COUNT = "haproxy.snat_conn_count"
	LB_METRIC_TYPE_NET_BPS_RX      = "haproxy.bin"
	LB_METRIC_TYPE_NET_BPS_TX      = "haproxy.bout"
	LB_METRIC_TYPE_CHC_STATUS      = "haproxy.check_status"
	LB_METRIC_TYPE_CHC_CODE        = "haproxy.check_code"
	LB_METRIC_TYPE_LAST_CHC        = "haproxy.last_chk"
	LB_METRIC_TYPE_REQ_RATE        = "haproxy.req_rate"
	LB_METRIC_TYPE_HRSP_COUNT      = "haproxy.hrsp_Nxx"

	BUCKET_METRIC_TYPE_NET_BPS_TX = "oss_netio.bps_sent"
	BUCKET_METRIC_TYPE_NET_BPS_RX = "oss_netio.bps_recv"
	BUCKET_METRIC_TYPE_LATECY     = "oss_latency.req_late"
	BUCKET_METRYC_TYPE_REQ_COUNT  = "oss_req.req_count"
)

var (
	ALL_RDS_METRIC_TYPES = []TMetricType{
		RDS_METRIC_TYPE_CPU_USAGE,
		RDS_METRIC_TYPE_MEM_USAGE,
		RDS_METRIC_TYPE_NET_BPS_RX,
		RDS_METRIC_TYPE_NET_BPS_TX,

		RDS_METRIC_TYPE_DISK_USAGE,
		RDS_METRIC_TYPE_DISK_READ_BPS,
		RDS_METRIC_TYPE_DISK_WRITE_BPS,
		RDS_METRIC_TYPE_DISK_IO_PERCENT,

		RDS_METRIC_TYPE_CONN_COUNT,
		RDS_METRIC_TYPE_CONN_ACTIVE,
		RDS_METRIC_TYPE_CONN_USAGE,
		RDS_METRIC_TYPE_CONN_FAILED,

		RDS_METRIC_TYPE_QPS,
		RDS_METRIC_TYPE_TPS,
		RDS_METRIC_TYPE_INNODB_READ_BPS,
		RDS_METRIC_TYPE_INNODB_WRITE_BPS,
	}

	ALL_VM_METRIC_TYPES = []TMetricType{
		VM_METRIC_TYPE_CPU_USAGE,
		VM_METRIC_TYPE_MEM_USAGE,
		VM_METRIC_TYPE_DISK_USAGE,

		VM_METRIC_TYPE_DISK_IO_READ_BPS,
		VM_METRIC_TYPE_DISK_IO_WRITE_BPS,
		VM_METRIC_TYPE_DISK_IO_READ_IOPS,
		VM_METRIC_TYPE_DISK_IO_WRITE_IOPS,

		VM_METRIC_TYPE_NET_BPS_RX,
		VM_METRIC_TYPE_NET_BPS_TX,

		VM_METRIC_TYPE_EIP_BPS_IN,
		VM_METRIC_TYPE_EIP_BPS_OUT,

		VM_METRIC_TYPE_EIP_PPS_IN,
		VM_METRIC_TYPE_EIP_PPS_OUT,
	}

	ALL_REDIS_METRIC_TYPES = []TMetricType{
		REDIS_METRIC_TYPE_CPU_USAGE,
		REDIS_METRIC_TYPE_MEM_USAGE,
		REDIS_METRIC_TYPE_NET_BPS_RX,
		REDIS_METRIC_TYPE_NET_BPS_TX,
		REDIS_METRIC_TYPE_CONN_USAGE,
		REDIS_METRIC_TYPE_OPT_SES,
		REDIS_METRIC_TYPE_CACHE_KEYS,
		REDIS_METRIC_TYPE_CACHE_EXP_KEYS,
		REDIS_METRIC_TYPE_DATA_MEM_USAGE,
		REDIS_METRIC_TYPE_SERVER_LOAD,
		REDIS_METRIC_TYPE_CONN_ERRORS,
	}

	ALL_LB_METRIC_TYPES = []TMetricType{
		LB_METRIC_TYPE_SNAT_PORT,
		LB_METRIC_TYPE_SNAT_CONN_COUNT,
		LB_METRIC_TYPE_NET_BPS_RX,
		LB_METRIC_TYPE_NET_BPS_TX,
		LB_METRIC_TYPE_CHC_STATUS,
		LB_METRIC_TYPE_CHC_CODE,
		LB_METRIC_TYPE_LAST_CHC,
		LB_METRIC_TYPE_REQ_RATE,
		LB_METRIC_TYPE_HRSP_COUNT,
	}

	ALL_BUCKET_TYPES = []TMetricType{
		BUCKET_METRIC_TYPE_NET_BPS_TX,
		BUCKET_METRIC_TYPE_NET_BPS_RX,
		BUCKET_METRIC_TYPE_LATECY,
		BUCKET_METRYC_TYPE_REQ_COUNT,
	}
)

type MetricListOptions struct {
	ResourceType TResourceType
	MetricType   TMetricType

	ResourceId string
	StartTime  time.Time
	EndTime    time.Time

	Interval int
	Engine   string
}

type MetricValue struct {
	Timestamp time.Time
	Value     float64
	Tags      map[string]string
}

type MetricValues struct {
	Id         string
	Unit       string
	MetricType TMetricType
	Values     []MetricValue
}
