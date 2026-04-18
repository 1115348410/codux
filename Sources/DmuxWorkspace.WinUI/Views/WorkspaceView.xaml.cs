using Microsoft.UI.Xaml;
using Microsoft.UI.Xaml.Controls;
using Microsoft.UI.Xaml.Media;
using Codux.WinUI.ViewModels;

namespace Codux.WinUI.Views;

public partial class WorkspaceView : Page
{
    public WorkspaceViewModel ViewModel { get; }

    private readonly Dictionary<Guid, TerminalView> _terminalViews = new();

    public WorkspaceView()
    {
        this.InitializeComponent();
        ViewModel = App.Services.GetRequiredService<WorkspaceViewModel>();
        this.DataContext = ViewModel;

        ViewModel.PropertyChanged += OnViewModelPropertyChanged;
        ViewModel.Panes.CollectionChanged += OnPanesCollectionChanged;

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

    private void OnPanesCollectionChanged(object? sender, System.Collections.Specialized.NotifyCollectionChangedEventArgs e)
    {
        if (e.Action == System.Collections.Specialized.NotifyCollectionChangedAction.Add)
        {
            foreach (TerminalPaneViewModel pane in e.NewItems!)
            {
                AddTerminalView(pane);
            }
        }
        else if (e.Action == System.Collections.Specialized.NotifyCollectionChangedAction.Remove)
        {
            foreach (TerminalPaneViewModel pane in e.OldItems!)
            {
                RemoveTerminalView(pane);
            }
        }
        UpdateEmptyState();
    }

    private void AddTerminalView(TerminalPaneViewModel pane)
    {
        var terminalView = new TerminalView(pane);
        terminalView.CloseRequested += OnTerminalCloseRequested;

        _terminalViews[pane.Id] = terminalView;

        if (SplitContainer.Children.Count == 0)
        {
            SplitContainer.Children.Add(terminalView);
        }
        else
        {
            var firstView = SplitContainer.Children[0] as TerminalView;
            if (firstView != null)
            {
                var grid = new Grid();
                grid.Children.Add(firstView);
                grid.Children.Add(terminalView);
                SplitContainer.Children.Clear();
                SplitContainer.Children.Add(grid);
            }
        }
    }

    private void RemoveTerminalView(TerminalPaneViewModel pane)
    {
        if (_terminalViews.TryGetValue(pane.Id, out var view))
        {
            view.CloseRequested -= OnTerminalCloseRequested;
            SplitContainer.Children.Remove(view);
            _terminalViews.Remove(pane.Id);
        }
        UpdateEmptyState();
    }

    private void OnTerminalCloseRequested(object? sender, EventArgs e)
    {
        if (sender is TerminalView view)
        {
            ViewModel.ClosePaneCommand.Execute(view.ViewModel);
        }
    }

    private void UpdateEmptyState()
    {
        EmptyState.Visibility = ViewModel.Panes.Count == 0
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

    private void OnProjectSelectionChanged(object sender, SelectionChangedEventArgs e)
    {
        if (ViewModel.SelectedProject != null && ViewModel.Panes.Count == 0)
        {
            ViewModel.NewTerminalPaneCommand.Execute(null);
        }
    }

    private void OnSplitHorizontalClick(object sender, RoutedEventArgs e)
    {
        if (ViewModel.SelectedProject == null) return;
        ViewModel.NewTerminalPaneCommand.Execute(null);
    }

    private void OnSplitVerticalClick(object sender, RoutedEventArgs e)
    {
        if (ViewModel.SelectedProject == null) return;
        ViewModel.NewTerminalPaneCommand.Execute(null);
    }
}
