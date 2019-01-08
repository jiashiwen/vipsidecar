// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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

package cmd

import (
	"errors"
	"fmt"
	common "github.com/jiashiwen/vipsidecar/common"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net"
	"os"
	"sync"
	// "github.com/jdcloud-api/jdcloud-sdk-go/services/vpc/models"
	"time"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vipsidecar",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
	Run: func(cmd *cobra.Command, args []string) {

		configfile, _ := cmd.Flags().GetString("config")
		var wg sync.WaitGroup
		var mutex = &sync.Mutex{}
		//健康检查map[ip地址]是否ping通
		// var healtcheckmap map[string]bool

		if configfile != "" {

			defer os.Exit(0)
			parameter := common.GetConfigParameters(configfile)
			CheckParameter(parameter)
			vpcclient := common.InitVpcClient(parameter.AccessKeyID, parameter.AccessKeySecret)

			for {
				//本地网卡绑定的vip
				vipsonlocal := []string{}

				//当前网络接口与vip绑定关系
				var networkinterfacevips = make(map[common.JdNetworkInterface][]string)
				for _, networkinterface := range parameter.Allnetworkinterfaces {
					wg.Add(1)
					nf := networkinterface
					go func() {
						defer wg.Done()
						secondaryips := common.GetNetworkInterfaceIps(vpcclient, nf.RangId, nf.NetWorkInterfaceId)
						ips := []string{}
						for i := 0; i < len(secondaryips); i++ {
							ips = append(ips, secondaryips[i].PrivateIpAddress)
						}
						mutex.Lock()
						networkinterfacevips[nf] = ips
						mutex.Unlock()

					}()

				}

				wg.Wait()
				// fmt.Println("networkinterfacevips", networkinterfacevips)
				//获取本地vip列表
				localips := GetIntranetIp()
				for _, ip := range localips {
					ok, _ := common.Contain(ip, parameter.Vips)
					if ok {
						vipsonlocal = append(vipsonlocal, ip)
					}
				}

				//如果在本地检查到vip,同时vip的注册网络端口不是本地网路端口，或所有网络端口中都没有注册，则注册vip到本地网络端口,同时删除老旧注册
				for _, localvip := range vipsonlocal {
					vipnotonanyinterface := true
					for k, v := range networkinterfacevips {
						ok, _ := common.Contain(v, localvip)
						if ok {
							vipnotonanyinterface = false
							if k.RangId != parameter.Localnetworkinterface.RangId || k.NetWorkInterfaceId != parameter.Localnetworkinterface.NetWorkInterfaceId {
								common.UnAssignVips(vpcclient, k.RangId, k.NetWorkInterfaceId, []string{localvip})
								common.AssignVips(vpcclient, parameter.Localnetworkinterface.RangId, parameter.Localnetworkinterface.NetWorkInterfaceId, []string{localvip})
							}
						}
					}
					if vipnotonanyinterface {
						common.AssignVips(vpcclient, parameter.Localnetworkinterface.RangId, parameter.Localnetworkinterface.NetWorkInterfaceId, []string{localvip})
					}
				}

				log.Println("vipsonlocal", vipsonlocal)
				log.Println("networkinterfacevips", networkinterfacevips)
				time.Sleep(time.Duration(parameter.Pollinginterval) * time.Second)
			}

		}
		cmd.Help()

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vipsidecar.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".vipsidecar" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".vipsidecar")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

//本地ip列表
func GetIntranetIp() []string {
	localips := []string{}
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		log.Println(err)
		return localips
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localips = append(localips, ipnet.IP.String())
				// fmt.Println("ip:", ipnet.IP.String())
			}
		}
	}
	return localips
}

//配置文件参数检查
func CheckParameter(p *common.Parameters) {
	//检查ak
	if p.AccessKeyID == "" {
		log.Println(errors.New("AccessKeyID must be set"))
		os.Exit(1)
	}

	//检查sk
	if p.AccessKeySecret == "" {
		log.Println(errors.New("AccessKeySecret must be set"))
		os.Exit(1)
	}

	if p.Pollinginterval <= 5 {
		p.Pollinginterval = 5
	}
}
