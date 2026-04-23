# Web App SDK Core Sample

A sample web application that demonstrates how to create sessions using the VideoOS Platform SDK Core.

## Overview

This sample is an ASP.NET Core web application with a static frontend. It exposes three API endpoints that create a `Session` using the VideoOS Platform SDK Core, and returns the session ID, server URI, and access token to the browser.

The three session creation methods demonstrated are:

| Method | Endpoint | Description |
|--------|----------|-------------|
| Server config provided | `POST /api/session/create-with-server-config` | Creates a session with an explicit server URL and credentials |
| Runtime config | `POST /api/session/create-with-runtime-config` | Creates a session using server config resolved from the App Center runtime environment |
| Runtime config + default user | `POST /api/session/create-with-runtime-config-default-user` | Creates a session using the runtime config and the default App Center user |

All three methods support four user types: `DefaultWindows`, `Windows`, `Basic`, and `External` (access token).

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
                ├── index.html             # Session Manager UI
                ├── app.js                 # Frontend form handling and API calls
                └── styles.css
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

## VideoOS Platform SDK Core