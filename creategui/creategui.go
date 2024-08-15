package creategui

import (
	"gosniper/getdata"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func CreateGUI() {
	// 创建 Fyne 应用程序
	myApp := app.New()
	myWindow := myApp.NewWindow("Flow Detext")
	
	// 网络接口
	interfaces := getdata.GetDevices()
	var interfaceList []string
	for interfaceDes := range interfaces {
		interfaceList = append(interfaceList, interfaceDes)
	}

	var stopCaptureChan chan bool
	var netcard string
	var capture *getdata.Catcher

	// 创建一个下拉列表来显示网络接口
	interfaceSelecter := widget.NewSelect(interfaceList, func(selected string) {
		netcard = interfaces[selected]
		capture = getdata.NewCatcher(netcard)
	})

	startCaptureButton := widget.NewButton("Start Capture", func() {
		if capture != nil {
            stopCaptureChan = make(chan bool) // 创建一个停止通道
            go func() {
                capture.Start() // 启动捕获
                getdata.DumpToPcap(capture, stopCaptureChan) // 抓包
            }()
        }
	})

	stopCaptureButton := widget.NewButton("Stop Capture", func() {
        if capture != nil {
            if capture.IsRunning() && stopCaptureChan != nil {
				capture.Stop() // 停止捕获
                stopCaptureChan <- true // 发送停止信号
                close(stopCaptureChan) // 关闭通道
				stopCaptureChan = nil
            }
        }
    })

	// // 设置选择的选项
	// interfaceSelecter.SetSelected(interfaceList[0])

	// 创建一个垂直布局的盒子
	vBox := container.NewVBox(
		widget.NewLabelWithStyle("Flow Detect", fyne.TextAlignCenter, fyne.TextStyle{}),

		// 网卡选择列表
		interfaceSelecter, 

		startCaptureButton,

		widget.NewLabelWithStyle("Capture Number: ", fyne.TextAlignCenter, fyne.TextStyle{}),

		stopCaptureButton,
	)

	// 设置窗口的内容
	myWindow.SetContent(vBox)
	myWindow.Resize(fyne.NewSize(400, 600))

	// 显示窗口
	myWindow.ShowAndRun()
}