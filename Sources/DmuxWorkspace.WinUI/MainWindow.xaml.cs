using System.Windows;
using System.Windows.Controls;
using System.Collections.ObjectModel;
using Codux.WinUI.Services;
using Codux.WinUI.ViewModels;
using Microsoft.Extensions.DependencyInjection;

namespace Codux.WinUI;

public partial class MainWindow : Window
{
    private readonly IProjectService _projectService;
    private readonly ITerminalService _terminalService;
    private ProjectInfo? _selectedProject;
    private readonly ObservableCollection<ProjectInfo> _projects = new();
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
        _projects.Clear();

        foreach (var project in projects)
        {
            _projects.Add(project);
        }

        ProjectsList.ItemsSource = _projects;
        if (_projects.Any())
        {
            ProjectsList.SelectedItem = _projects.First();
        }
    }

    private void OnProjectSelected(object sender, SelectionChangedEventArgs e)
    {
        _selectedProject = ProjectsList.SelectedItem as ProjectInfo;
    }

    private async void OnNewTerminal(object sender, RoutedEventArgs e)
    {
        if (_selectedProject == null) return;

        var terminal = new TerminalPaneViewModel(_terminalService)
        {
            WorkingDirectory = _selectedProject.Path
        };

        await terminal.CreateSessionAsync();
        _terminals.Add(terminal);

        var panel = CreateTerminalPanel(terminal);
        TerminalsContainer.Children.Add(panel);

        StatusText.Text = $"Terminals: {_terminals.Count}";
    }

    private async void OnAddProject(object sender, RoutedEventArgs e)
    {
        var dialog = new Microsoft.Win32.OpenFolderDialog { Title = "Select Project Folder" };
        if (dialog.ShowDialog() == true)
        {
            var project = await _projectService.AddProjectAsync(dialog.FolderName);
            _projects.Add(project);
            ProjectsList.SelectedItem = project;
        }
    }

    private Border CreateTerminalPanel(TerminalPaneViewModel terminal)
    {
        var border = new Border
        {
            Background = new System.Windows.Media.SolidColorBrush(System.Windows.Media.Color.FromRgb(0x0C, 0x0C, 0x0C)),
            Margin = new Thickness(8),
            Padding = new Thickness(8)
        };

        var grid = new Grid();
        grid.RowDefinitions.Add(new RowDefinition { Height = GridLength.Auto });
        grid.RowDefinitions.Add(new RowDefinition { Height = new GridLength(1, GridUnitType.Star) });
        grid.RowDefinitions.Add(new RowDefinition { Height = GridLength.Auto });

        var header = new Border { Background = new System.Windows.Media.SolidColorBrush(System.Windows.Media.Color.FromRgb(0x1E, 0x1E, 0x1E)), Padding = new Thickness(8, 4, 8, 4) };
        var headerText = new TextBlock { Text = terminal.WorkingDirectory, Foreground = new System.Windows.Media.SolidColorBrush(System.Windows.Media.Color.FromRgb(0xCC, 0xCC, 0xCC)), FontSize = 11 };
        header.Child = headerText;
        Grid.SetRow(header, 0);

        var output = new ScrollViewer { VerticalScrollBarVisibility = ScrollBarVisibility.Auto };
        var outputText = new TextBlock
        {
            Text = terminal.TerminalOutput,
            FontFamily = new System.Windows.Media.FontFamily("Cascadia Code, Consolas"),
            FontSize = 13,
            Foreground = new System.Windows.Media.SolidColorBrush(System.Windows.Media.Color.FromRgb(0xCC, 0xCC, 0xCC)),
            TextWrapping = TextWrapping.Wrap
        };
        output.Content = outputText;
        terminal.PropertyChanged += (s, e) =>
        {
            if (e.PropertyName == nameof(TerminalPaneViewModel.TerminalOutput))
            {
                Dispatcher.Invoke(() =>
                {
                    outputText.Text = terminal.TerminalOutput;
                    output.ScrollToEnd();
                });
            }
        };
        Grid.SetRow(output, 1);

        var inputPanel = new Border { Background = new System.Windows.Media.SolidColorBrush(System.Windows.Media.Color.FromRgb(0x1E, 0x1E, 0x1E)) };
        var inputGrid = new Grid();
        inputGrid.ColumnDefinitions.Add(new ColumnDefinition { Width = GridLength.Auto });
        inputGrid.ColumnDefinitions.Add(new ColumnDefinition { Width = new GridLength(1, GridUnitType.Star) });
        var prompt = new TextBlock { Text = ">", Foreground = new System.Windows.Media.SolidColorBrush(System.Windows.Media.Color.FromRgb(0xCC, 0xCC, 0xCC)), FontFamily = new System.Windows.Media.FontFamily("Cascadia Code, Consolas"), FontSize = 13, Margin = new Thickness(8, 6, 0, 6) };
        Grid.SetColumn(prompt, 0);
        var inputBox = new TextBox
        {
            Background = System.Windows.Media.Brushes.Transparent,
            Foreground = System.Windows.Media.Brushes.White,
            FontFamily = new System.Windows.Media.FontFamily("Cascadia Code, Consolas"),
            FontSize = 13,
            BorderThickness = new Thickness(0),
            Margin = new Thickness(0, 6, 8, 6)
        };
        inputBox.KeyDown += async (s, args) =>
        {
            if (args.Key == System.Windows.Input.Key.Enter)
            {
                var cmd = inputBox.Text;
                inputBox.Text = string.Empty;
                await terminal.SendInputAsync(cmd);
            }
        };
        Grid.SetColumn(inputBox, 1);
        inputGrid.Children.Add(prompt);
        inputGrid.Children.Add(inputBox);
        inputPanel.Child = inputGrid;
        Grid.SetRow(inputPanel, 2);

        grid.Children.Add(header);
        grid.Children.Add(output);
        grid.Children.Add(inputPanel);
        border.Child = grid;

        return border;
    }
}
