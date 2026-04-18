using CommunityToolkit.Mvvm.ComponentModel;

namespace Codux.WinUI.ViewModels;

public partial class SettingsViewModel : ObservableObject
{
    [ObservableProperty]
    private string _terminalShell = "cmd.exe";

    [ObservableProperty]
    private string _terminalFontFamily = "Cascadia Code";

    [ObservableProperty]
    private int _terminalFontSize = 12;

    [ObservableProperty]
    private bool _darkMode = true;

    [ObservableProperty]
    private bool _startMinimized;

    [ObservableProperty]
    private bool _checkForUpdates = true;

    public SettingsViewModel()
    {
        LoadSettings();
    }

    private void LoadSettings()
    {
        var settings = App.Services.GetRequiredService<Codux.WinUI.Services.ISettingsService>();
        
        TerminalShell = settings.GetSetting("TerminalShell", TerminalShell);
        TerminalFontFamily = settings.GetSetting("TerminalFontFamily", TerminalFontFamily);
        TerminalFontSize = settings.GetSetting("TerminalFontSize", TerminalFontSize);
        DarkMode = settings.GetSetting("DarkMode", DarkMode);
        StartMinimized = settings.GetSetting("StartMinimized", StartMinimized);
        CheckForUpdates = settings.GetSetting("CheckForUpdates", CheckForUpdates);
    }

    public void SaveSettings()
    {
        var settings = App.Services.GetRequiredService<Codux.WinUI.Services.ISettingsService>();
        
        settings.SetSetting("TerminalShell", TerminalShell);
        settings.SetSetting("TerminalFontFamily", TerminalFontFamily);
        settings.SetSetting("TerminalFontSize", TerminalFontSize);
        settings.SetSetting("DarkMode", DarkMode);
        settings.SetSetting("StartMinimized", StartMinimized);
        settings.SetSetting("CheckForUpdates", CheckForUpdates);
        
        settings.Save();
    }
}
