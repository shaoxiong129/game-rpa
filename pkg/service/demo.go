//go:build windows
// +build windows

package service

// 查找所有窗口标题包含指定关键字且可见的窗口句柄
func FindWindowsByTitleKeyword(keyword string) []win.HWND {
	var foundWindows []win.HWND

	cb := syscall.NewCallback(func(hwnd win.HWND, lparam uintptr) uintptr {
		if !win.IsWindowVisible(hwnd) {
			return 1 // 继续枚举
		}

		length := win.GetWindowTextLength(hwnd) + 1
		buf := make([]uint16, length)
		win.GetWindowText(hwnd, &buf[0], int32(length))
		title := syscall.UTF16ToString(buf)

		if strings.Contains(title, keyword) {
			fmt.Printf("Found window: hwnd=0x%x title=%s\n", hwnd, title)
			foundWindows = append(foundWindows, hwnd)
		}

		return 1 // 继续枚举
	})

	win.EnumWindows(cb, 0)
	return foundWindows
}

// 激活指定窗口
func ActivateWindow(hwnd win.HWND) bool {
	ret := win.SetForegroundWindow(hwnd)
	return ret != 0
}

// 调整窗口大小和位置
func MoveResizeWindow(hwnd win.HWND, x, y, width, height int32) bool {
	// SWP_NOZORDER = 0x4, SWP_SHOWWINDOW = 0x40
	flags := uint32(win.SWP_NOZORDER | win.SWP_SHOWWINDOW)
	return win.SetWindowPos(hwnd, 0, x, y, width, height, flags) != 0
}

// 模拟鼠标左键点击指定屏幕坐标
func MouseLeftClick(x, y int32) {
	// 设置鼠标位置
	win.SetCursorPos(int32(x), int32(y))
	time.Sleep(50 * time.Millisecond) // 稍微等待

	// 鼠标按下和抬起事件
	win.MouseEvent(win.MOUSEEVENTF_LEFTDOWN, 0, 0, 0, 0)
	time.Sleep(50 * time.Millisecond)
	win.MouseEvent(win.MOUSEEVENTF_LEFTUP, 0, 0, 0, 0)
}

// 获取窗口矩形（位置和大小）
func GetWindowRect(hwnd win.HWND) (x, y, width, height int32, ok bool) {
	var rect win.RECT
	if win.GetWindowRect(hwnd, &rect) {
		x = rect.Left
		y = rect.Top
		width = rect.Right - rect.Left
		height = rect.Bottom - rect.Top
		return x, y, width, height, true
	}
	return 0, 0, 0, 0, false
}

func main() {
	keyword := "雷电模拟器"
	windows := FindWindowsByTitleKeyword(keyword)
	if len(windows) == 0 {
		fmt.Printf("没有找到包含关键字 [%s] 的窗口\n", keyword)
		return
	}

	hwnd := windows[0]

	// 激活窗口
	if ActivateWindow(hwnd) {
		fmt.Println("窗口已激活")
	} else {
		fmt.Println("激活窗口失败")
	}

	// 调整窗口大小位置
	if MoveResizeWindow(hwnd, 100, 100, 800, 600) {
		fmt.Println("窗口大小位置调整成功")
	} else {
		fmt.Println("窗口调整失败")
	}

	// 取窗口当前坐标用于模拟点击（示范点击窗口左上角偏右下点）
	x, y, w, h, ok := GetWindowRect(hwnd)
	if ok {
		clickX := x + w/4
		clickY := y + h/4
		fmt.Printf("模拟点击窗口坐标 (%d, %d)\n", clickX, clickY)
		MouseLeftClick(clickX, clickY)
		fmt.Println("模拟点击完成")
	} else {
		fmt.Println("获取窗口位置失败，无法模拟点击")
	}
}
