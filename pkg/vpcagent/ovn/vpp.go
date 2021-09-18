package ovn

import (
	"context"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	apis "yunion.io/x/onecloud/pkg/apis/compute"
	agentmodels "yunion.io/x/onecloud/pkg/vpcagent/models"
	"yunion.io/x/ovsdb/schema/ovn_nb"
	"yunion.io/x/ovsdb/types"
)

type VppLogicalRouterStaticRoute []struct {
	VpcID string  `yaml:"vpcID"`
	Route []Route `yaml:"route"`
}
type Route struct {
	Policy                   string `yaml:"policy"`
	IPPrefix                 string `yaml:"ipPrefix"`
	Name                     string `yaml:"name"`
	LogicalRouterStaticRoute *ovn_nb.LogicalRouterStaticRoute
}

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

func (keeper *OVNNorthboundKeeper) ClaimVppLogicalRouterStaticRoute(ctx context.Context, vpc *agentmodels.Vpc) error {

	var (
		args      []string
		ocVersion = fmt.Sprintf("%s.%d", vpc.UpdatedAt, vpc.UpdateVersion)
	)

	yamlFile, err := ioutil.ReadFile("/etc/yunion/vpp/vpp.yaml")
	if err != nil {
		return errors.New(fmt.Sprintf("Error reading YAML file: %s", err))

	}
	var yamlConfig VppLogicalRouterStaticRoute
	err = yaml.Unmarshal(yamlFile, &yamlConfig)
	if err != nil {
		return errors.New(fmt.Sprintf("Error parsing YAML file: %s", err))
	}
	routes := []Route{}

	for i := range yamlConfig {
		if yamlConfig[i].VpcID == vpc.Id {
			routes = yamlConfig[i].Route
		}
	}
	if len(routes) == 0 {
		return nil
	}

	log.Println(routes)

	irows := []types.IRow{}

	for i := range routes {
		routes[i].LogicalRouterStaticRoute = &ovn_nb.LogicalRouterStaticRoute{
			Policy:     ptr(routes[i].Policy),
			IpPrefix:   routes[i].IPPrefix,
			Nexthop:    apis.VpcInterVppIP2().String(),
			OutputPort: ptr(vpcRvpppName(vpc.Id)),
		}
		irows = append(irows,
			routes[i].LogicalRouterStaticRoute,
		)
	}

	allFound, args := cmp(&keeper.DB, ocVersion, irows...)
	if allFound {
		return nil
	}

	for i := range routes {
		args = append(args, ovnCreateArgs(routes[i].LogicalRouterStaticRoute, fmt.Sprintf("%s%d", routes[i].Name, i))...)
		args = append(args, "--", "add", "Logical_Router", vpcLrName(vpc.Id), "static_routes", fmt.Sprintf("@%s%d", routes[i].Name, i))
	}
	return keeper.cli.Must(ctx, "ClaimVppLogicalRouterStaticRoute", args)

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
