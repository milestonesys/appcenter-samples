# PostgreSQL Sample

A sample application that demonstrates how to use a PostgreSQL database in Runtime Platform. This sample includes a Go-based todo application that connects to a PostgreSQL database to manage tasks.

## Overview

This sample demonstrates:
- Database schema initialization and management
- Basic CRUD operations (Create, Read, Update, Delete) with PostgreSQL

## Sample Structure

```
postgresql/
├── app-definition.yaml          # App Center application definition with database config
├── Makefile                    # Build and deployment commands
├── README.md                   # This documentation
└── containers/
    └── todo/
        ├── Dockerfile          # Go application container build instructions
        └── src/
            ├── go.mod          # Go module dependencies
            └── main.go         # Todo application source code
```

## Database Configuration

The application automatically creates a PostgreSQL database with:
- **Database name**: `samples.todo-db`
- **Table**: `tasks` with columns:
  - `id` (serial primary key)
  - `description` (text, not null)

The database schema is defined in `app-definition.yaml` and automatically initialized when the app is deployed.

## How to Build and Deploy

First use the login command to connect your running cluster.

```bash
cd src/postgresql
make login
make build # builds both the image and the chart
```

### Deploy to App Center
Push the application to your cluster's Sandbox registry and install it:

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

## Using the Todo Application

The sample includes a command-line todo application that demonstrates database operations. When deployed, the application automatically:

1. Adds three sample tasks (`item1`, `item2`, `item3`)
2. Lists all tasks
3. Keeps running for demonstration purposes

## Connecting to the Database with psql

To connect directly to the PostgreSQL database using `psql`:
```bash
# Get DB-related pods 
kubectl get pods -n postgresql
# Access the pods in shell
kubectl exec --stdin --tty pg-cluster-1 -n postgresql -- /bin/bash
```

You can always use `k9s` or `Kubernetes Dashboard` to access the pods in shell mode.

Once you are in shell mode in the Postgresql pod access the `psql` 
```bash
> psql
# List all DBs
> \l
# Output
                                                      List of databases
      Name       |      Owner      | Encoding | Locale Provider | Collate | Ctype | Locale | ICU Rules |   Access privileges
-----------------+-----------------+----------+-----------------+---------+-------+--------+-----------+-----------------------
 samples.todo-db | samples.todo-db | UTF8     | libc            | C       | C     |        |           | 
```

## Database Environment Variables

The application uses the following environment variable automatically provided by App Center:

- `PGDB_URI`: Complete PostgreSQL connection string (includes host, port, username, password, and database name)
