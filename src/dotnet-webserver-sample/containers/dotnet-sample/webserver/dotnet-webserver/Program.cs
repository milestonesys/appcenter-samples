using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Http;
using Microsoft.Data.Sqlite;
using Microsoft.OpenApi.Models;
using Swashbuckle.AspNetCore.Annotations;
using Microsoft.AspNetCore.Http.HttpResults;
using System.Collections.Generic;
using System.IO;
using System.Text.Json;
using System.Threading.Tasks;

// Top-level statements - Start here

var builder = WebApplication.CreateBuilder(args);

// Add services for Swagger
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen(options =>
{
    options.SwaggerDoc("v1", new OpenApiInfo
    {
        Title = "Key-Value Store API",
        Version = "v1",
        Description = "A simple key-value store API using .NET and SQLite"
    });
    options.EnableAnnotations(); // Enable annotations
});

var app = builder.Build();

// Enable Swagger and Swagger UI
app.UseSwagger();
app.UseSwaggerUI(options =>
{
    if (app.Environment.IsDevelopment())
    {
        Console.WriteLine("Swagger UI in dev mode");
        options.SwaggerEndpoint("/swagger/v1/swagger.json", "Key-Value Store API v1");
    } else {
        Console.WriteLine("Swagger UI in production mode");
        options.SwaggerEndpoint("/dotnet-sample/swagger/v1/swagger.json", "Key-Value Store API v1");
        options.RoutePrefix = "swagger"; // Set Swagger UI at /swagger
    }
});

// SQLite Connection
string connectionString = "Data Source=/app/data.db";
using (var connection = new SqliteConnection(connectionString))
{
    connection.Open();
    var command = connection.CreateCommand();
    command.CommandText = "CREATE TABLE IF NOT EXISTS KeyValueStore (Key TEXT PRIMARY KEY, Value TEXT)";
    command.ExecuteNonQuery();
}

// Get all key-value pairs
app.MapGet("/items", async () =>
{
    var items = new List<object>();
    using var connection = new SqliteConnection(connectionString);
    await connection.OpenAsync();

    var command = connection.CreateCommand();
    command.CommandText = "SELECT Key, Value FROM KeyValueStore";
    using var reader = await command.ExecuteReaderAsync();
    while (await reader.ReadAsync())
    {
        items.Add(new { Key = reader.GetString(0), Value = reader.GetString(1) });
    }

    return Results.Json(items);
})
.WithName("GetAllItems")
.WithMetadata(new SwaggerOperationAttribute("Get all items", "Returns all stored key-value pairs"));

// Add a key-value pair
app.MapPost("/items", async (KeyValueRequest request) =>
{
    using var connection = new SqliteConnection(connectionString);
    await connection.OpenAsync();

    var command = connection.CreateCommand();
    command.CommandText = "INSERT INTO KeyValueStore (Key, Value) VALUES (@key, @value)";
    command.Parameters.AddWithValue("@key", request.Key);
    command.Parameters.AddWithValue("@value", request.Value);

    try
    {
        await command.ExecuteNonQueryAsync();
        return Results.Created($"/items/{request.Key}", new { Key = request.Key, Value = request.Value });
    }
    catch
    {
        return Results.BadRequest("Key already exists.");
    }
})
.WithName("AddItem")
.WithMetadata(new SwaggerOperationAttribute("Add a new key-value pair", "Creates a new key-value pair"));

// Update a value
app.MapPut("/items/{key}", async (string key, KeyValueRequest request) =>
{
    using var connection = new SqliteConnection(connectionString);
    await connection.OpenAsync();

    var command = connection.CreateCommand();
    command.CommandText = "UPDATE KeyValueStore SET Value = @value WHERE Key = @key";
    command.Parameters.AddWithValue("@value", request.Value);
    command.Parameters.AddWithValue("@key", key);

    int rowsAffected = await command.ExecuteNonQueryAsync();
    return rowsAffected > 0 ? Results.Ok(new { Key = key, Value = request.Value }) : Results.NotFound();
})
.WithName("UpdateItem")
.WithMetadata(new SwaggerOperationAttribute("Update a value for a given key", "Updates the value of an existing key-value pair"));

// Delete a key
app.MapDelete("/items/{key}", async (string key) =>
{
    using var connection = new SqliteConnection(connectionString);
    await connection.OpenAsync();

    var command = connection.CreateCommand();
    command.CommandText = "DELETE FROM KeyValueStore WHERE Key = @key";
    command.Parameters.AddWithValue("@key", key);

    int rowsAffected = await command.ExecuteNonQueryAsync();
    return rowsAffected > 0 ? Results.Ok($"Deleted {key}") : Results.NotFound();
})
.WithName("DeleteItem")
.WithMetadata(new SwaggerOperationAttribute("Delete an item", "Deletes a key-value pair based on the provided key"));

app.Urls.Add("http://0.0.0.0:80");
app.Run();

// Define classes below the top-level statements (after the WebApplication setup)
public class KeyValueRequest
{
    public string Key { get; set; }
    public string Value { get; set; }
}
