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

package compute

import (
	"yunion.io/x/pkg/util/netutils"
)

const (
	VpcInterVppMask = 30
	sVpcInterVppIP1 = "100.66.0.1"
	sVpcInterVppIP2 = "100.66.0.2"
	VpcInterVppMac1 = "ee:ee:ee:ee:e0:f0"
	VpcInterVppMac2 = "ee:ee:ee:ee:e0:f1"
)

var (
	vpcInterVppIP1 netutils.IPV4Addr
	vpcInterVppIP2 netutils.IPV4Addr
)

func VpcInterVppIP1() netutils.IPV4Addr {
	return vpcInterVppIP1
}
func VpcInterVppIP2() netutils.IPV4Addr {
	return vpcInterVppIP2
}

func init() {
	mi := func(v netutils.IPV4Addr, err error) netutils.IPV4Addr {
		if err != nil {
			panic(err.Error())
		}
		return v
	}

	vpcInterVppIP1 = mi(netutils.NewIPV4Addr(sVpcInterVppIP1))
	vpcInterVppIP2 = mi(netutils.NewIPV4Addr(sVpcInterVppIP2))

}
