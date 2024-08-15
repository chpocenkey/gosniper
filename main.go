package main

import (
	"gosniper/creategui"
	"os"
	"strings"
	"github.com/flopp/go-findfont"
)

func init() {
    fontPaths := findfont.List() // findfont 是用于查找字体的库
    for _, path := range fontPaths {
        if strings.Contains(path, "msyh.ttf") || // 微软雅黑
            strings.Contains(path, "simhei.ttf") || // 黑体
            strings.Contains(path, "simsun.ttc") || // 宋体
            strings.Contains(path, "simkai.ttf") {  // 楷体
            os.Setenv("FYNE_FONT", path)
            break
        }
    }
}

func main() {
	creategui.CreateGUI()
}
