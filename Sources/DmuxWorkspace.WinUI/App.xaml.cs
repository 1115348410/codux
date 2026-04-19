using System.Windows;
using Serilog;
using System.IO;
using Codux.WinUI.Services;
using Microsoft.Extensions.DependencyInjection;

namespace Codux.WinUI;

public partial class App : Application
{
    public static IServiceProvider Services { get; private set; } = null!;

    protected override void OnStartup(StartupEventArgs e)
    {
        InitializeLogging();
        Log.Information("Codux starting...");

        this.DispatcherUnhandledException += (s, args) =>
        {
            Log.Fatal(args.Exception, "Unhandled exception");
            Log.CloseAndFlush();
        };

        ConfigureServices();
        base.OnStartup(e);
    }

    private static void InitializeLogging()
    {
        var logPath = Path.Combine(
            Environment.GetFolderPath(Environment.SpecialFolder.LocalApplicationData),
            "Codux", "Logs", "codux-.log");

        Directory.CreateDirectory(Path.GetDirectoryName(logPath)!);

        Log.Logger = new LoggerConfiguration()
            .MinimumLevel.Debug()
            .WriteTo.File(logPath,
                rollingInterval: RollingInterval.Day,
                retainedFileCountLimit: 7,
                outputTemplate: "{Timestamp:yyyy-MM-dd HH:mm:ss.fff} [{Level:u3}] {Message:lj}{NewLine}{Exception}")
            .CreateLogger();
    }

    private static void ConfigureServices()
    {
        var services = new ServiceCollection();
        services.AddSingleton<ITerminalService, WindowsTerminalService>();
        services.AddSingleton<IProjectService, ProjectService>();
        services.AddSingleton<IGitService, GitService>();
        services.AddSingleton<IUsageService, UsageService>();
        services.AddSingleton<ISettingsService, SettingsService>();
        Services = services.BuildServiceProvider();
    }

    protected override void OnExit(ExitEventArgs e)
    {
        Log.Information("Codux shutting down");
        Log.CloseAndFlush();
        base.OnExit(e);
    }
}
