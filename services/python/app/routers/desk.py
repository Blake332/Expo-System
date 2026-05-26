from fastapi import APIRouter, Depends, HTTPException, Query, status
from sqlalchemy.ext.asyncio import AsyncSession

from app.database import get_db
from app.deps import require_desk_key
from app.schemas import RegisterIn, RegistrantOut
from app.services import lookup_by_name, register_person

router = APIRouter(prefix="/api/desk", tags=["desk"], dependencies=[Depends(require_desk_key)])


@router.post("/register", response_model=RegistrantOut)
async def desk_register(body: RegisterIn, db: AsyncSession = Depends(get_db)) -> RegistrantOut:
    try:
        p = await register_person(
            db,
            full_name=body.full_name,
            email=body.email,
            phone=body.phone,
            category=body.category,
            notes=body.notes,
        )
    except ValueError as e:
        if str(e) == "registration_closed":
            raise HTTPException(status_code=status.HTTP_403_FORBIDDEN, detail="Registration is closed")
        if str(e) == "capacity_reached":
            raise HTTPException(status_code=status.HTTP_409_CONFLICT, detail="Maximum capacity reached")
        raise
    return RegistrantOut.model_validate(p)


@router.get("/lookup", response_model=list[RegistrantOut])
async def desk_lookup(
    q: str = Query(..., min_length=2, max_length=200),
    db: AsyncSession = Depends(get_db),
) -> list[RegistrantOut]:
    rows = await lookup_by_name(db, q)
    return [RegistrantOut.model_validate(r) for r in rows]
