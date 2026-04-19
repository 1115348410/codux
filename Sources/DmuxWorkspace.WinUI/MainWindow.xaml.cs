using System.Windows;
using System.Windows.Controls;
using System.Windows.Media;
using Codux.WinUI.Services;
using Codux.WinUI.ViewModels;
using Microsoft.Extensions.DependencyInjection;

namespace Codux.WinUI;

public partial class MainWindow : Window
{
    private readonly IProjectService _projectService;
    private readonly ITerminalService _terminalService;
    private ProjectInfo? _selectedProject;
    private readonly List<TerminalPaneViewModel> _terminals = new();

    public MainWindow()
    {
        InitializeComponent();

        _projectService = App.Services.GetRequiredService<IProjectService>();
        _terminalService = App.Services.GetRequiredService<ITerminalService>();

        _ = LoadProjectsAsync();
    }

    private async Task LoadProjectsAsync()
    {
        var projects = await _projectService.GetProjectsAsync();
        var projectList = projects.ToList();
        ProjectsList.ItemsSource = projectList;
        if (projectList.Any())
        {
            ProjectsList.SelectedItem = projectList.First();
            _selectedProject = projectList.First();
            ProjectPathText.Text = _selectedProject.Path;
        }
    }

    private void OnProjectSelected(object sender, SelectionChangedEventArgs e)
    {
        _selectedProject = ProjectsList.SelectedItem as ProjectInfo;
        if (_selectedProject != null)
        {
            ProjectPathText.Text = _selectedProject.Path;
        }
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
        _terminals.Add(terminal);

        var panel = CreateTerminalPanel(terminal);
        TerminalsContainer.Children.Add(panel);

        StatusText.Text = $"Terminals: {_terminals.Count} | {_selectedProject.Path}";
    }

    private void OnAddProject(object sender, RoutedEventArgs e)
    {
        var dialog = new Microsoft.Win32.OpenFolderDialog
        {
            Title = "Select Project Folder"
        };

        if (dialog.ShowDialog() == true)
        {
            try
            {
                var project = _projectService.AddProjectAsync(dialog.FolderName).Result;
                var projectList = ((List<ProjectInfo>)ProjectsList.ItemsSource);
                projectList.Add(project);
                ProjectsList.Items.Refresh();
                ProjectsList.SelectedItem = project;
                _selectedProject = project;
                ProjectPathText.Text = project.Path;
                StatusText.Text = $"Added project: {project.Name}";
            }
            catch (Exception ex)
            {
                MessageBox.Show($"Failed to add project: {ex.Message}", "Error",
                    MessageBoxButton.OK, MessageBoxImage.Error);
            }
        }
    }

    private Border CreateTerminalPanel(TerminalPaneViewModel terminal)
    {
        var border = new Border
        {
            Background = new SolidColorBrush(Color.FromRgb(0x0C, 0x0C, 0x0C)),
            Margin = new Thickness(8),
            Padding = new Thickness(0),
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
        if (e.Key == System.Windows.Input.Key.Enter)
        {
            if (sender is TextBox textBox && textBox.Tag is TerminalPaneViewModel terminal)
            {
                var cmd = textBox.Text;
                textBox.Text = string.Empty;
                await terminal.SendInputAsync(cmd);
            }
            e.Handled = true;
        }
    }

    private void OnCloseTerminal(object sender, RoutedEventArgs e)
    {
        if (sender is Button button && button.Tag is TerminalPaneViewModel terminal)
        {
            _terminals.Remove(terminal);

            foreach (var child in TerminalsContainer.Children.OfType<Border>())
            {
                if (child.Tag == terminal)
                {
                    TerminalsContainer.Children.Remove(child);
                    break;
                }
            }

            StatusText.Text = $"Terminals: {_terminals.Count}";
        }
    }
}
