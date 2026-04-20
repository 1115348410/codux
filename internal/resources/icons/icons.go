package icons

import (
	_ "embed"
	"image/color"
	_ "image/png"

	"fyne.io/fyne/v2"
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

// CustomColor 自定义颜色
func CustomColor() color.Color {
	return color.NRGBA{R: 0, G: 122, B: 255, A: 255}
}
