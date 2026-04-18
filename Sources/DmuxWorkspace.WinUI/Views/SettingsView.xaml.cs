using Microsoft.UI.Xaml;
using Microsoft.UI.Xaml.Controls;
using Codux.WinUI.ViewModels;

namespace Codux.WinUI.Views;

public partial class SettingsView : Page
{
    public SettingsViewModel ViewModel { get; }

    public SettingsView()
    {
        this.InitializeComponent();
        ViewModel = App.Services.GetRequiredService<SettingsViewModel>();
        this.DataContext = ViewModel;
    }

    private void OnSaveClick(object sender, RoutedEventArgs e)
    {
        ViewModel.SaveSettings();
    }
}
