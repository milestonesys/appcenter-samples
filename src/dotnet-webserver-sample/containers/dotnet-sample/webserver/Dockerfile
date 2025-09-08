# Use official .NET SDK image to build the application
FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build

WORKDIR /app

# Copy the project and restore dependencies
COPY ./dotnet-webserver/ ./
RUN dotnet publish -c Release -o /out

# Use a runtime image for a smaller final container
FROM mcr.microsoft.com/dotnet/aspnet:8.0

WORKDIR /app
COPY --from=build /out .

# Expose the port
EXPOSE 80

# Command to run the application
ENTRYPOINT ["dotnet", "dotnet-webserver.dll"]
