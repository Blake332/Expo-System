# Python FastAPI Backend

FastAPI-based REST API for the Expo Registration System with support for admin, desk, and public registration flows.

## 🛠️ Tech Stack

- **Framework:** FastAPI
- **Server:** Uvicorn
- **Database:** SQLite with SQLAlchemy ORM
- **Validation:** Pydantic v2
- **Python:** 3.11+

## 📋 Features

- RESTful API endpoints
- Role-based access (admin, desk, public)
- Email validation
- Async/await support
- Automatic API documentation (Swagger UI)

## 🚀 Getting Started

### Installation

```bash
cd services/python
pip install -r requirements.txt
```

### Configuration

1. Copy the environment file:
```bash
cp ../../.env.example .env
```

2. Update `.env`:
```
DESK_API_KEY=your-desk-key
ADMIN_API_KEY=your-admin-key
DATABASE_URL=sqlite+aiosqlite:///./expo_registration.db
```

### Running Locally

```bash
python -m uvicorn app.main:app --reload
```

The API will be available at `http://localhost:8000`

### API Documentation

- **Swagger UI:** http://localhost:8000/docs
- **ReDoc:** http://localhost:8000/redoc

## 📁 Project Structure

```
app/
├── main.py           # Application entry point
├── config.py         # Configuration management
├── database.py       # Database setup
├── models.py         # SQLAlchemy models
├── schemas.py        # Pydantic schemas
├── services.py       # Business logic
├── deps.py           # Dependency injection
├── __init__.py
└── routers/
    ├── admin.py      # Admin endpoints
    ├── desk.py       # Desk endpoints
    ├── public.py     # Public endpoints
    └── __init__.py
```

## 🔌 API Endpoints

### Public
- `GET /` - Welcome message

### Desk
- `POST /desk/register` - Register via desk
- `GET /desk/registrations` - List registrations

### Admin
- `GET /admin/dashboard` - Admin dashboard
- `POST /admin/verify` - Verify registration
- `DELETE /admin/registration/{id}` - Delete registration

## 🗄️ Database

SQLite database: `expo_registration.db`

### Models
- User
- Registration
- Event
- Verification

## 📦 Requirements

See `requirements.txt` for all dependencies.

## 🧪 Testing

(Add testing setup as needed)

```bash
pip install pytest pytest-asyncio
pytest
```

## 🐳 Docker

Build and run with Docker:
```bash
docker build -t expo-python-backend .
docker run -p 8000:8000 expo-python-backend
```

## 📝 Notes

- Use `.env` for local configuration (not tracked in git)
- Database migrations should be managed via SQLAlchemy
- Keep API endpoints RESTful and well-documented

