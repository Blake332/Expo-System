Go backend replacement for the original Python FastAPI app.

Quick start:

1. Install Go (1.20+).
2. From this folder run:

```bash
go mod tidy
go run .
```

The server listens on `:8080`, serves static files from `../static`, and exposes the same `/api/*` endpoints.

Set `DESK_API_KEY` and `ADMIN_API_KEY` environment variables to secure the desk/admin routes.
