package getdata

import (
	"fmt"
	"log"
	"strings"

	"github.com/google/gopacket/pcap"
)


func GetAllDevices() (map[string]string){
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	deviceDescriptions := make(map[string]string)
	for _, device := range devices {
		deviceDescriptions[device.Name] = device.Description
	}

	return deviceDescriptions
}

func GetDevices() (ret map[string]string){
	devices, err := pcap.FindAllDevs()
	deviceDescriptions := make(map[string]string)
	if err != nil {
	    log.Fatal(err)
	}
	// fmt.Println("找到网卡信息: ")
	for _, device := range devices {
	//   fmt.Println("\n名字: ", device.Name)
	//   fmt.Println("描述信息: ", device.Description)
	  deviceDescriptions[device.Description] = device.Name
	//   fmt.Println("网卡地址信息: ", device.Addresses)
	//   for _, address := range device.Addresses {
	// 	fmt.Println("- IP 地址为: ", address.IP)
	// 	fmt.Println("- 掩码为: ", address.Netmask)
	// 	fmt.Println()
	//   }
	}
	return deviceDescriptions
}
 
//获取所有有IP的网络设备名称
func GetAllDevsHaveAddress() (string, error) {
	pcapDevices, err := pcap.FindAllDevs()
	if err != nil {
		return "", fmt.Errorf("获取失败: %s", err.Error())
	}
	var buf strings.Builder
	for _, dev := range pcapDevices {
		fmt.Println("Dev:", dev.Name, "\tDes:", dev.Description)
		buf.WriteString(dev.Name)
		if len(dev.Addresses) > 0 {
			for _, ips := range dev.Addresses {
				fmt.Println("\tAddr:", ips.IP.String())
				//buf.WriteString(ips.IP.String())
			}
 
		}
		buf.WriteString("\n")
	}
	return buf.String(), nil
}
