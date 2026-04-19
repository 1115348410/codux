using Serilog;
using System.IO;
using Newtonsoft.Json;

namespace Codux.WinUI.Services;

public class UsageService : IUsageService
{
    private readonly string _dataPath;
    private readonly List<UsageRecord> _records = new();

    public UsageService()
    {
        _dataPath = Path.Combine(
            Environment.GetFolderPath(Environment.SpecialFolder.LocalApplicationData),
            "Codux", "usage.json");
        LoadRecords();
    }

    private void LoadRecords()
    {
        try
        {
            if (File.Exists(_dataPath))
            {
                var json = File.ReadAllText(_dataPath);
                var records = JsonConvert.DeserializeObject<List<UsageRecord>>(json);
                if (records != null) _records.AddRange(records);
            }
        }
        catch (Exception ex)
        {
            Log.Warning(ex, "Failed to load usage records");
        }
    }

    private void SaveRecords()
    {
        try
        {
            var dir = Path.GetDirectoryName(_dataPath);
            if (dir != null) Directory.CreateDirectory(dir);
            var json = JsonConvert.SerializeObject(_records, Formatting.Indented);
            File.WriteAllText(_dataPath, json);
        }
        catch (Exception ex)
        {
            Log.Warning(ex, "Failed to save usage records");
        }
    }

    public Task RecordUsageAsync(int inputTokens, int outputTokens, string model, CancellationToken ct = default)
    {
        _records.Add(new UsageRecord(inputTokens, outputTokens, model, DateTime.Now));
        SaveRecords();
        return Task.CompletedTask;
    }

    public Task<int> GetTotalUsageAsync(CancellationToken ct = default)
        => Task.FromResult(_records.Sum(r => r.InputTokens + r.OutputTokens));
}
