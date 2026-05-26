# Expo Registration System

A full-stack monorepo application with FastAPI (Python) and Go backends, featuring admin, desk, and public registration interfaces.

## 📁 Project Structure

```
.
├── frontend/              # UI static files
│   └── static/
│       ├── admin.html     # Admin dashboard
│       ├── desk.html      # Desk registration
│       └── public.html    # Public registration
├── services/
│   ├── python/            # FastAPI backend
│   └── go/                # Go backend
├── .env.example           # Environment template
└── docker-compose.yml     # Local development setup
```

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose (recommended)
- Python 3.11+ (for local Python dev)
- Go 1.21+ (for local Go dev)

### With Docker
```bash
docker-compose up
```

### Without Docker

**Python backend:**
```bash
cd services/python
pip install -r requirements.txt
python -m uvicorn app.main:app --reload
```

**Go backend:**
```bash
cd services/go
go run main.go
```

**Frontend:**
Open `frontend/static/public.html` in your browser.

## 📝 Environment Setup

1. Copy the environment template:
```bash
cp .env.example .env
```

2. Update `.env` with your configuration:
```
DESK_API_KEY=your-desk-key
ADMIN_API_KEY=your-admin-key
DATABASE_URL=sqlite+aiosqlite:///./expo_registration.db
```

## 📚 Documentation

- [Python Backend](services/python/README.md)
- [Go Backend](services/go/README.md)
- [API Documentation](#) - Available at `http://localhost:8000/docs`

## 🛠️ Development

### Python Backend
- Framework: FastAPI
- Database: SQLite with SQLAlchemy ORM
- See [services/python/README.md](services/python/README.md) for details

### Go Backend
- See [services/go/README.md](services/go/README.md) for details

### Frontend
Static HTML files served separately. Communicates with backend APIs.

## 📦 Dependencies

All dependencies are locked in version files:
- Python: `services/python/requirements.txt`
- Go: `services/go/go.mod`

## 📧 License

[Add your license here]
