using System.Windows;
using System.Windows.Controls;
using System.Windows.Media;
using System.Collections.ObjectModel;
using System.IO;
using Codux.WinUI.Services;
using Codux.WinUI.ViewModels;
using Microsoft.Extensions.DependencyInjection;

namespace Codux.WinUI;

public partial class MainWindow : Window
{
    private readonly IProjectService _projectService;
    private readonly ITerminalService _terminalService;
    private readonly IGitService _gitService;
    private readonly IUsageService _usageService;
    private ProjectInfo? _selectedProject;
    private readonly ObservableCollection<ProjectInfo> _projects = new();
    private readonly List<TerminalPaneViewModel> _terminals = new();
    private readonly Dictionary<TerminalPaneViewModel, List<string>> _terminalInputHistory = new();
    private readonly Dictionary<TerminalPaneViewModel, int> _terminalInputHistoryIndex = new();
    private bool _isProjectSelectionInternalUpdate;

    public MainWindow()
    {
        InitializeComponent();
        Closing += OnWindowClosing;

        _projectService = App.Services.GetRequiredService<IProjectService>();
        _terminalService = App.Services.GetRequiredService<ITerminalService>();
        _gitService = App.Services.GetRequiredService<IGitService>();
        _usageService = App.Services.GetRequiredService<IUsageService>();

        _ = LoadProjectsAsync();
    }

    private async Task LoadProjectsAsync()
    {
        var projects = await _projectService.GetProjectsAsync();
        _projects.Clear();

        foreach (var project in projects)
        {
            _projects.Add(project);
        }

        ProjectsList.ItemsSource = _projects;
        if (_projects.Any())
        {
            ProjectsList.SelectedItem = _projects.First();
            _selectedProject = _projects.First();
            ProjectPathText.Text = _selectedProject.Path;
        }

        await RefreshInspectorAsync();
    }

    private async void OnProjectSelected(object sender, SelectionChangedEventArgs e)
    {
        if (_isProjectSelectionInternalUpdate)
        {
            return;
        }

        _selectedProject = ProjectsList.SelectedItem as ProjectInfo;
        if (_selectedProject != null)
        {
            ProjectPathText.Text = _selectedProject.Path;
            await _projectService.TouchProjectAsync(_selectedProject.Id);
            UpdateProjectLastAccessed(_selectedProject.Id);
        }

        _ = RefreshInspectorAsync();
    }

    private async void OnNewTerminal(object sender, RoutedEventArgs e)
    {
        if (_selectedProject == null)
        {
            MessageBox.Show("Please select or add a project first.", "No Project Selected",
                MessageBoxButton.OK, MessageBoxImage.Information);
            return;
        }

        var terminal = new TerminalPaneViewModel(_terminalService)
        {
            WorkingDirectory = _selectedProject.Path
        };

        await terminal.CreateSessionAsync();
        terminal.AppendClientOutput($"Session started in {_selectedProject.Path}");
        _terminals.Add(terminal);
        _terminalInputHistory[terminal] = new List<string>();
        _terminalInputHistoryIndex[terminal] = 0;

        var panel = CreateTerminalPanel(terminal);
        TerminalsContainer.Children.Add(panel);

        StatusText.Text = $"Terminals: {_terminals.Count} | {_selectedProject.Path}";
        await RefreshInspectorAsync();
    }

    private async void OnAddProject(object sender, RoutedEventArgs e)
    {
        var dialog = new Microsoft.Win32.OpenFolderDialog
        {
            Title = "Select Project Folder"
        };

        if (dialog.ShowDialog() == true)
        {
            try
            {
                var project = await _projectService.AddProjectAsync(dialog.FolderName);
                _projects.Add(project);
                ProjectsList.SelectedItem = project;
                _selectedProject = project;
                ProjectPathText.Text = project.Path;
                StatusText.Text = $"Added project: {project.Name}";
                await RefreshInspectorAsync();
            }
            catch (Exception ex)
            {
                MessageBox.Show($"Failed to add project: {ex.Message}", "Error",
                    MessageBoxButton.OK, MessageBoxImage.Error);
            }
        }
    }

