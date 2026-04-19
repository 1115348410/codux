using Serilog;
using System.IO;
using Newtonsoft.Json;

namespace Codux.WinUI.Services;

public class ProjectService : IProjectService
{
    private readonly string _dataPath;
    private readonly Dictionary<Guid, ProjectInfo> _projects = new();

    public ProjectService()
    {
        _dataPath = Path.Combine(
            Environment.GetFolderPath(Environment.SpecialFolder.LocalApplicationData),
            "Codux", "projects.json");
        LoadProjects();
    }

    private void LoadProjects()
    {
        try
        {
            if (File.Exists(_dataPath))
            {
                var json = File.ReadAllText(_dataPath);
                var projects = JsonConvert.DeserializeObject<List<ProjectInfo>>(json);
                if (projects != null)
                    foreach (var p in projects)
                        _projects[p.Id] = p;
            }
        }
        catch (Exception ex)
        {
            Log.Warning(ex, "Failed to load projects");
        }
    }

    private void SaveProjects()
    {
        try
        {
            var dir = Path.GetDirectoryName(_dataPath);
            if (dir != null) Directory.CreateDirectory(dir);
            var json = JsonConvert.SerializeObject(_projects.Values.ToList(), Formatting.Indented);
            File.WriteAllText(_dataPath, json);
        }
        catch (Exception ex)
        {
            Log.Warning(ex, "Failed to save projects");
        }
    }

    public Task<IEnumerable<ProjectInfo>> GetProjectsAsync(CancellationToken ct = default)
        => Task.FromResult<IEnumerable<ProjectInfo>>(_projects.Values.OrderByDescending(p => p.LastAccessed));

    public Task<ProjectInfo> AddProjectAsync(string path, CancellationToken ct = default)
    {
        var project = new ProjectInfo(Guid.NewGuid(), Path.GetFileName(path), path, DateTime.Now);
        _projects[project.Id] = project;
        SaveProjects();
        return Task.FromResult(project);
    }

    public Task RemoveProjectAsync(Guid id, CancellationToken ct = default)
    {
        if (_projects.Remove(id)) SaveProjects();
        return Task.CompletedTask;
    }
}
