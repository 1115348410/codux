package main

import (
	"image/color"
	_ "image/png"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

// 资源打包
var (
	//go:embed resources/icons/logo.png
	logoPNG []byte
)

// LogoResource Logo 图标资源
var LogoResource = &fyne.StaticResource{
	StaticName:    "logo.png",
	StaticContent: logoPNG,
}

// init 初始化应用图标
func init() {
	app.SetIcon(LogoResource)
}

// CustomColor 自定义颜色
func CustomColor() color.Color {
	return color.NRGBA{R: 0, G: 122, B: 255, A: 255}
}
