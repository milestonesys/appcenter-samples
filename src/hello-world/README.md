# Hello World Sample

A simple web application that demonstrates the basic structure and deployment of an App Center application. This sample serves a static "Hello World" HTML page using Apache HTTP Server.

## Overview

This is the simplest possible App Center application that:
- Uses a standard Apache HTTP Server container
- Serves a static HTML page
- Demonstrates the basic app structure required for App Center deployment

## Sample Structure

```
hello-world/
├── app-definition.yaml          # App Center application definition
├── Makefile                    # Build and deployment commands
├── README.md                   # This documentation
└── containers/
    └── hello-world/
        ├── Dockerfile          # Container build instructions
        └── html/
            └── index.html      # Static HTML content
```

## How to Build and Deploy

First use the login command to connect your running cluster.

```bash
cd src/hello-world
# Login is only required once to connect to the cluster
make login # “logging in” to the cluster from your machine 
make build # builds both the image and the chart
```

### Deploy to App Center
Push the application to your cluster's `sandbox` registry and install it:

```bash
# Push to registry and repository
make push # (helm and docker image)

# Install the application (similar to install button in the UI)
make install-from-repo
```

### Verify Deployment
Check that the application is running:

```bash
# List installed applications
make list

# Check application events
make events
```

### Accessing the Application

Once deployed, the application will be available at:
```
http://<your-cluster-domain>/api/samples/hello-world/
```

Replace `<your-cluster-domain>` with your actual cluster domain.
