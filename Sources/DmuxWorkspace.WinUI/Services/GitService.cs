using LibGit2Sharp;
using Serilog;

namespace Codux.WinUI.Services;

public class GitService : IGitService
{
    public Task<IEnumerable<GitBranch>> GetBranchesAsync(string repoPath, CancellationToken ct = default)
    {
        return Task.Run(() =>
        {
            var branches = new List<GitBranch>();
            try
            {
                using var repo = new Repository(repoPath);
                foreach (var branch in repo.Branches)
                    branches.Add(new GitBranch(branch.FriendlyName, branch.IsRemote, branch.FriendlyName == repo.Head.FriendlyName));
            }
            catch (Exception ex)
            {
                Log.Warning(ex, "Failed to get branches for {RepoPath}", repoPath);
            }
            return (IEnumerable<GitBranch>)branches;
        }, ct);
    }

    public Task<IEnumerable<GitChange>> GetChangesAsync(string repoPath, CancellationToken ct = default)
    {
        return Task.Run(() =>
        {
            var changes = new List<GitChange>();
            try
            {
                using var repo = new Repository(repoPath);
                var status = repo.RetrieveStatus();
                foreach (var entry in status)
                {
                    if (entry.State == FileStatus.Ignored) continue;
                    changes.Add(new GitChange(entry.FilePath, entry.State.ToString(), entry.State.HasFlag(FileStatus.ModifiedInIndex)));
                }
            }
            catch (Exception ex)
            {
                Log.Warning(ex, "Failed to get changes for {RepoPath}", repoPath);
            }
            return (IEnumerable<GitChange>)changes;
        }, ct);
    }

    public Task CommitAsync(string repoPath, string message, CancellationToken ct = default)
    {
        return Task.Run(() =>
        {
            try
            {
                using var repo = new Repository(repoPath);
                var signature = repo.Config.BuildSignature(DateTime.Now);
                repo.Commit(message, signature, signature);
            }
            catch (Exception ex)
            {
                Log.Warning(ex, "Failed to commit in {RepoPath}", repoPath);
            }
        }, ct);
    }
}
