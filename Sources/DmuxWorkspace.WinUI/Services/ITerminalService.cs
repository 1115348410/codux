namespace Codux.WinUI.Services;

public interface ITerminalService
{
    Task<Guid> CreateSessionAsync(string workingDirectory, CancellationToken ct = default);
    Task WriteAsync(Guid sessionId, string input, CancellationToken ct = default);
    Task CloseSessionAsync(Guid sessionId, CancellationToken ct = default);
    Task ResizeAsync(Guid sessionId, int width, int height, CancellationToken ct = default);
    event EventHandler<TerminalOutputEventArgs>? OutputReceived;
}

public class TerminalOutputEventArgs : EventArgs
{
    public Guid SessionId { get; set; }
    public string Output { get; set; } = string.Empty;
}
