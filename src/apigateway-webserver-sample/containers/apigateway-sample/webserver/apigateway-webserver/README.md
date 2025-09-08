# API Gateway webserver

This app is a golang cmd application. The app when run will start a server and start listening to port 8080.
When the app is run locally and not on the App Center runtime, it will only support login via basic user credentials as the OAuth2 Client Credentials Flow setup won't be available in this environment.

## Requirements

Installation of Golang 1.23.6 compiler. For more information follow the official Golang installation [docs](https://go.dev/doc/install).
