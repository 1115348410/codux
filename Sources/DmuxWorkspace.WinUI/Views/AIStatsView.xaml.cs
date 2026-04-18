using Microsoft.UI.Xaml.Controls;
using Codux.WinUI.ViewModels;

namespace Codux.WinUI.Views;

public partial class AIStatsView : Page
{
    public AIStatsViewModel ViewModel { get; }

    public AIStatsView()
    {
        this.InitializeComponent();
        ViewModel = App.Services.GetRequiredService<AIStatsViewModel>();
        this.DataContext = ViewModel;
    }

    public async Task LoadDataAsync()
    {
        await ViewModel.LoadDataAsync();
    }
}
