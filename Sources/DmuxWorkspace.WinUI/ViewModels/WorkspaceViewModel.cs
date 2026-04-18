using CommunityToolkit.Mvvm.ComponentModel;
using CommunityToolkit.Mvvm.Input;
using Codux.WinUI.Services;
using System.Collections.ObjectModel;

namespace Codux.WinUI.ViewModels;

public partial class WorkspaceViewModel : ObservableObject
{
    private readonly IProjectService _projectService;
    private readonly ITerminalService _terminalService;

    [ObservableProperty]
    private ObservableCollection<ProjectInfo> _projects = new();

    [ObservableProperty]
    private ProjectInfo? _selectedProject;

    [ObservableProperty]
    private ObservableCollection<TerminalPaneViewModel> _panes = new();

    [ObservableProperty]
    private TerminalPaneViewModel? _activePane;

    public WorkspaceViewModel()
    {
        _projectService = App.Services.GetRequiredService<IProjectService>();
        _terminalService = App.Services.GetRequiredService<ITerminalService>();

        _ = LoadProjectsAsync();
    }

    private async Task LoadProjectsAsync()
    {
        var projects = await _projectService.GetProjectsAsync();
        Projects = new ObservableCollection<ProjectInfo>(projects);

        if (Projects.Count > 0)
        {
            SelectedProject = Projects.First();
        }
    }

    [RelayCommand]
    private async Task AddProject()
    {
        var dialog = new Microsoft.Win32.OpenFolderDialog
        {
            Title = "Select Project Folder"
        };

        if (dialog.ShowDialog() == true)
        {
            var project = await _projectService.AddProjectAsync(dialog.FolderName);
            Projects.Add(project);
            SelectedProject = project;
        }
    }

    [RelayCommand]
    private async Task RemoveProject(ProjectInfo? project)
    {
        if (project == null) return;

        var panesToClose = Panes.Where(p => p.WorkingDirectory.StartsWith(project.Path)).ToList();
        foreach (var pane in panesToClose)
        {
            await ClosePaneInternalAsync(pane);
        }

        await _projectService.RemoveProjectAsync(project.Id);
        Projects.Remove(project);

        if (SelectedProject == project)
        {
            SelectedProject = Projects.FirstOrDefault();
        }
    }

    partial void OnSelectedProjectChanged(ProjectInfo? value)
    {
        if (value != null && Panes.Count == 0)
        {
            _ = NewTerminalPaneInternalAsync();
        }
    }

    [RelayCommand]
    private async Task NewTerminalPane()
    {
        await NewTerminalPaneInternalAsync();
    }

    private async Task NewTerminalPaneInternalAsync()
    {
        if (SelectedProject == null) return;

        var pane = new TerminalPaneViewModel
        {
            Id = Guid.NewGuid(),
            WorkingDirectory = SelectedProject.Path
        };

        await pane.CreateSessionAsync();
        Panes.Add(pane);
        ActivePane = pane;
    }

    [RelayCommand]
    private async Task ClosePane(TerminalPaneViewModel? pane)
    {
        if (pane == null) return;
        await ClosePaneInternalAsync(pane);
    }

    private async Task ClosePaneInternalAsync(TerminalPaneViewModel pane)
    {
        await pane.CloseSessionAsync();
        Panes.Remove(pane);

        if (ActivePane == pane)
        {
            ActivePane = Panes.FirstOrDefault();
        }
    }
}

public partial class TerminalPaneViewModel : ObservableObject
{
    private readonly ITerminalService _terminalService;

    public Guid Id { get; set; }

    [ObservableProperty]
    private string _workingDirectory = string.Empty;

    [ObservableProperty]
    private string _terminalOutput = string.Empty;

    private Guid _sessionId;
    private bool _sessionCreated;

    public TerminalPaneViewModel()
    {
        _terminalService = App.Services.GetRequiredService<ITerminalService>();
        _terminalService.OutputReceived += OnOutputReceived;
    }

    public async Task CreateSessionAsync()
    {
        if (_sessionCreated) return;

        _sessionId = await _terminalService.CreateSessionAsync(WorkingDirectory);
        _sessionCreated = true;

        TerminalOutput = $"PowerShell 7+ - {WorkingDirectory}\n";
        TerminalOutput += "Type commands below and press Enter to execute.\n";
        TerminalOutput += "--------------------------------------------\n\n";
    }

    public Task CloseSessionAsync()
    {
        if (!_sessionCreated) return Task.CompletedTask;

        _terminalService.OutputReceived -= OnOutputReceived;
        return _terminalService.CloseSessionAsync(_sessionId);
    }

    public async Task SendInputAsync(string input)
    {
        if (!_sessionCreated) return;

        TerminalOutput += $"PS> {input}\n";
        await _terminalService.WriteAsync(_sessionId, input);
    }

    private void OnOutputReceived(object? sender, TerminalOutputEventArgs e)
    {
        if (e.SessionId != _sessionId) return;

        Microsoft.UI.Xaml.Application.Current?.DispatcherQueue.TryEnqueue(() =>
        {
            TerminalOutput += e.Output;
        });
    }
}
