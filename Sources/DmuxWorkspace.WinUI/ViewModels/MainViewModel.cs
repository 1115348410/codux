using CommunityToolkit.Mvvm.ComponentModel;
using CommunityToolkit.Mvvm.Input;
using Codux.WinUI.Services;
using Microsoft.Extensions.DependencyInjection;

namespace Codux.WinUI.ViewModels;

public partial class MainViewModel : ObservableObject
{
    [ObservableProperty]
    private string _statusText = "Ready";

    [ObservableProperty]
    private string _activeProjectName = string.Empty;

    [ObservableProperty]
    private bool _isGitPanelVisible;

    [ObservableProperty]
    private bool _isAIStatsPanelVisible;

    private readonly IProjectService _projectService;
    private readonly INotificationService _notificationService;

    public MainViewModel()
    {
        _projectService = App.Services.GetRequiredService<IProjectService>();
        _notificationService = App.Services.GetRequiredService<INotificationService>();
    }

    [RelayCommand]
    private void ToggleGitPanel()
    {
        IsGitPanelVisible = !IsGitPanelVisible;
    }

    [RelayCommand]
    private void ToggleAIStatsPanel()
    {
        IsAIStatsPanelVisible = !IsAIStatsPanelVisible;
    }
}
