//go:build windows
// +build windows

package service

import (
	"fmt"
	"strings"
)

func run() {
	// windows窗口标题
	winTitle = "雷电模拟器"
	// 参数 nil 代表不指定窗口类名，只按窗口标题匹配。
	// syscall.StringToUTF16Ptr 把 Go 字符串转换成 Windows API 需要的 UTF16 指针。
	// 这个句柄 hwnd 只是演示用，实际后续枚举用的才是你真正想找的模拟器窗口。
	hwnd := win.FindWindow(nil, syscall.StringToUTF16Ptr(winTitle))

	// 声明一个切片用来保存符合条件的窗口句柄。
	var foundWindows []win.HWND

	// 创建一个回调函数，传给 EnumWindows 枚举所有顶层窗口时调用。
	// syscall.NewCallback 是 Go 绑定 Windows API 回调的标准做法。
	// 这个回调参数：
	// hwnd 是当前枚举到的窗口句柄
	// lparam 是你传给 EnumWindows 的自定义参数（这里是0，没用）
	cb := syscall.NewCallback(func(hwnd win.HWND, lparam uintptr) uintptr {
		// 先用 GetWindowTextLength 获取当前窗口标题长度，+1留给结尾符。
		length := win.GetWindowTextLength(hwnd) + 1
		// 申请一个 uint16 类型的切片作为缓冲区（Windows 字符是 UTF-16）。
		buf := make([]uint16, length)
		// 用 GetWindowText 读取窗口标题到缓冲区。
		win.GetWindowText(hwnd, &buf[0], int32(length))
		// 转成 Go 字符串。
		title := syscall.UTF16ToString(buf)

		// 判断窗口是否可见 IsWindowVisible(hwnd)，并且窗口标题里包含子串 "雷电模拟器"。
		if win.IsWindowVisible(hwnd) && strings.Contains(title, "雷电模拟器") {
			fmt.Printf("Found window: hwnd=0x%x title=%s\n", hwnd, title)
			// 如果符合条件，则打印窗口信息，并把句柄加入 foundWindows。
			foundWindows = append(foundWindows, hwnd)
		}
		// EnumWindows 的回调必须返回非0值继续枚举，返回0则停止。这里返回1表示继续枚举剩下的窗口。
		return 1 // continue enumeration
	})
	// 调用 Windows API EnumWindows，让系统遍历所有顶层窗口，针对每个窗口调用上面定义的回调 cb。
	win.EnumWindows(cb, 0)

	// 枚举完成后判断有没有找到符合条件的窗口。
	// 如果没找到，打印提示。
	// 找到的话，调用 SetForegroundWindow 把第一个符合条件的窗口激活（置于最前）。
	if len(foundWindows) == 0 {
		fmt.Println("No 雷电模拟器 window found.")
	} else {
		// 示例：激活第一个找到的窗口
		win.SetForegroundWindow(foundWindows[0])
	}

}
