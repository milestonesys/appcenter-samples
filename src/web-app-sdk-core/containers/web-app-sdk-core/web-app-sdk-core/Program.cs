using System.Collections.Concurrent;
using TestWebApp;
using VideoOS.Platform.SDK.Core.Configuration;
using VideoOS.Platform.SDK.Core.Configuration.Filtering;
using VideoOS.Platform.SDK.Core.Configuration.Items;
using VideoOS.Platform.SDK.Core.Extensions;
using VideoOS.Platform.SDK.Core.Media;

var builder = WebApplication.CreateBuilder(args);
builder.Services.AddMipServices();
var app = builder.Build();

// Serve wwwroot/index.html as the default page for the browser UI.
app.UseDefaultFiles();
app.UseStaticFiles();

// Returns the VMS server URL assembled from runtime environment variables injected by App Center.
app.MapGet("/config", () =>
{
    var mgmtServer = Environment.GetEnvironmentVariable("LEGACY_MANAGEMENT_SERVER");
    var useTls = Environment.GetEnvironmentVariable("LEGACY_USE_TLS");

    if (string.IsNullOrEmpty(mgmtServer))
        return Results.Ok(new ConfigResponse(null));

    var scheme = string.Equals(useTls, "true", StringComparison.OrdinalIgnoreCase) ? "https" : "http";
    var serverUrl = $"{scheme}://{mgmtServer}";
    return Results.Ok(new ConfigResponse(serverUrl));
});

// Creates a session and returns the access token.
// Useful for verifying that credentials and server connectivity are correct before performing other operations.
app.MapPost("/session/create-with-server-config", async (SessionRequest request, IServiceProvider serviceProvider) =>
{
    try
    {
        var session = SessionHelper.CreateSessionWithServerConfigurationProvided(
            serviceProvider, request.ServerUrl, request.UserType, request.Username, request.Password);

        var token = await SessionHelper.WaitForTokenAsync(session);

        return Results.Ok(new SessionResponse(session.Id.ToString(), session.ServerConfiguration.ServerUri.ToString(), token));
    }
    catch (Exception ex)
    {
        return Results.BadRequest(new { error = ex.Message });
    }
});

// Returns all cameras visible to the authenticated user.
// ConfigurationService.Get<Camera>() queries the VMS configuration for all Camera items.
app.MapPost("/cameras", async (SessionRequest request, IServiceProvider serviceProvider) =>
{
    try
    {
        var session = SessionHelper.CreateSessionWithServerConfigurationProvided(
            serviceProvider, request.ServerUrl, request.UserType, request.Username, request.Password);

        await SessionHelper.WaitForTokenAsync(session);

        var cameras = await new ConfigurationService(session).Get<Camera>();

        return Results.Ok(cameras.Select(c => new CameraResponse(c.Id.ToString(), c.Name ?? string.Empty, c.Description ?? string.Empty)));
    }
    catch (Exception ex)
    {
        return Results.BadRequest(new { error = ex.Message });
    }
});

// Updates the Name and/or Description of a single camera.
// Demonstrates the Filter pattern for fetching a specific item by ID, and camera.Save() for persisting changes.
app.MapPost("/cameras/update", async (CameraUpdateRequest request, IServiceProvider serviceProvider) =>
{
    try
    {
        var session = SessionHelper.CreateSessionWithServerConfigurationProvided(
            serviceProvider, request.ServerUrl, request.UserType, request.Username, request.Password);

        await SessionHelper.WaitForTokenAsync(session);

        var configuration = new ConfigurationService(session);

        // Use a Filter to fetch only the target camera rather than loading all cameras.
        var cameras = await configuration.Get<Camera>(
            [new Filter { Field = "Id", Value = request.CameraId, Operator = FilterOperator.Equal }]);

        var camera = cameras.FirstOrDefault()
            ?? throw new Exception($"Camera with ID '{request.CameraId}' not found.");

        if (request.Name is not null) camera.Name = request.Name;
        camera.Description = request.Description ?? "";

        // Save() persists the changes back to the VMS.
        await camera.Save();

        return Results.Ok(new CameraResponse(camera.Id.ToString(), camera.Name ?? string.Empty, camera.Description ?? string.Empty));
    }
    catch (Exception ex)
    {
        return Results.BadRequest(new { error = ex.Message });
    }
});

// Starts live streaming for a camera (if not already started) and returns total bytes received.
app.MapPost("/cameras/stream", async (CameraStreamRequest request, IServiceProvider serviceProvider) =>
{
    try
    {
        if (!Guid.TryParse(request.CameraId, out var cameraGuid))
            return Results.BadRequest(new { error = $"Invalid camera ID: '{request.CameraId}'" });

        if (CameraStreams.ActiveStreams.TryGetValue(request.CameraId, out var existing))
            return Results.Ok(new CameraStreamResponse(request.CameraId, existing.TotalBytesReceived, false));

        var session = SessionHelper.CreateSessionWithServerConfigurationProvided(
            serviceProvider, request.ServerUrl, request.UserType, request.Username, request.Password);

        await SessionHelper.WaitForTokenAsync(session);

        var state = new CameraStreamState();
        var rawSource = new RawSource(session, RawSourceMode.Live, cameraGuid, null);
        rawSource.DataReady += (_, e) => Interlocked.Add(ref state.TotalBytesReceivedField, e.Data?.Length ?? 0);

        state.Source = rawSource;
        CameraStreams.ActiveStreams[request.CameraId] = state;

        _ = rawSource.Start(CancellationToken.None);

        return Results.Ok(new CameraStreamResponse(request.CameraId, 0, true));
    }
    catch (Exception ex)
    {
        return Results.BadRequest(new { error = ex.Message });
    }
});

// Stops an active stream for the specified camera.
app.MapPost("/cameras/stream/stop", (CameraStreamRequest request) =>
{
    if (CameraStreams.ActiveStreams.TryRemove(request.CameraId, out var state))
    {
        state.Source?.Stop();
        return Results.Ok(new { cameraId = request.CameraId, stopped = true });
    }
    return Results.Ok(new { cameraId = request.CameraId, stopped = false });
});

app.Run();

public record SessionRequest(string ServerUrl, string UserType, string? Username, string? Password);
public record ConfigResponse(string? ServerUrl);
public record SessionResponse(string Id, string ServerUri, string Token);
public record CameraResponse(string Id, string Name, string Description);
public record CameraUpdateRequest(string ServerUrl, string UserType, string? Username, string? Password, string CameraId, string? Name, string? Description);
public record CameraStreamRequest(string ServerUrl, string UserType, string? Username, string? Password, string CameraId);
public record CameraStreamResponse(string CameraId, long TotalBytesReceived, bool StreamStarted);

public class CameraStreamState
{
    public long TotalBytesReceivedField;
    public long TotalBytesReceived => Interlocked.Read(ref TotalBytesReceivedField);
    public RawSource? Source { get; set; }
}

public static class CameraStreams
{
    public static readonly ConcurrentDictionary<string, CameraStreamState> ActiveStreams = new();
}