package main

import (
	"fyne.io/fyne/v2"
	fyneApp "fyne.io/fyne/v2/app"
	"github.com/duxweb/codux/internal/app"
	coduxtheme "github.com/duxweb/codux/internal/ui/theme"
	"github.com/duxweb/codux/internal/ui/views"
)

func main() {
	// 创建 Fyne 应用
	a := fyneApp.New()

	// 创建自定义主题
	coduxTheme := coduxtheme.NewCoduxTheme()
	coduxTheme.SetMode(coduxtheme.ThemeModeDark)
	a.Settings().SetTheme(coduxTheme)

	// 创建应用 Store
	store, err := app.NewStore()
	if err != nil {
		panic(err)
	}
	defer store.Close()

	// 创建主窗口
	w := a.NewWindow("Codux")

	// 创建主视图
	mainView := views.NewMainView(w, store)

	// 设置窗口内容
	w.SetContent(mainView)

	// 设置窗口大小
	w.Resize(fyne.NewSize(1200, 800))

	// 显示窗口
	w.ShowAndRun()
}
