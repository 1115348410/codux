namespace Codux.WinUI.Services;

public interface IGitService
{
    Task<IEnumerable<GitBranch>> GetBranchesAsync(string repoPath, CancellationToken ct = default);
    Task<GitBranch?> GetCurrentBranchAsync(string repoPath, CancellationToken ct = default);
    Task<IEnumerable<GitChange>> GetChangesAsync(string repoPath, CancellationToken ct = default);
    Task<IEnumerable<GitCommit>> GetCommitsAsync(string repoPath, int count = 50, CancellationToken ct = default);
    Task StageAsync(string repoPath, string filePath, CancellationToken ct = default);
    Task UnstageAsync(string repoPath, string filePath, CancellationToken ct = default);
    Task CommitAsync(string repoPath, string message, CancellationToken ct = default);
    Task PushAsync(string repoPath, CancellationToken ct = default);
    Task PullAsync(string repoPath, CancellationToken ct = default);
    Task FetchAsync(string repoPath, CancellationToken ct = default);
}

public record GitBranch(string Name, bool IsRemote, bool IsCurrent);
public record GitChange(string FilePath, GitChangeStatus Status, bool IsStaged);
public record GitCommit(string Sha, string Message, string Author, DateTime Date);

public enum GitChangeStatus { Modified, Added, Deleted, Renamed, Untracked }
