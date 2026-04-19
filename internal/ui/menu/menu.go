package menu

import (
	"fyne.io/fyne/v2"
)

// NewAppMenu 创建应用菜单
func NewAppMenu(app fyne.App) *fyne.MainMenu {
	// 文件菜单
	fileMenu := fyne.NewMenu("文件",
		fyne.NewMenuItem("新建项目", func() {}),
		fyne.NewMenuItem("打开文件夹", func() {}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("关闭项目", func() {}),
		fyne.NewMenuItem("关闭所有项目", func() {}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("退出", func() {
			app.Quit()
		}),
	)

	// 编辑菜单
	editMenu := fyne.NewMenu("编辑",
		fyne.NewMenuItem("撤销", func() {}),
		fyne.NewMenuItem("重做", func() {}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("剪切", func() {}),
		fyne.NewMenuItem("复制", func() {}),
		fyne.NewMenuItem("粘贴", func() {}),
	)

	// 视图菜单
	viewMenu := fyne.NewMenu("视图",
		fyne.NewMenuItem("放大", func() {}),
		fyne.NewMenuItem("缩小", func() {}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("重置缩放", func() {}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("全屏", func() {}),
	)

	// 帮助菜单
	helpMenu := fyne.NewMenu("帮助",
		fyne.NewMenuItem("关于 Codux", func() {}),
		fyne.NewMenuItem("检查更新", func() {}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("文档", func() {}),
	)

	return fyne.NewMainMenu(
		fileMenu,
		editMenu,
		viewMenu,
		helpMenu,
	)
}
