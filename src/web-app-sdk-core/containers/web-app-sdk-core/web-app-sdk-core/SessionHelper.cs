using VideoOS.Platform.SDK.Core;

namespace TestWebApp;

public class SessionHelper
{
    public static Session CreateSessionWithServerConfigurationProvided(IServiceProvider serviceProvider, string serverUrl, string userType, string? username = null, string? password = null)
    {
        var serverUri = new Uri(serverUrl);
        var serverConfiguration = new ServerConfiguration(serverUri, new Uri(serverUri + "idp"));

        return userType switch
        {
            "DefaultWindows" => new Session(serverConfiguration, serviceProvider, new DefaultWindowsUser()),
            "Windows" => new Session(serverConfiguration, serviceProvider, new WindowsUser(username ?? "", password ?? "")),
            "Basic" => new Session(serverConfiguration, serviceProvider, new BasicUser(username ?? "", password ?? "")),
            "External" => new Session(serverConfiguration, serviceProvider, password ?? ""),
            _ => throw new ArgumentException($"Invalid user type: {userType}. Valid types: DefaultWindows, Windows, Basic, External")
        };
    }

    public static async Task<string> WaitForTokenAsync(Session session, int timeoutSeconds = 30)
    {
        var existing = session.MipTokenCache.Token;
        if (!string.IsNullOrEmpty(existing))
            return existing;

        var tcs = new TaskCompletionSource<string>(TaskCreationOptions.RunContinuationsAsynchronously);

        void OnToken(string token) { session.MipTokenCache.OnNewTokenAvailable -= OnToken; tcs.TrySetResult(token); }
        void OnErr(Exception ex) { session.MipTokenCache.OnError -= OnErr; tcs.TrySetException(ex); }

        session.MipTokenCache.OnNewTokenAvailable += OnToken;
        session.MipTokenCache.OnError += OnErr;

        using var cts = new CancellationTokenSource(TimeSpan.FromSeconds(timeoutSeconds));
        cts.Token.Register(() => tcs.TrySetException(new TimeoutException($"Token acquisition timed out after {timeoutSeconds}s.")));

        return await tcs.Task;
    }
}
