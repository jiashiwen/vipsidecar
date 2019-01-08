使用方面vipsidecar遵循简洁原则，只需要编写配置文件然后启动即可，我们先看看配置文件的格式
config.yaml
```
accessskeyid: xxxxxxxxxxxxxxxxxxxxxxxx
accesskeysecret: xxxxxxxxxxxxxxxxxxxxxx
vips: 
- 10.0.0.30
- 10.0.0.11
- 10.0.0.21
- 192.168.1.50
allnetworkinterfaces: 
- rangid: cn-east-2
  networkinterfaceid: port-fet9c4pvzt
- rangid: cn-east-2
  networkinterfaceid: port-pig3p7864x

localnetworkinterface: 
  rangid: cn-east-2
  networkinterfaceid: port-fet9c4pvzt
```

|参数|描述|
|---|---|
|accessskeyid|访问密钥ID|
|accesskeysecret|与访问密钥ID结合使用的密钥|
|vips|vip列表|
|allnetworkinterfaces|各个节点上所有可能绑定vip的portid,相关信息可以在控制台查询|
|localnetworkinterface|本机用于绑定vip的网络设备pordid|

* 测试方法
* 京东云申请两台云主机，并保证两台主机可以访问公网，并绑定弹性网卡，此时每台云主机上应该有两块网卡(eth0、eth1),eth1为弹性网卡。
* 编写配置文件config.yaml
```
accessskeyid: your_ak
accesskeysecret: your_sk
vips: 
- vip1
- vip2
- vip3
allnetworkinterfaces: 
- rangid: like_cn-east-2
  networkinterfaceid: like_port-fet9c4pvzt
- rangid: cn-east-2
  networkinterfaceid: like_port-pig3p7864x

localnetworkinterface: 
  rangid: like_cn-east-2
  networkinterfaceid: like_port-fet9c4pvzt
```
* 启动vipsidecar
```
./vipsidecar --config config.yaml
```

* node1手动为网卡添加vip,命令如下
```
ifconfig eth1:0 vip1 netmask 255.255.255.0 up
```
此时可以查看控制台是否为指定的弹性网卡绑定了vip,或者直接通过同网段其他主机ping vip1

* 模拟vip漂移
    * 删除node1上的vip1
        ```
        ifconfig eth1:0 down
        ```
    * node2上添加vip1
        ```
        ifconfig eth1:0 vip1 netmask 255.255.255.0 up
        ```
此时可以查看控制台是否为指定的弹性网卡绑定了vip,或者直接通过同网段其他主机ping vip1
