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

package notify

type SWorkwxSendMessageResp struct {
	ErrCode        int    `json:"errcode"`
	ErrMsg         string `json:"errmsg"`
	InvalidUser    string `json:"invaliduser"`
	InvalidParty   string `json:"invalidparty" `
	InvalidTag     string `json:"invalidtag"`
	UnlicensedUser string `json:"unlicenseduser" `
	MsgId          string `json:"msgid"`
	ResponseCode   string `json:"response_code"`
}
