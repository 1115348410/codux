using System.ComponentModel;
using System.Runtime.CompilerServices;
using Codux.WinUI.Services;

namespace Codux.WinUI.ViewModels;

public sealed class TerminalPaneViewModel : INotifyPropertyChanged
{
    private readonly ITerminalService _terminalService;
    private Guid _sessionId;
    private string _workingDirectory = string.Empty;
    private string _terminalOutput = string.Empty;

    public TerminalPaneViewModel(ITerminalService terminalService)
    {
        _terminalService = terminalService;
        _terminalService.OutputReceived += HandleOutputReceived;
    }

    public event PropertyChangedEventHandler? PropertyChanged;

    public string WorkingDirectory
    {
        get => _workingDirectory;
        set
        {
            if (_workingDirectory == value)
            {
                return;
            }

            _workingDirectory = value;
            OnPropertyChanged();
        }
    }

    public string TerminalOutput
    {
        get => _terminalOutput;
        private set
        {
            if (_terminalOutput == value)
            {
                return;
            }

            _terminalOutput = value;
            OnPropertyChanged();
        }
    }

    public async Task CreateSessionAsync(CancellationToken ct = default)
    {
        _sessionId = await _terminalService.CreateSessionAsync(WorkingDirectory, ct);
    }

    public Task SendInputAsync(string input, CancellationToken ct = default)
    {
        return _terminalService.WriteAsync(_sessionId, input, ct);
    }

    public void AppendClientOutput(string message)
    {
        if (string.IsNullOrWhiteSpace(message))
        {
            return;
        }

        TerminalOutput += $"> {message}\n";
    }

    public async Task CloseAsync(CancellationToken ct = default)
    {
        if (_sessionId != Guid.Empty)
        {
            await _terminalService.CloseSessionAsync(_sessionId, ct);
            _sessionId = Guid.Empty;
        }

        _terminalService.OutputReceived -= HandleOutputReceived;
    }

    private void HandleOutputReceived(object? sender, TerminalOutputEventArgs e)
    {
        if (e.SessionId != _sessionId)
        {
            return;
        }

        TerminalOutput += e.Output;
    }

    private void OnPropertyChanged([CallerMemberName] string? propertyName = null)
    {
        PropertyChanged?.Invoke(this, new PropertyChangedEventArgs(propertyName));
    }
}
