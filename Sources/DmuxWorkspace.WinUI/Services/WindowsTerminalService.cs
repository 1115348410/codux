using Serilog;
using System.Collections.Concurrent;
using System.Diagnostics;
using System.Runtime.InteropServices;

namespace Codux.WinUI.Services;

public class WindowsTerminalService : ITerminalService
{
    private readonly ConcurrentDictionary<Guid, TerminalSession> _sessions = new();

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

        var process = new Process { StartInfo = startInfo };
        var session = new TerminalSession
        {
            Id = sessionId,
            WorkingDirectory = workingDirectory,
            Process = process
        };

        process.OutputDataReceived += (s, e) =>
        {
            if (e.Data != null)
            {
                OutputReceived?.Invoke(this, new TerminalOutputEventArgs
                {
                    SessionId = sessionId,
                    Output = e.Data + Environment.NewLine
                });
            }
        };

        process.ErrorDataReceived += (s, e) =>
        {
            if (e.Data != null)
            {
                OutputReceived?.Invoke(this, new TerminalOutputEventArgs
                {
                    SessionId = sessionId,
                    Output = "[ERROR] " + e.Data + Environment.NewLine
                });
            }
        };

        process.Start();
        process.BeginOutputReadLine();
        process.BeginErrorReadLine();

        _sessions[sessionId] = session;
        Log.Information("Terminal session created: {SessionId} in {WorkingDirectory}", sessionId, workingDirectory);

        return Task.FromResult(sessionId);
    }

    public async Task WriteAsync(Guid sessionId, string input, CancellationToken ct = default)
    {
        if (_sessions.TryGetValue(sessionId, out var session))
        {
            if (!session.Process.HasExited)
            {
                await session.Process.StandardInput.WriteLineAsync(input);
                await session.Process.StandardInput.FlushAsync(ct);
            }
        }
    }

    public Task CloseSessionAsync(Guid sessionId, CancellationToken ct = default)
    {
        if (_sessions.TryRemove(sessionId, out var session))
        {
            try
            {
                if (!session.Process.HasExited)
                {
                    session.Process.Kill(entireProcessTree: true);
                }
                session.Process.Dispose();
                Log.Information("Terminal session closed: {SessionId}", sessionId);
            }
            catch (Exception ex)
            {
                Log.Warning(ex, "Error closing terminal session: {SessionId}", sessionId);
            }
        }
        return Task.CompletedTask;
    }

    public Task ResizeAsync(Guid sessionId, int width, int height, CancellationToken ct = default)
    {
        return Task.CompletedTask;
    }

    private class TerminalSession
    {
        public Guid Id { get; set; }
        public string WorkingDirectory { get; set; } = string.Empty;
        public Process Process { get; set; } = null!;
    }
}
