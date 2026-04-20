package main

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	internalapp "github.com/duxweb/codux/internal/app"
	"github.com/duxweb/codux/internal/ui/menu"
	coduxtheme "github.com/duxweb/codux/internal/ui/theme"
	"github.com/duxweb/codux/internal/ui/views"
)

var Version = "0.2.0"

func main() {
	// 创建 Fyne 应用
	a := app.NewWithID("io.github.duxweb.codux")

	// 创建自定义主题
	coduxTheme := coduxtheme.NewCoduxTheme()
	coduxTheme.SetMode(coduxtheme.ThemeModeDark)
	a.Settings().SetTheme(coduxTheme)

	// 创建应用 Store
	store, err := internalapp.NewStore()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create store: %v\n", err)
		os.Exit(1)
	}
	defer store.Close()

	// 创建主窗口
	w := a.NewWindow("Codux")

	// 设置菜单栏
	mainMenu := menu.NewAppMenu(a)
	w.SetMainMenu(mainMenu)

	// 创建主视图
	mainView := views.NewMainView(w, store)

	// 设置窗口内容
	w.SetContent(mainView)

	// 设置窗口大小
	w.Resize(fyne.NewSize(1200, 800))

	// 窗口关闭处理器
	w.SetOnClosed(func() {
		store.Close()
	})

	// 显示窗口
	w.ShowAndRun()
}
