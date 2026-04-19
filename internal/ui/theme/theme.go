package theme

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// ThemeMode 主题模式
type ThemeMode string

const (
	ThemeModeSystem ThemeMode = "system"
	ThemeModeLight  ThemeMode = "light"
	ThemeModeDark   ThemeMode = "dark"
)

// CoduxTheme 自定义主题
type CoduxTheme struct {
	mode ThemeMode
}

// NewCoduxTheme 创建 Codux 主题
func NewCoduxTheme() *CoduxTheme {
	return &CoduxTheme{
		mode: ThemeModeSystem,
	}
}

// Color 实现 fyne.Theme 接口
func (t *CoduxTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if t.mode == ThemeModeDark || (t.mode == ThemeModeSystem && variant == theme.VariantDark) {
		return t.darkColor(name)
	}
	return t.lightColor(name)
}

func (t *CoduxTheme) darkColor(name fyne.ThemeColorName) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 30, G: 30, B: 30, A: 255}
	case theme.ColorNameButton:
		return color.NRGBA{R: 50, G: 50, B: 50, A: 255}
	case theme.ColorNameDisabled:
		return color.NRGBA{R: 80, G: 80, B: 80, A: 255}
	case theme.ColorNameDisabledButton:
		return color.NRGBA{R: 40, G: 40, B: 40, A: 255}
	case theme.ColorNameError:
		return color.NRGBA{R: 239, G: 98, B: 97, A: 255}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 235, G: 235, B: 235, A: 255}
	case theme.ColorNameHover:
		return color.NRGBA{R: 60, G: 60, B: 60, A: 255}
	case theme.ColorNameHyperlink:
		return color.NRGBA{R: 100, G: 149, B: 237, A: 255}
	case theme.ColorNameInputBackground:
		return color.NRGBA{R: 40, G: 40, B: 40, A: 255}
	case theme.ColorNameInputBorder:
		return color.NRGBA{R: 60, G: 60, B: 60, A: 255}
	case theme.ColorNameOverlayBackground:
		return color.NRGBA{R: 45, G: 45, B: 45, A: 255}
	case theme.ColorNamePlaceHolder:
		return color.NRGBA{R: 120, G: 120, B: 120, A: 255}
	case theme.ColorNamePressed:
		return color.NRGBA{R: 70, G: 70, B: 70, A: 255}
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 0, G: 122, B: 255, A: 255}
	case theme.ColorNameScrollBar:
		return color.NRGBA{R: 70, G: 70, B: 70, A: 200}
	case theme.ColorNameSeparator:
		return color.NRGBA{R: 50, G: 50, B: 50, A: 255}
	case theme.ColorNameShadow:
		return color.NRGBA{R: 0, G: 0, B: 0, A: 100}
	default:
		return theme.DefaultTheme().Color(name, theme.VariantDark)
	}
}

func (t *CoduxTheme) lightColor(name fyne.ThemeColorName) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	case theme.ColorNameButton:
		return color.NRGBA{R: 235, G: 235, B: 235, A: 255}
	case theme.ColorNameDisabled:
		return color.NRGBA{R: 200, G: 200, B: 200, A: 255}
	case theme.ColorNameDisabledButton:
		return color.NRGBA{R: 220, G: 220, B: 220, A: 255}
	case theme.ColorNameError:
		return color.NRGBA{R: 255, G: 59, B: 48, A: 255}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 30, G: 30, B: 30, A: 255}
	case theme.ColorNameHover:
		return color.NRGBA{R: 240, G: 240, B: 240, A: 255}
	case theme.ColorNameHyperlink:
		return color.NRGBA{R: 0, G: 122, B: 255, A: 255}
	case theme.ColorNameInputBackground:
		return color.NRGBA{R: 245, G: 245, B: 245, A: 255}
	case theme.ColorNameInputBorder:
		return color.NRGBA{R: 210, G: 210, B: 210, A: 255}
	case theme.ColorNameOverlayBackground:
		return color.NRGBA{R: 255, G: 255, B: 255, A: 250}
	case theme.ColorNamePlaceHolder:
		return color.NRGBA{R: 150, G: 150, B: 150, A: 255}
	case theme.ColorNamePressed:
		return color.NRGBA{R: 220, G: 220, B: 220, A: 255}
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 0, G: 122, B: 255, A: 255}
	case theme.ColorNameScrollBar:
		return color.NRGBA{R: 180, G: 180, B: 180, A: 200}
	case theme.ColorNameSeparator:
		return color.NRGBA{R: 220, G: 220, B: 220, A: 255}
	case theme.ColorNameShadow:
		return color.NRGBA{R: 0, G: 0, B: 0, A: 50}
	default:
		return theme.DefaultTheme().Color(name, theme.VariantLight)
	}
}

// Font 实现 fyne.Theme 接口
func (t *CoduxTheme) Font(style fyne.TextStyle) fyne.Resource {
	return nil
}

// Icon 实现 fyne.Theme 接口
func (t *CoduxTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size 实现 fyne.Theme 接口
func (t *CoduxTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

// SetMode 设置主题模式
func (t *CoduxTheme) SetMode(mode ThemeMode) {
	t.mode = mode
}

// Mode 获取当前主题模式
func (t *CoduxTheme) Mode() ThemeMode {
	return t.mode
}
