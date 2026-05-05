# Web App SDK Core Sample

A sample web application that demonstrates how to create sessions using the VideoOS Platform SDK Core.

## Overview

This sample is an ASP.NET Core web application with a static frontend. It exposes API endpoints for creating a `Session` using the VideoOS Platform SDK Core, querying cameras, and updating camera metadata.

The implemented endpoints are:

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/config` | Returns the VMS server URL assembled from runtime environment variables |
| POST | `/session/create-with-server-config` | Creates a session with an explicit server URL and credentials |
| POST | `/cameras` | Returns all cameras visible to the authenticated user |
| POST | `/cameras/update` | Updates the `Name` and/or `Description` of a single camera |

The session creation endpoint supports four user types: `DefaultWindows`, `Windows`, `Basic`, and `External` (access token).

## Sample Structure

```
web-app-sdk-core/
├── app-definition.yaml                    # App Center application definition
├── Makefile                               # Build and deployment commands
└── containers/
    └── web-app-sdk-core/
        └── web-app-sdk-core/
            ├── Program.cs                 # API endpoints
            ├── SessionHelper.cs           # SDK session creation logic
            ├── TestWebApp.csproj
            └── wwwroot/
                └── index.html             # Session Manager UI
```

## How to Build and Deploy

First use the login command to connect your running cluster.

```bash
cd src/web-app-sdk-core
make login   # connect to the cluster (only required once)
make build   # build the container image and Helm chart
```

### Deploy to App Center

```bash
make push              # push image and chart to the registry
make install-from-repo # install the application
```

### Verify Deployment

```bash
make list    # list installed applications
make events  # check application events
```

### Accessing the Application

Once deployed, the application is available at:

```
http://<system-ip>/api/samples/sdk-test/
```

Replace `<system-ip>` with your actual cluster domain.

## API Endpoints

This sample demonstrates how to use the VideoOS Platform SDK Core in a containerized application. It provides a simple browser UI and four endpoints covering session creation, configuration queries, and camera updates.

### `GET /config`
Returns the VMS server URL assembled from the `LEGACY_MANAGEMENT_SERVER` and `LEGACY_USE_TLS` runtime environment variables. Used by the UI to pre-fill the server URL on load.

### `POST /session/create-with-server-config`
Creates an SDK session with the supplied credentials and returns the session ID, server URI, and access token. Useful for verifying connectivity before performing other operations.

### `POST /cameras`
Returns all cameras visible to the authenticated user. Uses `ConfigurationService.Get<Camera>()` to query the VMS configuration.

### `POST /cameras/update`
Updates the `Name` and/or `Description` of a single camera identified by its ID. Uses a `Filter` query to fetch the specific camera and calls `camera.Save()` to persist the changes. Fields omitted from the request are left unchanged (partial update).