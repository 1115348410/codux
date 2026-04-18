namespace Codux.WinUI.Services;

public interface IUsageService
{
    Task<IEnumerable<UsageRecord>> GetUsageRecordsAsync(DateTime since, CancellationToken ct = default);
    Task RecordUsageAsync(UsageRecord record, CancellationToken ct = default);
    Task<UsageSummary> GetSummaryAsync(DateTime since, CancellationToken ct = default);
}

public record UsageRecord(
    string Model,
    int InputTokens,
    int OutputTokens,
    DateTime Timestamp,
    string SessionId);

public record UsageSummary(
    int TotalInputTokens,
    int TotalOutputTokens,
    Dictionary<string, int> TokensByModel,
    int TotalSessions);
