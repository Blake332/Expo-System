from fastapi import APIRouter, Depends
from sqlalchemy.ext.asyncio import AsyncSession

from app.database import get_db
from app.schemas import PublicStatus
from app.services import public_status

router = APIRouter(prefix="/api/public", tags=["public"])


@router.get("/status", response_model=PublicStatus)
async def get_status(db: AsyncSession = Depends(get_db)) -> PublicStatus:
    data = await public_status(db)
    return PublicStatus(**data)