    private async void OnRemoveProject(object sender, RoutedEventArgs e)
    {
        if (_selectedProject == null)
        {
            return;
        }

        var toRemove = _selectedProject;
        await _projectService.RemoveProjectAsync(toRemove.Id);

        _projects.Remove(toRemove);
        _selectedProject = _projects.FirstOrDefault();
        ProjectsList.SelectedItem = _selectedProject;
        ProjectPathText.Text = _selectedProject?.Path ?? string.Empty;
        StatusText.Text = $"Removed project: {toRemove.Name}";
        await RefreshInspectorAsync();
    }

    private Border CreateTerminalPanel(TerminalPaneViewModel terminal)
    {
        var border = new Border
        {
            Background = new SolidColorBrush(Color.FromRgb(0x0C, 0x0C, 0x0C)),
            Margin = new Thickness(8),
            Padding = new Thickness(0),
            BorderBrush = new SolidColorBrush(Color.FromRgb(0x2F, 0x2F, 0x2F)),
            BorderThickness = new Thickness(1),
            Tag = terminal
        };

        var grid = new Grid();
        grid.RowDefinitions.Add(new RowDefinition { Height = GridLength.Auto });
        grid.RowDefinitions.Add(new RowDefinition { Height = new GridLength(1, GridUnitType.Star) });
        grid.RowDefinitions.Add(new RowDefinition { Height = GridLength.Auto });

        // Header with close button
        var header = new Border
        {
            Background = new SolidColorBrush(Color.FromRgb(0x1E, 0x1E, 0x1E)),
            Padding = new Thickness(8, 4, 4, 4)
        };
        var headerGrid = new Grid();
        headerGrid.ColumnDefinitions.Add(new ColumnDefinition { Width = new GridLength(1, GridUnitType.Star) });
        headerGrid.ColumnDefinitions.Add(new ColumnDefinition { Width = GridLength.Auto });

        var headerText = new TextBlock
        {
            Text = terminal.WorkingDirectory,
            Foreground = new SolidColorBrush(Color.FromRgb(0xCC, 0xCC, 0xCC)),
            FontSize = 11,
            VerticalAlignment = VerticalAlignment.Center
        };
        Grid.SetColumn(headerText, 0);

        var closeButton = new Button
        {
            Content = "X",
            Width = 20,
            Height = 20,
            FontSize = 10,
            Background = Brushes.Transparent,
            Foreground = new SolidColorBrush(Color.FromRgb(0x99, 0x99, 0x99)),
            BorderThickness = new Thickness(0),
            Cursor = System.Windows.Input.Cursors.Hand,
            Tag = terminal
        };
        closeButton.Click += OnCloseTerminal;
        Grid.SetColumn(closeButton, 1);

        headerGrid.Children.Add(headerText);
        headerGrid.Children.Add(closeButton);
        header.Child = headerGrid;
        Grid.SetRow(header, 0);

        // Output area - using TextBox for selectable text
        var outputScrollViewer = new ScrollViewer
        {
            VerticalScrollBarVisibility = ScrollBarVisibility.Auto,
            HorizontalScrollBarVisibility = ScrollBarVisibility.Auto,
            Background = Brushes.Transparent
        };

        var outputTextBox = new TextBox
        {
            Text = terminal.TerminalOutput,
            FontFamily = new FontFamily("Cascadia Code, Consolas, Courier New"),
            FontSize = 13,
            Foreground = new SolidColorBrush(Color.FromRgb(0xCC, 0xCC, 0xCC)),
            Background = Brushes.Transparent,
            BorderThickness = new Thickness(0),
            IsReadOnly = true,
            TextWrapping = TextWrapping.Wrap,
            AcceptsReturn = true,
            VerticalScrollBarVisibility = ScrollBarVisibility.Disabled,
            HorizontalScrollBarVisibility = ScrollBarVisibility.Disabled,
            Margin = new Thickness(8, 4, 8, 4),
            Padding = new Thickness(0)
        };

        terminal.PropertyChanged += (s, e) =>
        {
            if (e.PropertyName == nameof(TerminalPaneViewModel.TerminalOutput))
            {
                Dispatcher.Invoke(() =>
                {
                    outputTextBox.Text = terminal.TerminalOutput;
                    outputScrollViewer.ScrollToEnd();
                });
            }
        };

        outputScrollViewer.Content = outputTextBox;
        Grid.SetRow(outputScrollViewer, 1);

        // Input area
        var inputBorder = new Border
        {
            Background = new SolidColorBrush(Color.FromRgb(0x1E, 0x1E, 0x1E)),
            Padding = new Thickness(0)
        };
        var inputGrid = new Grid();
        inputGrid.ColumnDefinitions.Add(new ColumnDefinition { Width = GridLength.Auto });
        inputGrid.ColumnDefinitions.Add(new ColumnDefinition { Width = new GridLength(1, GridUnitType.Star) });

        var promptText = new TextBlock
        {
            Text = "> ",
            Foreground = new SolidColorBrush(Color.FromRgb(0xCC, 0xCC, 0xCC)),
            FontFamily = new FontFamily("Cascadia Code, Consolas, Courier New"),
            FontSize = 13,
            Margin = new Thickness(8, 6, 0, 6),
            VerticalAlignment = VerticalAlignment.Center
        };
        Grid.SetColumn(promptText, 0);

        var inputBox = new TextBox
        {
            Background = Brushes.Transparent,
            Foreground = Brushes.White,
            FontFamily = new FontFamily("Cascadia Code, Consolas, Courier New"),
            FontSize = 13,
            BorderThickness = new Thickness(0),
            Margin = new Thickness(0, 6, 8, 6),
            VerticalAlignment = VerticalAlignment.Center,
            Tag = terminal
        };
        inputBox.KeyDown += OnInputKeyDown;
        Grid.SetColumn(inputBox, 1);

        inputGrid.Children.Add(promptText);
        inputGrid.Children.Add(inputBox);
        inputBorder.Child = inputGrid;
        Grid.SetRow(inputBorder, 2);

        grid.Children.Add(header);
        grid.Children.Add(outputScrollViewer);
        grid.Children.Add(inputBorder);
        border.Child = grid;

        return border;
    }

