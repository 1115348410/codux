package settings

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/duxweb/codux/internal/app"
)

// SettingsWindow 设置窗口
type SettingsWindow struct {
	window fyne.Window
	store  *app.Store
}

// NewSettingsWindow 创建设置窗口
func NewSettingsWindow(parent fyne.Window, store *app.Store) *SettingsWindow {
	return &SettingsWindow{
		window: parent,
		store:  store,
	}
}

// Show 显示设置窗口
func (s *SettingsWindow) Show() {
	w := fyne.CurrentApp().NewWindow("Codux 设置")
	w.Resize(fyne.NewSize(600, 400))

	// 创建设置内容
	content := s.createSettingsContent()
	w.SetContent(content)
	w.Show()
}

func (s *SettingsWindow) createSettingsContent() fyne.CanvasObject {
	// General 设置
	generalTab := s.createGeneralTab()

	// Terminal 设置
	terminalTab := s.createTerminalTab()

	// Git 设置
	gitTab := s.createGitTab()

	// AI 设置
	aiTab := s.createAITab()

	tabs := container.NewAppTabs(
		container.NewTabItem("通用", generalTab),
		container.NewTabItem("终端", terminalTab),
		container.NewTabItem("Git", gitTab),
		container.NewTabItem("AI", aiTab),
	)

	return container.NewMax(tabs)
}

func (s *SettingsWindow) createGeneralTab() fyne.CanvasObject {
	settings := s.store.Settings()

	themeSelect := widget.NewSelect([]string{"明亮", "暗黑", "跟随系统"}, func(value string) {
		switch value {
		case "明亮":
			settings.Theme = "light"
		case "暗黑":
			settings.Theme = "dark"
		default:
			settings.Theme = "system"
		}
		s.store.UpdateSettings(settings)
	})

	switch settings.Theme {
	case "light":
		themeSelect.SetSelected("明亮")
	case "dark":
		themeSelect.SetSelected("暗黑")
	default:
		themeSelect.SetSelected("跟随系统")
	}

	return container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("主题", themeSelect),
		),
		widget.NewSeparator(),
		widget.NewLabel("其他设置将在这里添加"),
	)
}

func (s *SettingsWindow) createTerminalTab() fyne.CanvasObject {
	settings := s.store.Settings()

	// 字体大小
	fontSizeEntry := widget.NewEntry()
	fontSizeEntry.SetText("14")

	fontSizeSlider := widget.NewSlider(10, 24)
	fontSizeSlider.Value = float64(settings.TerminalFontSize)
	fontSizeSlider.OnChanged = func(value float64) {
		settings.TerminalFontSize = value
		s.store.UpdateSettings(settings)
		fontSizeEntry.SetText(string(rune(int(value)+'0') + '0'))
	}

	// 终端快捷键
	shortcutKeyEntry := widget.NewEntry()
	shortcutKeyEntry.SetPlaceHolder("按组合键设置...")
	shortcutKeyEntry.Disable()

	clearHistoryBtn := widget.NewButton("清空终端历史", func() {
		// TODO: 清空历史
	})

	return container.NewVBox(
		widget.NewLabel("外观"),
		widget.NewForm(
			widget.NewFormItem("字体大小", container.NewHBox(fontSizeSlider, fontSizeEntry)),
		),
		widget.NewSeparator(),
		widget.NewLabel("快捷键"),
		widget.NewForm(
			widget.NewFormItem("搜索终端", shortcutKeyEntry),
		),
		widget.NewSeparator(),
		widget.NewLabel("其他"),
		container.NewHBox(clearHistoryBtn),
	)
}

func (s *SettingsWindow) createGitTab() fyne.CanvasObject {
	return container.NewVBox(
		widget.NewLabel("Git 设置将在这里添加"),
		widget.NewForm(
			widget.NewFormItem("自动刷新间隔", widget.NewSelect([]string{"5 秒", "10 秒", "30 秒", "60 秒"}, func(string) {})),
		),
	)
}

func (s *SettingsWindow) createAITab() fyne.CanvasObject {
	return container.NewVBox(
		widget.NewLabel("AI 使用统计设置"),
		widget.NewCheck("后台自动刷新", func(bool) {}),
		widget.NewForm(
			widget.NewFormItem("刷新间隔", widget.NewSelect([]string{"5 秒", "10 秒", "30 秒"}, func(string) {})),
		),
	)
}
