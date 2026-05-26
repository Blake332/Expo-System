from fastapi import APIRouter, Depends
from sqlalchemy.ext.asyncio import AsyncSession

from app.database import get_db
from app.deps import require_admin_key
from app.schemas import AdminConfigUpdate, AdminStats
from app.services import admin_stats, admin_update_config, public_status

router = APIRouter(prefix="/api/admin", tags=["admin"], dependencies=[Depends(require_admin_key)])


@router.get("/stats", response_model=AdminStats)
async def admin_get_stats(db: AsyncSession = Depends(get_db)) -> AdminStats:
    data = await admin_stats(db)
    return AdminStats(**data)


@router.put("/config", response_model=AdminStats)
async def admin_put_config(body: AdminConfigUpdate, db: AsyncSession = Depends(get_db)) -> AdminStats:
    payload = body.model_dump(exclude_unset=True)
    await admin_update_config(db, payload)
    data = await admin_stats(db)
    return AdminStats(**data)


@router.get("/status", response_model=AdminStats)
async def admin_get_status(db: AsyncSession = Depends(get_db)) -> AdminStats:
    data = await public_status(db)
    return AdminStats(**data)
