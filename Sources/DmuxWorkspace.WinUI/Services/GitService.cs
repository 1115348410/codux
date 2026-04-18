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
                var current = repo.Head.FriendlyName;
                
                foreach (var branch in repo.Branches)
                {
                    branches.Add(new GitBranch(
                        branch.FriendlyName,
                        branch.IsRemote,
                        branch.FriendlyName == current));
                }
            }
            catch (Exception ex)
            {
                Log.Warning(ex, "Failed to get branches for {RepoPath}", repoPath);
            }
            return (IEnumerable<GitBranch>)branches;
        }, ct);
    }

    public Task<GitBranch?> GetCurrentBranchAsync(string repoPath, CancellationToken ct = default)
    {
        return Task.Run(() =>
        {
            try
            {
                using var repo = new Repository(repoPath);
                return new GitBranch(repo.Head.FriendlyName, false, true);
            }
            catch (Exception ex)
            {
                Log.Warning(ex, "Failed to get current branch for {RepoPath}", repoPath);
                return null;
            }
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
                    
                    var changeStatus = entry.State switch
                    {
                        FileStatus.ModifiedInIndex or FileStatus.ModifiedInWorkdir => GitChangeStatus.Modified,
                        FileStatus.NewInIndex or FileStatus.NewInWorkdir => GitChangeStatus.Added,
                        FileStatus.DeletedInIndex or FileStatus.DeletedInWorkdir => GitChangeStatus.Deleted,
                        FileStatus.RenamedInIndex or FileStatus.RenamedInWorkdir => GitChangeStatus.Renamed,
                        _ => GitChangeStatus.Modified
                    };

                    changes.Add(new GitChange(entry.FilePath, changeStatus, 
                        entry.State.HasFlag(FileStatus.ModifiedInIndex)));
                }
            }
            catch (Exception ex)
            {
                Log.Warning(ex, "Failed to get changes for {RepoPath}", repoPath);
            }
            return (IEnumerable<GitChange>)changes;
        }, ct);
    }

    public Task<IEnumerable<GitCommit>> GetCommitsAsync(string repoPath, int count = 50, CancellationToken ct = default)
    {
        return Task.Run(() =>
        {
            var commits = new List<GitCommit>();
            try
            {
                using var repo = new Repository(repoPath);
                foreach (var commit in repo.Commits.Take(count))
                {
                    commits.Add(new GitCommit(
                        commit.Sha[..7],
                        commit.MessageShort,
                        commit.Author.Name,
                        commit.Author.When.DateTime));
                }
            }
            catch (Exception ex)
            {
                Log.Warning(ex, "Failed to get commits for {RepoPath}", repoPath);
            }
            return (IEnumerable<GitCommit>)commits;
        }, ct);
    }

    public Task StageAsync(string repoPath, string filePath, CancellationToken ct = default)
    {
        return Task.Run(() =>
        {
            try
            {
                using var repo = new Repository(repoPath);
                Commands.Stage(repo, filePath);
            }
            catch (Exception ex)
            {
                Log.Warning(ex, "Failed to stage {FilePath} in {RepoPath}", filePath, repoPath);
            }
        }, ct);
    }

    public Task UnstageAsync(string repoPath, string filePath, CancellationToken ct = default)
    {
        return Task.Run(() =>
        {
            try
            {
                using var repo = new Repository(repoPath);
                Commands.Unstage(repo, filePath);
            }
            catch (Exception ex)
            {
                Log.Warning(ex, "Failed to unstage {FilePath} in {RepoPath}", filePath, repoPath);
            }
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

    public Task PushAsync(string repoPath, CancellationToken ct = default)
    {
        return Task.Run(() =>
        {
            try
            {
                using var repo = new Repository(repoPath);
                var origin = repo.Network.Remotes["origin"];
                if (origin != null)
                {
                    repo.Network.Push(origin, repo.Head.TrackedBranch, 
                        new PushOptions());
                }
            }
            catch (Exception ex)
            {
                Log.Warning(ex, "Failed to push in {RepoPath}", repoPath);
            }
        }, ct);
    }

    public Task PullAsync(string repoPath, CancellationToken ct = default)
    {
        return Task.Run(() =>
        {
            try
            {
                using var repo = new Repository(repoPath);
                var signature = repo.Config.BuildSignature(DateTime.Now);
                Commands.Pull(repo, signature, new PullOptions());
            }
            catch (Exception ex)
            {
                Log.Warning(ex, "Failed to pull in {RepoPath}", repoPath);
            }
        }, ct);
    }

    public Task FetchAsync(string repoPath, CancellationToken ct = default)
    {
        return Task.Run(() =>
        {
            try
            {
                using var repo = new Repository(repoPath);
                foreach (var remote in repo.Network.Remotes)
                {
                    var refSpecs = remote.FetchRefSpecs.Select(x => x.Specification);
                    Commands.Fetch(repo, remote.Name, refSpecs, new FetchOptions(), "");
                }
            }
            catch (Exception ex)
            {
                Log.Warning(ex, "Failed to fetch in {RepoPath}", repoPath);
            }
        }, ct);
    }
}
