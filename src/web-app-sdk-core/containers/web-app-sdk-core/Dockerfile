FROM mcr.microsoft.com/dotnet/sdk:10.0 AS build
WORKDIR /App

# Copy the web app source
COPY web-app-sdk-core ./web-app-sdk-core

WORKDIR /App/web-app-sdk-core

RUN dotnet restore
RUN dotnet publish -o /App/out

FROM mcr.microsoft.com/dotnet/aspnet:10.0
WORKDIR /App
COPY --from=build /App/out .
EXPOSE 8080
ENV ASPNETCORE_URLS=http://+:8080
ENTRYPOINT ["dotnet", "TestWebApp.dll"]
