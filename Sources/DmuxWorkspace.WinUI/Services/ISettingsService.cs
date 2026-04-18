namespace Codux.WinUI.Services;

public interface ISettingsService
{
    T GetSetting<T>(string key, T defaultValue);
    void SetSetting<T>(string key, T value);
    void Save();
}
