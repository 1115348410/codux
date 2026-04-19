package resources

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed icons/logo.png
var logo []byte

// Logo 获取 Logo 图标
func Logo() fyne.Resource {
	return &fyne.StaticResource{
		StaticName:    "logo.png",
		StaticContent: logo,
	}
}
