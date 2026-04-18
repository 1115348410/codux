using Microsoft.UI.Xaml;
using Microsoft.UI.Xaml.Controls;
using Codux.WinUI.Services;
using Codux.WinUI.ViewModels;

namespace Codux.WinUI.Views;

public partial class GitPanelView : Page
{
    public GitPanelViewModel ViewModel { get; }

    public GitPanelView()
    {
        this.InitializeComponent();
        ViewModel = App.Services.GetRequiredService<GitPanelViewModel>();
        this.DataContext = ViewModel;
    }

    public async Task LoadRepositoryAsync(string repoPath)
    {
        await ViewModel.LoadRepositoryAsync(repoPath);
    }

    private void OnStageClick(object sender, RoutedEventArgs e)
    {
        if (sender is Button button && button.Tag is GitChange change)
        {
            ViewModel.StageCommand.Execute(change);
        }
    }

    private void OnUnstageClick(object sender, RoutedEventArgs e)
    {
        if (sender is Button button && button.Tag is GitChange change)
        {
            ViewModel.UnstageCommand.Execute(change);
        }
    }
}
