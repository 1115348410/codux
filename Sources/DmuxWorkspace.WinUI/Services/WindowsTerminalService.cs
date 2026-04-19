using Serilog;
using System.Diagnostics;
using System.IO;
using System.Text;

namespace Codux.WinUI.Services;

public class WindowsTerminalService : ITerminalService
{
    private readonly Dictionary<Guid, Process> _processes = new();

    public event EventHandler<TerminalOutputEventArgs>? OutputReceived;

    public Task<Guid> CreateSessionAsync(string workingDirectory, CancellationToken ct = default)
    {
        var sessionId = Guid.NewGuid();

        workingDirectory = workingDirectory.Trim();
        if (!Directory.Exists(workingDirectory))
        {
            Log.Warning("Directory does not exist: {WorkingDirectory}, using user profile", workingDirectory);
            workingDirectory = Environment.GetFolderPath(Environment.SpecialFolder.UserProfile);
        }

        var startInfo = new ProcessStartInfo
        {
            FileName = "cmd.exe",
            Arguments = "/Q",
            WorkingDirectory = workingDirectory,
            UseShellExecute = false,
            RedirectStandardInput = true,
            RedirectStandardOutput = true,
            RedirectStandardError = true,
            CreateNoWindow = true,
            StandardOutputEncoding = Encoding.UTF8,
            StandardErrorEncoding = Encoding.UTF8
        };

        var process = new Process { StartInfo = startInfo };
        _processes[sessionId] = process;

        process.OutputDataReceived += (s, e) =>
        {
            if (e.Data != null)
            {
                OutputReceived?.Invoke(this, new TerminalOutputEventArgs { SessionId = sessionId, Output = e.Data + "\n" });
            }
        };

        process.ErrorDataReceived += (s, e) =>
        {
            if (e.Data != null)
            {
                OutputReceived?.Invoke(this, new TerminalOutputEventArgs { SessionId = sessionId, Output = "[ERROR] " + e.Data + "\n" });
            }
        };

        process.Start();
        process.BeginOutputReadLine();
        process.BeginErrorReadLine();

        // Send initial prompt command to show working directory
        var initCmd = $"@echo off & cd /d \"{workingDirectory}\" & prompt $G$S & echo.\r\n";
        process.StandardInput.Write(initCmd);
        process.StandardInput.Flush();

        Log.Information("Terminal session created: {SessionId} in {WorkingDirectory}", sessionId, workingDirectory);
        return Task.FromResult(sessionId);
    }

    public async Task WriteAsync(Guid sessionId, string input, CancellationToken ct = default)
    {
        if (_processes.TryGetValue(sessionId, out var process) && !process.HasExited)
        {
            // Wrap command in cmd.exe format
            var cmd = input + "\r\n";
            await process.StandardInput.WriteAsync(cmd);
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
