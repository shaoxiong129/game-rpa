package service

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"syscall"
	"unsafe"

	"github.com/go-vgo/robotgo"
	"github.com/kbinani/screenshot"
	"gocv.io/x/gocv"
	"golang.org/x/sys/windows"
)

var (
	user32            = windows.NewLazySystemDLL("user32.dll")
	procGetWindowRect = user32.NewProc("GetWindowRect")
)

// RECT结构体
type RECT struct {
	Left, Top, Right, Bottom int32
}

// 根据窗口句柄获取窗口矩形（屏幕坐标）
func GetWindowRect(hwnd uintptr) (image.Rectangle, error) {
	var rect RECT
	ret, _, err := procGetWindowRect.Call(hwnd, uintptr(unsafe.Pointer(&rect)))
	if ret == 0 {
		return image.Rectangle{}, err
	}
	return image.Rect(int(rect.Left), int(rect.Top), int(rect.Right), int(rect.Bottom)), nil
}

// 保存图片，调试用
func saveImage(img image.Image, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

// 通过窗口句柄+模板路径匹配点击，hwnd=0则全屏截取
func ClickByImageWithWindow(hwnd uintptr, templatePath string, threshold float32) (bool, error) {
	var bounds image.Rectangle
	var err error

	if hwnd == 0 {
		bounds = screenshot.GetDisplayBounds(0) // 全屏
	} else {
		bounds, err = GetWindowRect(hwnd)
		if err != nil {
			return false, fmt.Errorf("获取窗口位置失败: %w", err)
		}
	}

	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return false, fmt.Errorf("截屏失败: %w", err)
	}

	// 可选保存调试截图
	// _ = saveImage(img, "screen.png")

	imgMat, err := gocv.ImageToMatRGB(img)
	if err != nil {
		return false, fmt.Errorf("截图转换失败: %w", err)
	}
	defer imgMat.Close()

	template := gocv.IMRead(templatePath, gocv.IMReadColor)
	if template.Empty() {
		return false, fmt.Errorf("模板图读取失败: %s", templatePath)
	}
	defer template.Close()

	result := gocv.NewMat()
	defer result.Close()

	gocv.MatchTemplate(imgMat, template, &result, gocv.TmCcoeffNormed, gocv.NewMat())

	_, maxVal, _, maxLoc := gocv.MinMaxLoc(result)

	if maxVal >= threshold {
		// 计算相对窗口的点击点，加上窗口左上角坐标得到屏幕坐标
		clickX := maxLoc.X + template.Cols()/2 + bounds.Min.X
		clickY := maxLoc.Y + template.Rows()/2 + bounds.Min.Y

		robotgo.MoveClick(clickX, clickY, "left")
		return true, nil
	}

	return false, fmt.Errorf("未找到匹配目标，最大匹配度: %f", maxVal)
}

func main() {
	// 示例：你得先用你已有方法获取雷电模拟器窗口句柄 hwnd
	var hwnd uintptr = 0x00123456 // 替换成真实窗口句柄

	ok, err := ClickByImageWithWindow(hwnd, "button_template.png", 0.8)
	if ok {
		fmt.Println("点击成功")
	} else {
		fmt.Println("点击失败:", err)
	}
}
