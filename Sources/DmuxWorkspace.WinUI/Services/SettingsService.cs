using Serilog;
using System.IO;
using Newtonsoft.Json;

namespace Codux.WinUI.Services;

public class SettingsService : ISettingsService
{
    private readonly string _settingsPath;
    private Dictionary<string, object> _settings = new();

    public SettingsService()
    {
        _settingsPath = Path.Combine(
            Environment.GetFolderPath(Environment.SpecialFolder.LocalApplicationData),
            "Codux", "settings.json");
        LoadSettings();
    }

    private void LoadSettings()
    {
        try
        {
            if (File.Exists(_settingsPath))
            {
                var json = File.ReadAllText(_settingsPath);
                _settings = JsonConvert.DeserializeObject<Dictionary<string, object>>(json) ?? new();
            }
        }
        catch (Exception ex)
        {
            Log.Warning(ex, "Failed to load settings");
        }
    }

    public T GetSetting<T>(string key, T defaultValue)
    {
        if (_settings.TryGetValue(key, out var value))
        {
            try
            {
                return JsonConvert.DeserializeObject<T>(value?.ToString() ?? "") ?? defaultValue;
            }
            catch { return defaultValue; }
        }
        return defaultValue;
    }

    public void SetSetting<T>(string key, T value) => _settings[key] = value!;

    public void Save()
    {
        try
        {
            var dir = Path.GetDirectoryName(_settingsPath);
            if (dir != null) Directory.CreateDirectory(dir);
            var json = JsonConvert.SerializeObject(_settings, Formatting.Indented);
            File.WriteAllText(_settingsPath, json);
        }
        catch (Exception ex)
        {
            Log.Warning(ex, "Failed to save settings");
        }
    }
}
