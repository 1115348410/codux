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
                if (records != null)
                {
                    _records.AddRange(records);
                }
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

    public Task<IEnumerable<UsageRecord>> GetUsageRecordsAsync(DateTime since, CancellationToken ct = default)
    {
        var records = _records.Where(r => r.Timestamp >= since).OrderByDescending(r => r.Timestamp);
        return Task.FromResult<IEnumerable<UsageRecord>>(records);
    }

    public Task RecordUsageAsync(UsageRecord record, CancellationToken ct = default)
    {
        _records.Add(record);
        SaveRecords();
        return Task.CompletedTask;
    }

    public Task<UsageSummary> GetSummaryAsync(DateTime since, CancellationToken ct = default)
    {
        var filtered = _records.Where(r => r.Timestamp >= since).ToList();
        
        var summary = new UsageSummary(
            filtered.Sum(r => r.InputTokens),
            filtered.Sum(r => r.OutputTokens),
            filtered.GroupBy(r => r.Model)
                    .ToDictionary(g => g.Key, g => g.Sum(r => r.InputTokens + r.OutputTokens)),
            filtered.Select(r => r.SessionId).Distinct().Count());

        return Task.FromResult(summary);
    }
}
