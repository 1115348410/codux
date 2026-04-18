using CommunityToolkit.Mvvm.ComponentModel;
using CommunityToolkit.Mvvm.Input;
using Codux.WinUI.Services;
using System.Collections.ObjectModel;

namespace Codux.WinUI.ViewModels;

public partial class GitPanelViewModel : ObservableObject
{
    private readonly IGitService _gitService;
    private string _currentRepoPath = string.Empty;

    [ObservableProperty]
    private ObservableCollection<GitBranch> _branches = new();

    [ObservableProperty]
    private GitBranch? _currentBranch;

    [ObservableProperty]
    private ObservableCollection<GitChange> _changes = new();

    [ObservableProperty]
    private ObservableCollection<GitCommit> _commits = new();

    [ObservableProperty]
    private string _commitMessage = string.Empty;

    [ObservableProperty]
    private bool _isLoading;

    [ObservableProperty]
    private string _statusMessage = string.Empty;

    public GitPanelViewModel()
    {
        _gitService = App.Services.GetRequiredService<IGitService>();
    }

    public async Task LoadRepositoryAsync(string repoPath)
    {
        if (string.IsNullOrEmpty(repoPath)) return;
        
        _currentRepoPath = repoPath;
        IsLoading = true;
        StatusMessage = "Loading...";

        try
        {
            var branches = await _gitService.GetBranchesAsync(repoPath);
            Branches = new ObservableCollection<GitBranch>(branches);
            CurrentBranch = Branches.FirstOrDefault(b => b.IsCurrent);

            var changes = await _gitService.GetChangesAsync(repoPath);
            Changes = new ObservableCollection<GitChange>(changes);

            var commits = await _gitService.GetCommitsAsync(repoPath);
            Commits = new ObservableCollection<GitCommit>(commits);

            StatusMessage = $"{Changes.Count} changes";
        }
        finally
        {
            IsLoading = false;
        }
    }

    [RelayCommand]
    private async Task Stage(GitChange? change)
    {
        if (change == null || string.IsNullOrEmpty(_currentRepoPath)) return;
        
        await _gitService.StageAsync(_currentRepoPath, change.FilePath);
        await RefreshChangesAsync();
    }

    [RelayCommand]
    private async Task Unstage(GitChange? change)
    {
        if (change == null || string.IsNullOrEmpty(_currentRepoPath)) return;
        
        await _gitService.UnstageAsync(_currentRepoPath, change.FilePath);
        await RefreshChangesAsync();
    }

    [RelayCommand]
    private async Task Commit()
    {
        if (string.IsNullOrEmpty(_currentRepoPath) || string.IsNullOrWhiteSpace(CommitMessage)) return;
        
        IsLoading = true;
        StatusMessage = "Committing...";

        try
        {
            await _gitService.CommitAsync(_currentRepoPath, CommitMessage);
            CommitMessage = string.Empty;
            await RefreshCommitsAsync();
            await RefreshChangesAsync();
            StatusMessage = "Committed successfully";
        }
        finally
        {
            IsLoading = false;
        }
    }

    [RelayCommand]
    private async Task Push()
    {
        if (string.IsNullOrEmpty(_currentRepoPath)) return;
        
        IsLoading = true;
        StatusMessage = "Pushing...";

        try
        {
            await _gitService.PushAsync(_currentRepoPath);
            StatusMessage = "Pushed successfully";
        }
        finally
        {
            IsLoading = false;
        }
    }

    [RelayCommand]
    private async Task Pull()
    {
        if (string.IsNullOrEmpty(_currentRepoPath)) return;
        
        IsLoading = true;
        StatusMessage = "Pulling...";

        try
        {
            await _gitService.PullAsync(_currentRepoPath);
            await RefreshCommitsAsync();
            await RefreshChangesAsync();
            StatusMessage = "Pulled successfully";
        }
        finally
        {
            IsLoading = false;
        }
    }

    [RelayCommand]
    private async Task Fetch()
    {
        if (string.IsNullOrEmpty(_currentRepoPath)) return;
        
        IsLoading = true;
        StatusMessage = "Fetching...";

        try
        {
            await _gitService.FetchAsync(_currentRepoPath);
            await LoadRepositoryAsync(_currentRepoPath);
            StatusMessage = "Fetched successfully";
        }
        finally
        {
            IsLoading = false;
        }
    }

    private async Task RefreshChangesAsync()
    {
        if (string.IsNullOrEmpty(_currentRepoPath)) return;
        
        var changes = await _gitService.GetChangesAsync(_currentRepoPath);
        Changes = new ObservableCollection<GitChange>(changes);
    }

    private async Task RefreshCommitsAsync()
    {
        if (string.IsNullOrEmpty(_currentRepoPath)) return;
        
        var commits = await _gitService.GetCommitsAsync(_currentRepoPath);
        Commits = new ObservableCollection<GitCommit>(commits);
    }
}