    private async void OnInputKeyDown(object sender, System.Windows.Input.KeyEventArgs e)
    {
        if (sender is not TextBox textBox || textBox.Tag is not TerminalPaneViewModel terminal)
        {
            return;
        }

        if ((System.Windows.Input.Keyboard.Modifiers & System.Windows.Input.ModifierKeys.Control) == System.Windows.Input.ModifierKeys.Control &&
            e.Key == System.Windows.Input.Key.C)
        {
            terminal.AppendClientOutput("^C");
            await terminal.SendInputAsync("\u0003");
            e.Handled = true;
            return;
        }

        if (e.Key == System.Windows.Input.Key.Up)
        {
            NavigateTerminalHistory(textBox, terminal, -1);
            e.Handled = true;
            return;
        }

        if (e.Key == System.Windows.Input.Key.Down)
        {
            NavigateTerminalHistory(textBox, terminal, 1);
            e.Handled = true;
            return;
        }

        if (e.Key == System.Windows.Input.Key.Enter)
        {
            var cmd = textBox.Text;
            textBox.Text = string.Empty;

            if (!string.IsNullOrWhiteSpace(cmd))
            {
                terminal.AppendClientOutput(cmd);
                await terminal.SendInputAsync(cmd);

                if (_terminalInputHistory.TryGetValue(terminal, out var history))
                {
                    if (history.Count == 0 || history[^1] != cmd)
                    {
                        history.Add(cmd);
                    }

                    _terminalInputHistoryIndex[terminal] = history.Count;
                }
            }
            e.Handled = true;
        }
    }

    private async void OnCloseTerminal(object sender, RoutedEventArgs e)
    {
        if (sender is Button button && button.Tag is TerminalPaneViewModel terminal)
        {
            await terminal.CloseAsync();
            _terminals.Remove(terminal);
            _terminalInputHistory.Remove(terminal);
            _terminalInputHistoryIndex.Remove(terminal);

            foreach (var child in TerminalsContainer.Children.OfType<Border>())
            {
                if (child.Tag == terminal)
                {
                    TerminalsContainer.Children.Remove(child);
                    break;
                }
            }

            StatusText.Text = $"Terminals: {_terminals.Count}";
            await RefreshInspectorAsync();
        }
    }

    private async void OnGitStageAll(object sender, RoutedEventArgs e)
    {
        if (_selectedProject == null)
        {
            return;
        }

        try
        {
            await _gitService.StageAllAsync(_selectedProject.Path);
            StatusText.Text = "Git: staged all changes";
            await RefreshInspectorAsync();
        }
        catch (Exception ex)
        {
            MessageBox.Show($"Failed to stage changes: {ex.Message}", "Git Error", MessageBoxButton.OK, MessageBoxImage.Error);
        }
    }

