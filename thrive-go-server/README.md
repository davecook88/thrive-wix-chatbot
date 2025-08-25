# thrive-go-server

Lightweight Go backend used by the Thrive Wix chatbot projects. This README describes how to build, configure and run the server located in this folder.

## Overview

thrive-go-server implements API endpoints and integrations used by the Wix chatbot frontend(s). It contains authentication helpers, ChatGPT client logic, Wix client helpers, handlers for chat requests, and administrative routes.

Key packages and areas:
- `main.go` — server entrypoint
- `auth/` — Wix auth helpers
- `chatgpt/` — ChatGPT client and helpers
- `wix/` — Wix client and type definitions
- `handlers/` — request handlers and tool handlers
- `admin/` — admin routes

## Prerequisites

- Go 1.20+ (recommended)
- A valid `creds.json` (used by some parts of the server, see configuration)
- `app.yaml` if you plan to deploy to App Engine (this repo includes one)

Verify Go is installed:

```bash
go version
```

## Configuration

The repo uses a small set of files and environment variables for configuration.

- `creds.json` — local credential file used by parts of the server. Keep this out of source control for production.
- `app.yaml` — App Engine configuration included for deployment.

Environment variables (common):

- `PORT` — port to run the server on (defaults to 8080 in many setups)
- `OPENAI_API_KEY` — API key used by the ChatGPT client (if applicable)

Check the code for other environment variables used by custom integrations.

## Build & Run (local)

From the `thrive-go-server` directory:

```bash
# build
go build ./...

# run directly
go run ./main.go

# or build a binary and run it
go build -o thrive-server ./
./thrive-server
```

You can also use your preferred live-reload tool (the repository previously used `air` in development).

## Endpoints (high level)

The server exposes multiple HTTP routes across different packages. Major responsibilities:

- Authentication endpoints for validating Wix members and tokens (`auth/`, `chrome-ext/`)
- Chat request handlers that take user input and forward to ChatGPT-related logic (`handlers/chatgpt-requests.go`)
- Admin routes for internal tooling and diagnostics (`admin/routes.go`)

For exact paths and handler signatures, inspect the `handlers/` and `admin/` packages.

## Development notes

- Follow standard Go module workflows. The module is declared in `go.mod` at the repo root.
- Add unit tests alongside packages and run `go test ./...` from the `thrive-go-server` root.
- Keep secrets out of the git history. Use local env files or secret managers for production.

## Testing

Run unit tests:

```bash
go test ./... -v
```

## Deployment

This project includes an `app.yaml` for Google App Engine. To deploy with App Engine:

1. Ensure `gcloud` CLI is configured with the correct project.
2. Set any required environment variables in the App Engine configuration or via Secret Manager.
3. Deploy:

```bash
gcloud app deploy app.yaml
```

Adjust deployment steps for other hosting providers as needed.

## Troubleshooting

- Build errors: run `go build ./...` and inspect the output. Missing dependencies will appear at build time.
- Runtime errors: check the server logs (stdout/stderr). If using App Engine, check the Logs Viewer.
