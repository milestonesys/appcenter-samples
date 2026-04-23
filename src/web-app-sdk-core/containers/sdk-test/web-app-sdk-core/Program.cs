using TestWebApp;

var builder = WebApplication.CreateBuilder(args);
var app = builder.Build();

// Enable static file serving
app.UseStaticFiles();

// Redirect root to index.html
app.MapGet("", () => Results.Redirect("/index.html"));

app.MapPost("/api/session/create-with-server-config", (SessionRequest request, IServiceProvider serviceProvider) =>
{
    try
    {
        var session = SessionHelper.CreateSessionWithServerConfigurationProvided(
            serviceProvider, 
            request.ServerUrl, 
            request.UserType, 
            request.Username, 
            request.Password);
            
        
        return Results.Ok(new SessionResponse(
            session.Id.ToString(), 
            session.ServerConfiguration.ServerUri.ToString(),
            session.MipTokenCache.Token
        ));
    }
    catch (Exception ex)
    {
        return Results.BadRequest(new { error = ex.Message });
    }
});

app.MapPost("/api/session/create-with-runtime-config", (SessionRuntimeRequest request, IServiceProvider serviceProvider) =>
{
    try
    {
        var session = SessionHelper.CreateSessionWithRuntimeConfiguration(serviceProvider,
            request.UserType, 
            request.Username, 
            request.Password);
        
        return Results.Ok(new SessionResponse(
            session.Id.ToString(), 
            session.ServerConfiguration.ServerUri.ToString(),
            session.MipTokenCache.Token
        ));
    }
    catch (Exception ex)
    {
        return Results.BadRequest(new { error = ex.Message });
    }
});

app.MapPost("/api/session/create-with-runtime-config-default-user", (IServiceProvider serviceProvider) =>
{
    try
    {
        var session = SessionHelper.CreateSessionWithRuntimeConfigurationDefaultUser(serviceProvider);
        
        return Results.Ok(new SessionResponse(
            session.Id.ToString(), 
            session.ServerConfiguration.ServerUri.ToString(),
            session.MipTokenCache.Token
        ));
    }
    catch (Exception ex)
    {
        return Results.BadRequest(new { error = ex.Message });
    }
});

app.Run();

public record SessionRequest(string ServerUrl, string UserType, string Username, string Password);
public record SessionRuntimeRequest(string UserType, string Username, string Password);
public record SessionResponse(string Id, string ServerUri, string Token);