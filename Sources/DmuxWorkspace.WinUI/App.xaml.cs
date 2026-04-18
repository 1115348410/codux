using Microsoft.UI.Xaml;
using Microsoft.UI.Xaml.Navigation;
using Microsoft.Windows.AppLifecycle;
using Serilog;
using System.IO;
using Codux.WinUI.Services;
using Codux.WinUI.ViewModels;
using Microsoft.Extensions.DependencyInjection;

namespace Codux.WinUI;

public partial class App : Application
{
    private static IServiceProvider? _serviceProvider;
    public static IServiceProvider Services => _serviceProvider!;

    public App()
    {
        InitializeLogging();
        Log.Information("Codux starting...");

        this.UnhandledException += OnUnhandledException;
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

    protected override void OnLaunched(LaunchActivatedEventArgs args)
    {
        ConfigureServices();

        var mainWindow = new MainWindow();
        mainWindow.Activate();

        Log.Information("Codux launched successfully");
    }

    private static void ConfigureServices()
    {
        var services = new ServiceCollection();

        // Core Services
        services.AddSingleton<ITerminalService, WindowsTerminalService>();
        services.AddSingleton<INotificationService, ToastNotificationService>();
        services.AddSingleton<IProjectService, ProjectService>();
        services.AddSingleton<IGitService, GitService>();
        services.AddSingleton<IUsageService, UsageService>();
        services.AddSingleton<ISettingsService, SettingsService>();

        // ViewModels
        services.AddTransient<MainViewModel>();
        services.AddTransient<WorkspaceViewModel>();
        services.AddTransient<SettingsViewModel>();
        services.AddTransient<GitPanelViewModel>();
        services.AddTransient<AIStatsViewModel>();

        _serviceProvider = services.BuildServiceProvider();
    }

    private void OnUnhandledException(object sender, Microsoft.UI.Xaml.UnhandledExceptionEventArgs e)
    {
        Log.Fatal(e.Exception, "Unhandled exception occurred");
        Log.CloseAndFlush();
    }

    protected override void OnExit(ExitEventArgs e)
    {
        Log.Information("Codux shutting down");
        Log.CloseAndFlush();
        base.OnExit(e);
    }
}
