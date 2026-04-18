namespace Codux.WinUI.Services;

public interface IProjectService
{
    Task<IEnumerable<ProjectInfo>> GetProjectsAsync(CancellationToken ct = default);
    Task<ProjectInfo> AddProjectAsync(string path, CancellationToken ct = default);
    Task RemoveProjectAsync(Guid id, CancellationToken ct = default);
    event EventHandler<ProjectChangedEventArgs>? ProjectChanged;
}

public record ProjectInfo(Guid Id, string Name, string Path, DateTime LastAccessed);
public record ProjectChangedEventArgs(Guid ProjectId, ProjectChangeType ChangeType);

public enum ProjectChangeType { Added, Removed, Updated }
