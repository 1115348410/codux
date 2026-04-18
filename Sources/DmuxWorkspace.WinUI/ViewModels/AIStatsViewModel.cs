using CommunityToolkit.Mvvm.ComponentModel;
using CommunityToolkit.Mvvm.Input;
using Codux.WinUI.Services;
using System.Collections.ObjectModel;

namespace Codux.WinUI.ViewModels;

public partial class AIStatsViewModel : ObservableObject
{
    private readonly IUsageService _usageService;

    [ObservableProperty]
    private ObservableCollection<UsageRecord> _recentUsage = new();

    [ObservableProperty]
    private UsageSummary? _summary;

    [ObservableProperty]
    private int _totalInputTokens;

    [ObservableProperty]
    private int _totalOutputTokens;

    [ObservableProperty]
    private string _currentLevel = "Bronze";

    [ObservableProperty]
    private double _levelProgress;

    [ObservableProperty]
    private bool _isLoading;

    public AIStatsViewModel()
    {
        _usageService = App.Services.GetRequiredService<IUsageService>();
    }

    public async Task LoadDataAsync()
    {
        IsLoading = true;

        try
        {
            var since = DateTime.Today.AddDays(-7);
            
            var records = await _usageService.GetUsageRecordsAsync(since);
            RecentUsage = new ObservableCollection<UsageRecord>(records);

            Summary = await _usageService.GetSummaryAsync(since);
            
            TotalInputTokens = Summary?.TotalInputTokens ?? 0;
            TotalOutputTokens = Summary?.TotalOutputTokens ?? 0;

            CalculateLevel();
        }
        finally
        {
            IsLoading = false;
        }
    }

    private void CalculateLevel()
    {
        var totalTokens = TotalInputTokens + TotalOutputTokens;
        
        (CurrentLevel, LevelProgress) = totalTokens switch
        {
            < 100_000 => ("Iron", totalTokens / 100_000.0),
            < 500_000 => ("Bronze", (totalTokens - 100_000) / 400_000.0),
            < 2_000_000 => ("Silver", (totalTokens - 500_000) / 1_500_000.0),
            < 5_000_000 => ("Gold", (totalTokens - 2_000_000) / 3_000_000.0),
            < 10_000_000 => ("Platinum", (totalTokens - 5_000_000) / 5_000_000.0),
            < 50_000_000 => ("Diamond", (totalTokens - 10_000_000) / 40_000_000.0),
            < 100_000_000 => ("Master", (totalTokens - 50_000_000) / 50_000_000.0),
            _ => ("Grandmaster", Math.Min((totalTokens - 100_000_000) / 100_000_000.0, 1.0))
        };
    }
}
