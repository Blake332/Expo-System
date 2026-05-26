from fastapi import Depends, Header, HTTPException, status

from app.config import settings


def require_desk_key(x_desk_key: str | None = Header(default=None, alias="X-Desk-Key")) -> None:
    if not x_desk_key or x_desk_key != settings.desk_api_key:
        raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail="Invalid desk key")


def require_admin_key(x_admin_key: str | None = Header(default=None, alias="X-Admin-Key")) -> None:
    if not x_admin_key or x_admin_key != settings.admin_api_key:
        raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail="Invalid admin key")
