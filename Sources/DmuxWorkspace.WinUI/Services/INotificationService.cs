namespace Codux.WinUI.Services;

public interface INotificationService
{
    Task ShowToastAsync(string title, string message, CancellationToken ct = default);
    Task ShowProgressAsync(string title, double progress, CancellationToken ct = default);
}
