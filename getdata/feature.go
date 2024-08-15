package getdata

import (
	"time"
)

type basicPacketInfos struct {
	Id 				uint64
	TimeStamp		time.Time
	SrcIP 			string
	DstIP 			string
	Protocol		string
	Length			uint32
	// info 			string
}
