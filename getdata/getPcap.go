package getdata

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
)

// 启动数据包捕获
func (catcher *Catcher) Start() {
    catcher.mu.Lock()
    defer catcher.mu.Unlock()
    if catcher.isRunning {
        return
    }
    catcher.isRunning = true
}

// 停止数据包捕获
func (catcher *Catcher) Stop() {
    catcher.mu.Lock()
    defer catcher.mu.Unlock()
    if !catcher.isRunning {
        return
    }
    catcher.isRunning = false
}

// 检查 Catcher 是否正在运行
func (catcher *Catcher) IsRunning() bool {
    catcher.mu.Lock()
    defer catcher.mu.Unlock()
    return catcher.isRunning
}

// 将抓到的包保存到 pcap 文件中
func DumpToPcap(filepath string, catcher *Catcher, stopCaptureChan chan bool) {
	// 读取捕获的数据包
	handle, err := pcap.OpenLive(catcher.Device, catcher.Snapshot_len, catcher.Promiscuous, catcher.Timeout)
	fmt.Printf("open %s\n", catcher.Device)

	if err != nil {
		log.Fatal(err)
		fmt.Printf("Error opening device %s: %v", catcher.Device, err)
		os.Exit(1)
	}
	defer handle.Close()

	// 获取当前时间
	currentTime := time.Now()

	// 格式化时间字符串作为文件名
	formattedTime := currentTime.Format("2006-01-02_150405")
	filepath = filepath + "/" + formattedTime + ".pcap"

	// 处理从 handle 捕获的数据包
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}

	// 写入 pcap 文件
	w := pcapgo.NewWriter(file)
	w.WriteFileHeader(uint32(catcher.Snapshot_len), layers.LinkTypeEthernet)
	defer file.Close()

	// 执行捕获逻辑
	for packet := range packetSource.Packets() {
		// 写入数据包到 pcap 文件的代码
		fmt.Println(packet)
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
