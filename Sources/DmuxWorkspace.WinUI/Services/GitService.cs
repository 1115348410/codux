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
                if (repo.RetrieveStatus().IsDirty)
                {
                    Commands.Stage(repo, "*");
                }

                if (!repo.RetrieveStatus().Any(entry =>
                        entry.State.HasFlag(FileStatus.NewInIndex) ||
                        entry.State.HasFlag(FileStatus.ModifiedInIndex) ||
                        entry.State.HasFlag(FileStatus.DeletedFromIndex) ||
                        entry.State.HasFlag(FileStatus.RenamedInIndex) ||
                        entry.State.HasFlag(FileStatus.TypeChangeInIndex)))
                {
                    return;
                }

                var signature = repo.Config.BuildSignature(DateTimeOffset.Now)
                    ?? new Signature("Codux", "codux@local", DateTimeOffset.Now);
                repo.Commit(message, signature, signature);
            }
            catch (Exception ex)
            {
                Log.Warning(ex, "Failed to commit in {RepoPath}", repoPath);
                throw;
            }
        }, ct);
    }

    public Task StageAllAsync(string repoPath, CancellationToken ct = default)
    {
        return Task.Run(() =>
        {
            try
            {
                using var repo = new Repository(repoPath);
                Commands.Stage(repo, "*");
            }
            catch (Exception ex)
            {
                Log.Warning(ex, "Failed to stage changes in {RepoPath}", repoPath);
                throw;
            }
        }, ct);
    }
}
