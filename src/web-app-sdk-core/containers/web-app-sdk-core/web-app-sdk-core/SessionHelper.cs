using System;
using VideoOS.Platform.SDK.Core;

namespace TestWebApp;

public class SessionHelper
{
    public static Session CreateSessionWithServerConfigurationProvided(IServiceProvider serviceProvider, string serverUrl, string userType, string? username = null, string? password = null)
    {
        var serverUri = new Uri(serverUrl);
        var idpUri = new Uri(serverUri + "idp");
        var serverConfiguration = new ServerConfiguration(serverUri, idpUri);

        return userType switch
        {
            "DefaultWindows" => new Session(serverConfiguration, serviceProvider, new DefaultWindowsUser()),
            "Windows" => new Session(serverConfiguration, serviceProvider, new WindowsUser(username ?? "", password ?? "")),
            "Basic" => new Session(serverConfiguration, serviceProvider, new BasicUser(username ?? "", password ?? "")),
            "External" => new Session(serverConfiguration, serviceProvider, password ?? ""), // password field used for access token
            _ => throw new ArgumentException($"Invalid user type: {userType}. Valid types: DefaultWindows, Windows, Basic, External")
        };
    }

    public static Session CreateSessionWithRuntimeConfiguration(IServiceProvider serviceProvider, string userType, string? username = null, string? password = null)
    {
        var serverConfiguration = new RuntimeServerConfiguration();
        return userType switch
        {
            "DefaultWindows" => new Session(serverConfiguration, serviceProvider, new DefaultWindowsUser()),
            "Windows" => new Session(serverConfiguration, serviceProvider, new WindowsUser(username ?? "", password ?? "")),
            "Basic" => new Session(serverConfiguration, serviceProvider, new BasicUser(username ?? "", password ?? "")),
            "External" => new Session(serverConfiguration, serviceProvider, password ?? ""), // password field used for access token
            _ => throw new ArgumentException($"Invalid user type: {userType}. Valid types: DefaultWindows, Windows, Basic, External")
        };
    }

    public static Session CreateSessionWithRuntimeConfigurationDefaultUser(IServiceProvider serviceProvider)
    {
        var serverConfiguration = new RuntimeServerConfiguration();
        return new Session(serverConfiguration, serviceProvider);
    }
}
