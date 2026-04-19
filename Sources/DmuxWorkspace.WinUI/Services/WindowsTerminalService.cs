using Serilog;
using System.Diagnostics;

namespace Codux.WinUI.Services;

public class WindowsTerminalService : ITerminalService
{
    private readonly Dictionary<Guid, Process> _processes = new();

    public event EventHandler<TerminalOutputEventArgs>? OutputReceived;

    public Task<Guid> CreateSessionAsync(string workingDirectory, CancellationToken ct = default)
    {
        var sessionId = Guid.NewGuid();

        var startInfo = new ProcessStartInfo
        {
            FileName = "powershell.exe",
            Arguments = "-NoLogo -NoExit -Command -",
            WorkingDirectory = workingDirectory,
            UseShellExecute = false,
            RedirectStandardInput = true,
            RedirectStandardOutput = true,
            RedirectStandardError = true,
            CreateNoWindow = true,
            StandardOutputEncoding = System.Text.Encoding.UTF8,
            StandardErrorEncoding = System.Text.Encoding.UTF8
        };

        var process = Process.Start(startInfo);
        if (process == null) throw new Exception("Failed to start process");

        _processes[sessionId] = process;

        process.OutputDataReceived += (s, e) =>
        {
            if (e.Data != null)
                OutputReceived?.Invoke(this, new TerminalOutputEventArgs { SessionId = sessionId, Output = e.Data + "\n" });
        };

        process.ErrorDataReceived += (s, e) =>
        {
            if (e.Data != null)
                OutputReceived?.Invoke(this, new TerminalOutputEventArgs { SessionId = sessionId, Output = "[ERROR] " + e.Data + "\n" });
        };

        process.BeginOutputReadLine();
        process.BeginErrorReadLine();

        Log.Information("Terminal session created: {SessionId}", sessionId);
        return Task.FromResult(sessionId);
    }

    public async Task WriteAsync(Guid sessionId, string input, CancellationToken ct = default)
    {
        if (_processes.TryGetValue(sessionId, out var process) && !process.HasExited)
        {
            await process.StandardInput.WriteLineAsync(input);
            await process.StandardInput.FlushAsync(ct);
        }
    }

    public Task CloseSessionAsync(Guid sessionId, CancellationToken ct = default)
    {
        if (_processes.TryGetValue(sessionId, out var process))
        {
            try
            {
                if (!process.HasExited)
                    process.Kill(entireProcessTree: true);
                process.Dispose();
            }
            catch (Exception ex)
            {
                Log.Warning(ex, "Error closing session {SessionId}", sessionId);
            }
            _processes.Remove(sessionId);
        }
        return Task.CompletedTask;
    }
}
