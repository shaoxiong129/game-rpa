package service

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/go-vgo/robotgo"
	"github.com/kbinani/screenshot"
	"gocv.io/x/gocv"
)

// 保存截图到文件 (调试用)
func saveImage(img image.Image, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func main() {
	// 1. 截取主显示器全屏
	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		fmt.Println("截屏失败:", err)
		return
	}

	// 调试：保存截图到本地文件
	saveImage(img, "screen.png")

	// 2. 转换 image.Image 到 gocv.Mat
	imgMat, err := gocv.ImageToMatRGB(img)
	if err != nil {
		fmt.Println("转换截图失败:", err)
		return
	}
	defer imgMat.Close()

	// 3. 读取模板图像（比如按钮图标）
	template := gocv.IMRead("button_template.png", gocv.IMReadColor)
	if template.Empty() {
		fmt.Println("读取模板图失败")
		return
	}
	defer template.Close()

	// 4. 创建结果 Mat (用于存储匹配结果)
	result := gocv.NewMat()
	defer result.Close()

	// 5. 模板匹配，使用归一化相关系数法
	gocv.MatchTemplate(imgMat, template, &result, gocv.TmCcoeffNormed, gocv.NewMat())

	// 6. 找出最大匹配位置和匹配度
	_, maxVal, _, maxLoc := gocv.MinMaxLoc(result)

	fmt.Printf("最大匹配度: %f\n", maxVal)
	threshold := float32(0.8) // 设置阈值，匹配度大于才认为找到了

	if maxVal >= threshold {
		// 计算点击坐标（模板中心点）
		clickX := maxLoc.X + template.Cols()/2
		clickY := maxLoc.Y + template.Rows()/2

		fmt.Printf("找到目标，点击坐标：(%d, %d)\n", clickX, clickY)

		// 7. 模拟鼠标移动并点击
		robotgo.MoveClick(clickX, clickY, "left")
		fmt.Println("点击完成")
	} else {
		fmt.Println("未找到匹配的目标图像")
	}
}
