# RESTful Events API Sample

This sample shows how to trigger and retrieve events from a Python application. The sample logs into the server, retrieves any existing user-defined event, triggers it, retrieves the event by id and retrieves all stored events with their metadata. The sample uses [Streamlit](https://streamlit.io/) for creating a web interface.

![RESTful Events API Sample](./../../img/restful-events-api.png)

## Prerequisites

- XProtect installation
- The API Gateway installed on the same host as the management server
- An existing user defined event with event type retention policy greater than 0 days.
- Python version 3.10 or newer

## Environment Setup

### Environment Variables

This sample is designed to run in a container where the following environment variables are provided by the container runtime:

| Variable | Description | Example |
|----------|-------------|---------|
| `LEGACY_USE_TLS` | Whether XProtect server uses TLS | `true` or `false` |
| `LEGACY_MANAGEMENT_SERVER` | Management Server hostname | `your-management-server.com` |
| `CCF_CLIENT_ID` | The client ID to use for authenticating | `validClientId` |
| `CCF_CLIENT_SECRET` | The client secret to use for authenticating | `validClientSecret` |

For local development, set these variables in the `.env` file.

### Local Development Setup

1. **Copy environment configuration:**
   ```bash
   # Linux/macOS/Git Bash
   cp .env.example .env
   ```
   ```powershell
   # Windows PowerShell
   Copy-Item .env.example .env
   ```
   ```cmd
   # Windows Command Prompt
   copy .env.example .env
   ```

2. **Configure your environment:**
   Edit the `.env` file with your XProtect server details:
   ```env
   LEGACY_USE_TLS=true
   LEGACY_MANAGEMENT_SERVER=your-management-server.com
   CCF_CLIENT_ID=validClientId
   CCF_CLIENT_SECRET=validClientSecret
   ```

3. **Install dependencies:**
   ```bash
   cd containers/restful-events-api
   pip install -r requirements.txt
   ```

4. **Run the application:**
   ```bash
   streamlit run events_rest_api.py
   ```

### Container Deployment

**Build the Docker image:**
```bash
docker build -t restful-events-api:1.0.0 -f containers/Dockerfile .
```

**Run the container:**
```bash
docker run -p 8501:8501 \
  -e LEGACY_USE_TLS=true \
  -e LEGACY_MANAGEMENT_SERVER=your-management-server.com \
  -e CCF_CLIENT_ID=validClientId \
  -e CCF_CLIENT_SECRET=validClientSecret \
  restful-events-api:1.0.0
```

### App Center Deployment

Use the make commands described in the top-level README file to build and manage the sample.

When deployed through App Center, environment variables are automatically configured and should not be set manually. Remove the `.env` file before building and pushing your sample application to App Center.

## Project Structure

```
restful-events-api/
├── .env.example                    # Environment variables template, only needed for local development
├── app-definition.yaml             # App Center configuration
├── Makefile                        # Build automation
├── containers/
│   ├── Dockerfile                  # Container build instructions
│   └── restful-events-api/
│       ├── requirements.txt        # Python dependencies
│       ├── events_rest_api.py      # Main application code
│       └── identity_provider.py    # Authentication handler
└── README.md                       # This file
```

## Porting Existing Python MIP Integration

This sample was ported from the existing protocol integration sample [EventsRestApiPython](https://github.com/milestonesys/mipsdk-samples-protocol/tree/main/EventsRestApiPython).

### Steps to Port Your Python Integration:

1. **Create requirements file:** List all packages needed for pip installation based on your sample documentation

2. **Add Streamlit web interface:** 
   - Convert console application to web app using Streamlit
   - Add `streamlit` to your requirements.txt

3. **Test locally:**
   - Verify your sample works as a web application
   - Test all functionality through the web interface

4. **Containerize:**
   - Create Dockerfile following this sample's pattern
   - Test that your sample works in Docker container locally
   - Reference: [Streamlit Docker deployment guide](https://docs.streamlit.io/deploy/tutorials/docker)

5. **Deploy to App Center:**
   - Once container testing is successful, you can build, push, and install your application in the App Center