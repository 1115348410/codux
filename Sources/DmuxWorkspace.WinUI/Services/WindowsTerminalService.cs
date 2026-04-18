using Serilog;
using System.Collections.Concurrent;

namespace Codux.WinUI.Services;

public class WindowsTerminalService : ITerminalService
{
    private readonly ConcurrentDictionary<Guid, TerminalSession> _sessions = new();

    public event EventHandler<TerminalOutputEventArgs>? OutputReceived;

    public Task<Guid> CreateSessionAsync(string workingDirectory, CancellationToken ct = default)
    {
        var sessionId = Guid.NewGuid();
        
        var session = new TerminalSession
        {
            Id = sessionId,
            WorkingDirectory = workingDirectory,
            Process = new System.Diagnostics.Process()
        };

        session.Process.StartInfo = new System.Diagnostics.ProcessStartInfo
        {
            FileName = "cmd.exe",
            WorkingDirectory = workingDirectory,
            UseShellExecute = false,
            RedirectStandardInput = true,
            RedirectStandardOutput = true,
            RedirectStandardError = true,
            CreateNoWindow = true
        };

        session.Process.OutputDataReceived += (s, e) =>
        {
            if (!string.IsNullOrEmpty(e.Data))
            {
                OutputReceived?.Invoke(this, new TerminalOutputEventArgs
                {
                    SessionId = sessionId,
                    Output = e.Data + Environment.NewLine
                });
            }
        };

        session.Process.ErrorDataReceived += (s, e) =>
        {
            if (!string.IsNullOrEmpty(e.Data))
            {
                OutputReceived?.Invoke(this, new TerminalOutputEventArgs
                {
                    SessionId = sessionId,
                    Output = e.Data + Environment.NewLine
                });
            }
        };

        session.Process.Start();
        session.Process.BeginOutputReadLine();
        session.Process.BeginErrorReadLine();

        _sessions[sessionId] = session;
        Log.Information("Terminal session created: {SessionId}", sessionId);

        return Task.FromResult(sessionId);
    }

    public Task WriteAsync(Guid sessionId, string input, CancellationToken ct = default)
    {
        if (_sessions.TryGetValue(sessionId, out var session))
        {
            session.Process.StandardInput.WriteLine(input);
            session.Process.StandardInput.Flush();
        }
        return Task.CompletedTask;
    }

    public Task CloseSessionAsync(Guid sessionId, CancellationToken ct = default)
    {
        if (_sessions.TryRemove(sessionId, out var session))
        {
            if (!session.Process.HasExited)
            {
                session.Process.Kill();
            }
            session.Process.Dispose();
            Log.Information("Terminal session closed: {SessionId}", sessionId);
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
        public System.Diagnostics.Process Process { get; set; } = null!;
    }
}
