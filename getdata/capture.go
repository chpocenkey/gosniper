package getdata

import (
	"sync"
	"time"
)

// 抓包的基本设置
type Catcher struct {
	mu      		sync.Mutex           								// 用于同步的互斥锁
	isRunning 		bool                								// 表示是否正在抓包的标志
	Device       	string  											// 网卡名称
	Snapshot_len 	int32  				`default:"1024"`                // 每个数据包读取的最大长度
	Promiscuous  	bool   				`default:"false"`               // 是否将网口设置为混杂模式，即是否接收目的不为本机的包
	Timeout      	time.Duration 		`default:"-1 * time.Second"` 	// 设置抓包返回的超时时间，如果设置成30s，即每30s刷新下数据包，设置为负数，就立即刷新数据包
	PacketCount  	int           		`default:"-1"`					// 设置抓包的总包数量，设置为负数表示不限制抓包总量
}

// 初始化抓包设置（设置网卡）
func NewCatcher(device string) *Catcher {
	return &Catcher{
		mu:      sync.Mutex{},
		isRunning: false,
		Device:   device,
	}
}