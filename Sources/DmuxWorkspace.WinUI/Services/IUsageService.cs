namespace Codux.WinUI.Services;

public interface IUsageService
{
    Task RecordUsageAsync(int inputTokens, int outputTokens, string model, CancellationToken ct = default);
    Task<int> GetTotalUsageAsync(CancellationToken ct = default);
}

public record UsageRecord(int InputTokens, int OutputTokens, string Model, DateTime Timestamp);
