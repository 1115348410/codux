using Microsoft.UI.Xaml;
using Microsoft.UI.Xaml.Controls;
using Codux.WinUI.ViewModels;

namespace Codux.WinUI.Views;

public partial class WorkspaceView : Page
{
    public WorkspaceViewModel ViewModel { get; }

    public WorkspaceView()
    {
        this.InitializeComponent();
        ViewModel = App.Services.GetRequiredService<WorkspaceViewModel>();
        this.DataContext = ViewModel;

        ViewModel.PropertyChanged += OnViewModelPropertyChanged;
        UpdateEmptyState();
    }

    private void OnViewModelPropertyChanged(object? sender, System.ComponentModel.PropertyChangedEventArgs e)
    {
        if (e.PropertyName == nameof(ViewModel.SelectedProject) ||
            e.PropertyName == nameof(ViewModel.Panes))
        {
            UpdateEmptyState();
        }
    }

    private void UpdateEmptyState()
    {
        EmptyState.Visibility = ViewModel.SelectedProject == null
            ? Visibility.Visible
            : Visibility.Collapsed;
    }

    private void OnAddProjectClick(object sender, RoutedEventArgs e)
    {
        ViewModel.AddProjectCommand.Execute(null);
    }

    private void OnRemoveProjectClick(object sender, RoutedEventArgs e)
    {
        if (sender is Button button && button.Tag is Models.ProjectInfo project)
        {
            ViewModel.RemoveProjectCommand.Execute(project);
        }
    }

    private void OnNewTerminalClick(object sender, RoutedEventArgs e)
    {
        ViewModel.NewTerminalPaneCommand.Execute(null);
    }

    private void OnClosePaneClick(object sender, RoutedEventArgs e)
    {
        if (sender is Button button && button.Tag is TerminalPaneViewModel pane)
        {
            ViewModel.ClosePaneCommand.Execute(pane);
        }
    }
}
