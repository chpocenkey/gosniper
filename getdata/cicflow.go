package getdata

import (
	"fmt"
	"os/exec"
)

func Cicflow() {

	// 批处理脚本的路径
	batFilePath := "cicflow.bat"

	// 使用 exec.Command 执行批处理脚本
	cmd := exec.Command("cmd", "/c", batFilePath)

	// 获取命令的标准输出和标准错误
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("执行批处理脚本时出错: %s\n", err)
	}

	// 打印输出结果
	fmt.Printf("脚本输出: %s\n", output)
}