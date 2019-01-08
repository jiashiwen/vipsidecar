package common

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Parameters struct {
	AccessKeyID           string               `yaml:"accessskeyid"`
	AccessKeySecret       string               `yaml:"accesskeysecret"`
	Vips                  []string             `yaml:"vips"`
	Allnetworkinterfaces  []JdNetworkInterface `yaml:"allnetworkinterfaces"`
	Localnetworkinterface JdNetworkInterface   `yaml:"localnetworkinterface"`
	Pollinginterval       int                  `yaml:"pollinginterval"`
}

type JdNetworkInterface struct {
	RangId             string `yaml:"rangid"`
	NetWorkInterfaceId string `yaml:"networkinterfaceid"`
}

func GetConfigParameters(configfile string) *Parameters {

	parameters := new(Parameters)
	yamlFile, err := ioutil.ReadFile(configfile)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		os.Exit(1)
	}
	err = yaml.Unmarshal(yamlFile, parameters)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		os.Exit(1)
	}
	return parameters
}
