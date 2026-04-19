namespace Codux.WinUI.Services;

public interface IProjectService
{
    Task<IEnumerable<ProjectInfo>> GetProjectsAsync(CancellationToken ct = default);
    Task<ProjectInfo> AddProjectAsync(string path, CancellationToken ct = default);
    Task RemoveProjectAsync(Guid id, CancellationToken ct = default);
}

public record ProjectInfo(Guid Id, string Name, string Path, DateTime LastAccessed);