    private async void OnGitCommit(object sender, RoutedEventArgs e)
    {
        if (_selectedProject == null)
        {
            return;
        }

        var message = GitCommitMessageTextBox.Text.Trim();
        if (string.IsNullOrWhiteSpace(message))
        {
            MessageBox.Show("Please enter a commit message.", "Git Commit", MessageBoxButton.OK, MessageBoxImage.Information);
            return;
        }

        try
        {
            await _gitService.CommitAsync(_selectedProject.Path, message);
            StatusText.Text = "Git: commit completed";
            await RefreshInspectorAsync();
        }
        catch (Exception ex)
        {
            MessageBox.Show($"Failed to commit: {ex.Message}", "Git Error", MessageBoxButton.OK, MessageBoxImage.Error);
        }
    }

    private async void OnRefreshInspector(object sender, RoutedEventArgs e)
    {
        await RefreshInspectorAsync();
    }

    private async Task RefreshInspectorAsync()
    {
        OverviewSelectedProjectText.Text = _selectedProject?.Path ?? "No project selected";
        OverviewTerminalCountText.Text = _terminals.Count.ToString();

        if (_selectedProject == null)
        {
            GitCurrentBranchText.Text = "-";
            GitChangeCountText.Text = "0";
            GitChangesList.ItemsSource = null;
            UsageTotalTokensText.Text = "0";
            return;
        }

        try
        {
            if (Directory.Exists(Path.Combine(_selectedProject.Path, ".git")))
            {
                var branches = (await _gitService.GetBranchesAsync(_selectedProject.Path)).ToList();
                var currentBranch = branches.FirstOrDefault(b => !b.IsRemote && b.IsCurrent)?.Name ?? "detached";
                GitCurrentBranchText.Text = currentBranch;

                var changes = (await _gitService.GetChangesAsync(_selectedProject.Path)).ToList();
                GitChangeCountText.Text = changes.Count.ToString();
                GitChangesList.ItemsSource = changes.Select(c => $"[{c.Status}] {c.FilePath}").ToList();
            }
            else
            {
                GitCurrentBranchText.Text = "Not a git repository";
                GitChangeCountText.Text = "0";
                GitChangesList.ItemsSource = new[] { "No repository metadata found." };
            }

            var totalTokens = await _usageService.GetTotalUsageAsync();
            UsageTotalTokensText.Text = totalTokens.ToString("N0");
        }
        catch (Exception ex)
        {
            StatusText.Text = $"Inspector refresh failed: {ex.Message}";
        }
    }

    private async void OnWindowClosing(object? sender, System.ComponentModel.CancelEventArgs e)
    {
        var terminals = _terminals.ToList();
        foreach (var terminal in terminals)
        {
            await terminal.CloseAsync();
        }
        _terminals.Clear();
        _terminalInputHistory.Clear();
        _terminalInputHistoryIndex.Clear();
    }

    private void NavigateTerminalHistory(TextBox input, TerminalPaneViewModel terminal, int delta)
    {
        if (!_terminalInputHistory.TryGetValue(terminal, out var history) || history.Count == 0)
        {
            return;
        }

        var index = _terminalInputHistoryIndex.TryGetValue(terminal, out var currentIndex)
            ? currentIndex
            : history.Count;

        index = Math.Clamp(index + delta, 0, history.Count);
        _terminalInputHistoryIndex[terminal] = index;

        input.Text = index == history.Count ? string.Empty : history[index];
        input.CaretIndex = input.Text.Length;
    }

    private void UpdateProjectLastAccessed(Guid projectId)
    {
        var existing = _projects.FirstOrDefault(p => p.Id == projectId);
        if (existing == null)
        {
            return;
        }

        var index = _projects.IndexOf(existing);
        if (index < 0)
        {
            return;
        }

        var updated = existing with { LastAccessed = DateTime.Now };
        _projects[index] = updated;
        _selectedProject = updated;

        var ordered = _projects.OrderByDescending(p => p.LastAccessed).ToList();
        _projects.Clear();
        foreach (var item in ordered)
        {
            _projects.Add(item);
        }

        _isProjectSelectionInternalUpdate = true;
        try
        {
            ProjectsList.SelectedItem = _selectedProject;
        }
        finally
        {
            _isProjectSelectionInternalUpdate = false;
        }
    }
}
