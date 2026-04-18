using Microsoft.UI.Xaml;
using Microsoft.UI.Xaml.Controls;
using Microsoft.UI.Xaml.Navigation;
using Codux.WinUI.Views;
using Codux.WinUI.ViewModels;
using Microsoft.Extensions.DependencyInjection;

namespace Codux.WinUI;

public partial class MainWindow : Window
{
    public MainWindow()
    {
        this.InitializeComponent();
        this.ExtendsContentIntoTitleBar = true;

        NavView.SelectionChanged += OnNavigationSelectionChanged;

        // Set initial page
        if (NavView.MenuItems.Count > 0 && NavView.MenuItems[0] is NavigationViewItem firstItem)
        {
            NavView.SelectedItem = firstItem;
            ContentFrame.Navigate(typeof(WorkspaceView));
        }
    }

    private void OnNavigationSelectionChanged(NavigationView sender, NavigationViewSelectionChangedEventArgs args)
    {
        if (args.SelectedItem is NavigationViewItem item)
        {
            var tag = item.Tag?.ToString();
            Type pageType = tag switch
            {
                "workspace" => typeof(WorkspaceView),
                "git" => typeof(GitPanelView),
                "ai" => typeof(AIStatsView),
                "settings" => typeof(SettingsView),
                _ => typeof(WorkspaceView)
            };

            ContentFrame.Navigate(pageType);
        }
    }
}
