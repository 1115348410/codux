namespace Codux.WinUI.Services;

public interface IGitService
{
    Task<IEnumerable<GitBranch>> GetBranchesAsync(string repoPath, CancellationToken ct = default);
    Task<IEnumerable<GitChange>> GetChangesAsync(string repoPath, CancellationToken ct = default);
    Task CommitAsync(string repoPath, string message, CancellationToken ct = default);
}

public record GitBranch(string Name, bool IsRemote, bool IsCurrent);
public record GitChange(string FilePath, string Status, bool IsStaged);
