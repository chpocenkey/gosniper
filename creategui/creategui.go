package creategui

import (
	"gosniper/getdata"
	"os"

	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func CreateGUI() {
	// 创建应用程序
	myApp := app.New()
	myWindow := myApp.NewWindow("Flow Detect")
	
	// 获取网络接口
	interfaces := getdata.GetDevices()
	var interfaceList []string
	for interfaceDes := range interfaces {
		interfaceList = append(interfaceList, interfaceDes)
	}

	var stopCaptureChan chan bool
	var netcard string
	var capture *getdata.Catcher
	var savePath string

	// 创建一个下拉列表来显示网络接口
	interfaceSelecter := widget.NewSelect(interfaceList, func(selected string) {
		netcard = interfaces[selected]
		capture = getdata.NewCatcher(netcard)
	})

	// 开始抓包按钮
	startCaptureButton := widget.NewButton("Start Capture", func() {
		getdata.StrBind.Set("Packet Number: 0")
		capture.PacketCount = 0
		if capture != nil && savePath != ""{
            stopCaptureChan = make(chan bool) // 创建一个停止通道
            go func() {
                capture.Start() 							// 启动捕获
                getdata.DumpToPcap(savePath, capture, stopCaptureChan) // 抓包，并保存到 pcap 文件中
            }()
        }
	})

	// 停止抓包按钮
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

	// 选择保存路径按钮
	chooseSavePathButton := widget.NewButton("Choose Save Path", func() {
		folderDialog := dialog.NewFolderOpen(func(folder fyne.ListableURI, err error) {
			if err != nil {
				fyne.LogError("Folder selection error:", err)
				return
			}
			if folder == nil {
				log.Println("Folder selection cancelled")
				return
			}
			savePath = folder.Path() // 设置保存路径
		}, myWindow)
		
		// 获取当前路径
		currentPath, err := os.Getwd()
		if err != nil {
			log.Println("Folder selection cancelled")
			return
		}

		// 获取默认存储位置
		dirPath := currentPath + "/data"
		// 检查文件夹是否存在
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			// 文件夹不存在，创建文件夹
			err := os.MkdirAll(dirPath, 0755) // 0755 是文件夹权限，可根据需要设置
			if err != nil {
				panic(err)
			}
		}
		// 指定相对路径
		addr, err := storage.ListerForURI(storage.NewFileURI(dirPath))
		if err != nil {
			fyne.LogError("Error converting to ListableURI:", err)
			return
		}
		folderDialog.SetLocation(addr) // 设置默认路径
		folderDialog.Show()
	})
	
	// 创建一个标签，并使用WithData方法连接绑定
	captureNum := widget.NewLabelWithData(getdata.StrBind)

	// 创建一个垂直布局的盒子，用于防止布局
	vBox := container.NewVBox(
		widget.NewLabelWithStyle("Flow Detect", fyne.TextAlignCenter, fyne.TextStyle{}),

		// 网卡选择列表
		interfaceSelecter,

		// 选择文件保存路径
		chooseSavePathButton,

		// 抓包数量
		captureNum,

		// 开始抓包按钮
		startCaptureButton,

		// 停止抓包按钮
		stopCaptureButton,
	)

	// 设置窗口的内容
	myWindow.SetContent(vBox)
	myWindow.Resize(fyne.NewSize(400, 600))

	// 显示窗口
	myWindow.ShowAndRun()
}