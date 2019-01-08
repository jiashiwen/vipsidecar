package common

import (
	"github.com/jdcloud-api/jdcloud-sdk-go/core"
	"github.com/jdcloud-api/jdcloud-sdk-go/services/vpc/apis"
	"github.com/jdcloud-api/jdcloud-sdk-go/services/vpc/client"
	"github.com/jdcloud-api/jdcloud-sdk-go/services/vpc/models"
	"log"
)

type DefaultLogger struct {
	Level int
}

func (logger DefaultLogger) Log(level int, message ...interface{}) {
	if level <= logger.Level {
		// fmt.Println(message...)
	}
}

func InitVpcClient(accessKey string, secretKey string) *client.VpcClient {
	defaultlogger := DefaultLogger{}
	defaultlogger.Level = 1
	credentials := core.NewCredentials(accessKey, secretKey)
	vpcclient := client.NewVpcClient(credentials)
	vpcclient.SetLogger(defaultlogger)
	return vpcclient
}

//获取网卡上的SecondaryIps
func GetNetworkInterfaceIps(client *client.VpcClient, regionId string, network_interface_id string) []models.NetworkInterfacePrivateIp {
	networkinterfacereq := apis.NewDescribeNetworkInterfaceRequest(regionId, network_interface_id)
	nirespons, err := client.DescribeNetworkInterface(networkinterfacereq)
	if err != nil {
		log.Fatalln(err)
	}
	return nirespons.Result.NetworkInterface.SecondaryIps

}

//为网卡注册sencondaryip
func AssignVips(client *client.VpcClient, regionId string, network_interface_id string, ips []string) {
	assignsencondaryipsreq := apis.NewAssignSecondaryIpsRequest(regionId, network_interface_id)
	assignsencondaryipsreq.SecondaryIps = ips
	respons, err := client.AssignSecondaryIps(assignsencondaryipsreq)
	if err != nil {
		log.Println(err)
	}
	log.Println(respons)
}

//为网卡注销sencondaryip
func UnAssignVips(client *client.VpcClient, regionId string, network_interface_id string, ips []string) {
	unassignsecondaryipsreq := apis.NewUnassignSecondaryIpsRequest(regionId, network_interface_id)
	unassignsecondaryipsreq.SecondaryIps = ips
	client.UnassignSecondaryIps(unassignsecondaryipsreq)
}

//查看NetworkInterface是否绑定某一sencondaryip
func IpExistsOnInterface(client *client.VpcClient, regionId string, network_interface_id string, ip string) bool {
	exists := false
	networkinterfacereq := apis.NewDescribeNetworkInterfaceRequest(regionId, network_interface_id)
	nirespons, err := client.DescribeNetworkInterface(networkinterfacereq)
	if err != nil {
		log.Fatalln(err)
		return exists
	}
	ips := nirespons.Result.NetworkInterface.SecondaryIps
	for _, sechonderyip := range ips {
		if sechonderyip.PrivateIpAddress == ip {
			exists = true
		}
	}

	return exists
}
