package getdata

import (
	"fmt"
	"log"

	"github.com/google/gopacket/pcap"
)

// 获取网卡
func GetDevices() (map[string]string){
	devices, err := pcap.FindAllDevs()							// 获取所有网卡信息
	deviceDescriptions := make(map[string]string)				// 创建网卡信息和描述的 map
	if err != nil {
	    log.Fatal(err)
	}
	for _, device := range devices {
	  	deviceDescriptions[device.Description] = device.Name
	}
	return deviceDescriptions
}
 
//获取所有有IP的网络设备名称
func GetAllDevsHaveIPAddress() (map[string]string) {
	devices, err := pcap.FindAllDevs()							// 获取所有网卡信息
	deviceDescriptions := make(map[string]string)				// 创建网卡信息和描述的 map
	if err != nil {
	    log.Fatal(err)
	}
	for _, device := range devices {
		fmt.Println("Dev:", device.Name, "\tDes:", device.Description)
		if len(device.Addresses) > 0 {
			for _, ips := range device.Addresses {
				deviceDescriptions[device.Description] = device.Name
				fmt.Println("\tAddr:", ips.IP.String())
			}

		}
	}
	return deviceDescriptions
}
