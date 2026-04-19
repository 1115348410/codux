package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/duxweb/codux/internal/ui/views"
)

func main() {
	a := app.New()
	w := a.NewWindow("Codux")

	// 创建主视图
	mainView := views.NewMainView(w)

	// 设置窗口内容
	w.SetContent(mainView)

	// 设置最小窗口尺寸
	w.Resize(fyne.NewSize(1024, 768))

	// 显示窗口
	w.ShowAndRun()
}
