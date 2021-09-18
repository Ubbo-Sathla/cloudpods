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
