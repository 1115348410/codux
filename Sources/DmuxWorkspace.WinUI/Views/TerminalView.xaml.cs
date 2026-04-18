using Microsoft.UI.Xaml;
using Microsoft.UI.Xaml.Controls;
using Microsoft.UI.Xaml.Input;
using Codux.WinUI.ViewModels;

namespace Codux.WinUI.Views;

public partial class TerminalView : Page
{
    public TerminalPaneViewModel ViewModel { get; }

    public event EventHandler? CloseRequested;

    public TerminalView(TerminalPaneViewModel viewModel)
    {
        this.InitializeComponent();
        ViewModel = viewModel;
        this.DataContext = ViewModel;

        ViewModel.PropertyChanged += OnViewModelPropertyChanged;
        
        Loaded += OnLoaded;
    }

    private void OnLoaded(object sender, RoutedEventArgs e)
    {
        InputBox.Focus(FocusState.Programmatic);
    }

    private void OnViewModelPropertyChanged(object? sender, System.ComponentModel.PropertyChangedEventArgs e)
    {
        if (e.PropertyName == nameof(ViewModel.TerminalOutput))
        {
            OutputScroller.ScrollToEnd();
        }
    }

    private async void OnInputKeyDown(object sender, KeyRoutedEventArgs e)
    {
        if (e.Key == Windows.System.VirtualKey.Enter)
        {
            var input = InputBox.Text;
            InputBox.Text = string.Empty;

            if (!string.IsNullOrWhiteSpace(input))
            {
                await ViewModel.SendInputAsync(input);
            }
        }
    }

    private void OnCloseClick(object sender, RoutedEventArgs e)
    {
        CloseRequested?.Invoke(this, EventArgs.Empty);
    }
}
