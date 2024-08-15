package getdata

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
	"github.com/gorilla/websocket"
)

// Start 方法启动数据包捕获
func (c *Catcher) Start() {
    c.mu.Lock()
    defer c.mu.Unlock()
    if c.isRunning {
        return
    }
    c.isRunning = true
    // 初始化其他必要的资源
}

// Stop 方法停止数据包捕获
func (c *Catcher) Stop() {
    c.mu.Lock()
    defer c.mu.Unlock()
    if !c.isRunning {
        return
    }
    c.isRunning = false
    // 清理资源，例如关闭 pcap 句柄等
}

// IsRunning 方法检查 Catcher 是否正在运行
func (c *Catcher) IsRunning() bool {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.isRunning
}

func DumpToPcap(catcher *Catcher, stopCaptureChan chan bool) {
	// 使用网卡接收数据
	handle, err := pcap.OpenLive(catcher.Device, catcher.Snapshot_len, catcher.Promiscuous, catcher.Timeout)
	fmt.Printf("open %s\n", catcher.Device)
	if err != nil {
		log.Fatal(err)
		fmt.Printf("Error opening device %s: %v", catcher.Device, err)
		os.Exit(1)
	}
	defer handle.Close()
	
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	f, err := os.OpenFile("./data/test.pcap", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}

	w := pcapgo.NewWriter(f)
	w.WriteFileHeader(uint32(catcher.Snapshot_len), layers.LinkTypeEthernet)
	defer f.Close()

	// 执行捕获逻辑
	for packet := range packetSource.Packets() {
		// ... 写入数据包到 pcap 文件的代码 ...
		w.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
		catcher.PacketCount++

		select {
		case <-stopCaptureChan:
			fmt.Println("Stop")
			return
		default:
			continue
		}
	}

	// 检查 Catcher 是否仍在运行
	if !catcher.IsRunning() {
		fmt.Println("Packet capture is not running, stopping DumpToPcap.")
		return // Catcher 已停止运行，退出函数
	}
}

func GetPcap(catcher *Catcher, stopCaptureChan chan bool, conn *websocket.Conn) () {
	// 使用网卡接收数据
	handle, err := pcap.OpenLive(catcher.Device, catcher.Snapshot_len, catcher.Promiscuous, catcher.Timeout)
	fmt.Printf("open %s\n", catcher.Device)
	if err != nil {
		log.Fatal(err)
		fmt.Printf("Error opening device %s: %v", catcher.Device, err)
		os.Exit(1)
	}
	defer handle.Close()
	
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	f, err := os.OpenFile("../data/test.pcap", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}

	w := pcapgo.NewWriter(f)
	w.WriteFileHeader(uint32(catcher.Snapshot_len), layers.LinkTypeEthernet)
	defer f.Close()
	var no uint64 = 0
	// 执行捕获逻辑
	for packet := range packetSource.Packets() {
		isIPv6 := true

		for _, layer := range packet.Layers() {
			if layer.LayerType() == layers.LayerTypeIPv6 {
				isIPv6 = true
				break
			}
			isIPv6 = false
		}

		if (isIPv6) {
			continue
		}

		transportLayer := packet.TransportLayer()
		if (transportLayer == nil) {
			continue
		}

		no++
		// ... 写入数据包到 pcap 文件的代码 ...
		w.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
		
		var basicPacketInfo basicPacketInfos
		basicPacketInfo.Id = no
		basicPacketInfo.TimeStamp = packet.Metadata().Timestamp
		// 获取传输层协议（例如 TCP 或 UDP）
		for _, layer := range packet.Layers() {
			switch layer.LayerType() {
			case layers.LayerTypeIPv4:
				ipLayer, _ := layer.(*layers.IPv4)
				ipPayload := ipLayer.Payload
				payload := string(ipPayload)

				basicPacketInfo.Length = uint32(len((payload)))
				basicPacketInfo.SrcIP = ipLayer.SrcIP.String()
				basicPacketInfo.DstIP = ipLayer.DstIP.String()
			}
		}

		if transportLayer != nil {
			basicPacketInfo.Protocol = transportLayer.LayerType().String()
			for _, layer := range packet.Layers() {
				basicPacketInfo.Protocol = layer.LayerType().String()
			}
		}

		// 将数据包转换为 JSON 格式
		packetArrayJSON, err := json.Marshal(basicPacketInfo)
		if err != nil {
			log.Println("json marshal:", err)
			return
		}

		// 发送 JSON 数据通过 WebSocket
		if err := conn.WriteMessage(websocket.TextMessage, packetArrayJSON); err != nil {
			log.Println("write:", err)
			return
		}

		select {
		case <-stopCaptureChan:
			fmt.Println("Stop")
			return
		default:
			continue
		}
	}

	if err != nil {
		log.Fatalf("Failed to marshal packets to JSON: %v", err)
	}

	// 检查 Catcher 是否仍在运行
	if !catcher.IsRunning() {
		fmt.Println("Packet capture is not running, stopping DumpToPcap.")
		return// Catcher 已停止运行，退出函数
	}

	return
}