# Web App SDK Core Sample

A sample web application that demonstrates how to create sessions using the VideoOS Platform SDK Core.

## Overview

## Sample Structure

```
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
