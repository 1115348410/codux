using Microsoft.Toolkit.Uwp.Notifications;
using Serilog;

namespace Codux.WinUI.Services;

public class ToastNotificationService : INotificationService
{
    public Task ShowToastAsync(string title, string message, CancellationToken ct = default)
    {
        try
        {
            new ToastContentBuilder()
                .AddText(title)
                .AddText(message)
                .Show();
        }
        catch (Exception ex)
        {
            Log.Warning(ex, "Failed to show toast notification");
        }
        return Task.CompletedTask;
    }

    public Task ShowProgressAsync(string title, double progress, CancellationToken ct = default)
    {
        try
        {
            new ToastContentBuilder()
                .AddText(title)
                .AddVisualChild(new AdaptiveProgressBar
                {
                    Value = new BindableProgressBarValue("progress"),
                    Status = new BindableString("progressStatus")
                })
                .Show(toast =>
                {
                    toast.Tag = "progress";
                    toast.Data = new Windows.UI.Notifications.NotificationData();
                    toast.Data.Values["progress"] = progress.ToString();
                    toast.Data.Values["progressStatus"] = $"{progress * 100:F0}%";
                });
        }
        catch (Exception ex)
        {
            Log.Warning(ex, "Failed to show progress notification");
        }
        return Task.CompletedTask;
    }
}
