package ovn

import (
	"context"
	"fmt"
	apis "yunion.io/x/onecloud/pkg/apis/compute"
	agentmodels "yunion.io/x/onecloud/pkg/vpcagent/models"
	"yunion.io/x/ovsdb/schema/ovn_nb"
	"yunion.io/x/ovsdb/types"
)

func (keeper *OVNNorthboundKeeper) ClaimVpp(ctx context.Context, vpc *agentmodels.Vpc) error {
	var (
		args      []string
		ocVersion = fmt.Sprintf("%s.%d", vpc.UpdatedAt, vpc.UpdateVersion)
	)

	irows := []types.IRow{}
	var (
		vpp      = true
		vpcVppLs *ovn_nb.LogicalSwitch
		vpcRvppp *ovn_nb.LogicalRouterPort
		vpcVpprp *ovn_nb.LogicalSwitchPort
		vpcVppep *ovn_nb.LogicalSwitchPort
	)
	if vpp {
		vpcVppLs = &ovn_nb.LogicalSwitch{
			Name: vpcVppLsName(vpc.Id),
		}

		vpcRvppp = &ovn_nb.LogicalRouterPort{
			Name:     vpcRvpppName(vpc.Id),
			Mac:      apis.VpcInterVppMac1,
			Networks: []string{fmt.Sprintf("%s/%d", apis.VpcInterVppIP1(), apis.VpcInterVppMask)},
		}

		vpcVpprp = &ovn_nb.LogicalSwitchPort{
			Name:      vpcVpprpName(vpc.Id),
			Type:      "router",
			Addresses: []string{"router"},
			Options: map[string]string{
				"router-port": vpcRvpppName(vpc.Id),
			},
		}

		vpcVppep = &ovn_nb.LogicalSwitchPort{
			Name:      vpcVppepName(vpc.Id),
			Addresses: []string{fmt.Sprintf("%s %s", apis.VpcInterVppMac2, apis.VpcInterVppIP2().String())},
		}
		irows = append(irows,
			vpcVppLs,
			vpcRvppp,
			vpcVpprp,
			vpcVppep,
		)
	}
	allFound, args := cmp(&keeper.DB, ocVersion, irows...)
	if allFound {
		return nil
	}

	if vpp {
		args = append(args, ovnCreateArgs(vpcVppLs, vpcVppLs.Name)...)
		args = append(args, ovnCreateArgs(vpcRvppp, vpcRvppp.Name)...)
		args = append(args, ovnCreateArgs(vpcVpprp, vpcVpprp.Name)...)
		args = append(args, ovnCreateArgs(vpcVppep, vpcVppep.Name)...)
		args = append(args, "--", "add", "Logical_Switch", vpcVppLs.Name, "ports", "@"+vpcVpprp.Name)
		args = append(args, "--", "add", "Logical_Router", vpcLrName(vpc.Id), "ports", "@"+vpcRvppp.Name)
		args = append(args, "--", "add", "Logical_Switch", vpcVppLs.Name, "ports", "@"+vpcVppep.Name)

	}
	return keeper.cli.Must(ctx, "ClaimVpp", args)

}

// vpp

func vpcVppLsName(vpcId string) string {
	return fmt.Sprintf("vpc-vpp/%s", vpcId)
}

func vpcRvpppName(vpcId string) string {
	return fmt.Sprintf("vpc-rvpp/%s", vpcId)
}

func vpcVpprpName(vpcId string) string {
	return fmt.Sprintf("vpc-vppr/%s", vpcId)
}

func vpcVppepName(vpcId string) string {
	return fmt.Sprintf("vpc-vppe/%s", vpcId)
}
