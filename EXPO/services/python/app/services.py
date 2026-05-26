from datetime import datetime, timezone
import io
import openpyxlfrom sqlalchemy import func, select
from sqlalchemy.ext.asyncio import AsyncSession

from app.models import EventConfig, Registrant
from app.schemas import RegCategory

MAX_REGISTRANTS = 10_000


async def ensure_event_config(session: AsyncSession) -> EventConfig:
    r = await session.execute(select(EventConfig).order_by(EventConfig.id).limit(1))
    row = r.scalar_one_or_none()
    if row is None:
        row = EventConfig(
            countdown_target_utc=None,
            registration_open=True,
            public_message=None,
        )
        session.add(row)
        await session.commit()
        await session.refresh(row)
    return row


async def public_status(session: AsyncSession) -> dict:
    cfg = await ensure_event_config(session)
    total = await session.scalar(select(func.count()).select_from(Registrant)) or 0
    counts: dict[str, int] = {c.value: 0 for c in RegCategory}
    r = await session.execute(
        select(Registrant.category, func.count()).group_by(Registrant.category)
    )
    for cat, n in r.all():
        counts[str(cat)] = int(n)
    return {
        "countdown_target_utc": cfg.countdown_target_utc,
        "registration_open": cfg.registration_open,
        "public_message": cfg.public_message,
        "total_registered": int(total),
        "count_by_category": counts,
    }


async def register_person(
    session: AsyncSession,
    *,
    full_name: str,
    email: str | None,
    phone: str | None,
    category: RegCategory,
    notes: str | None,
) -> Registrant:
    cfg = await ensure_event_config(session)
    if not cfg.registration_open:
        raise ValueError("registration_closed")
    total = await session.scalar(select(func.count()).select_from(Registrant)) or 0
    if int(total) >= MAX_REGISTRANTS:
        raise ValueError("capacity_reached")
    p = Registrant(
        full_name=full_name.strip(),
        email=email,
        phone=phone,
        category=category.value,
        notes=notes,
    )
    session.add(p)
    await session.commit()
    await session.refresh(p)
    return p


async def lookup_by_name(session: AsyncSession, q: str, limit: int = 50) -> list[Registrant]:
    q = (q or "").strip()
    if len(q) < 2:
        return []
    like = f"%{q.lower()}%"
    stmt = (
        select(Registrant)
        .where(func.lower(Registrant.full_name).like(like))
        .order_by(Registrant.full_name.asc())
        .limit(limit)
    )
    r = await session.execute(stmt)
    return list(r.scalars().all())


async def admin_update_config(session: AsyncSession, updates: dict) -> EventConfig:
    cfg = await ensure_event_config(session)
    for key in ("countdown_target_utc", "registration_open", "public_message"):
        if key in updates:
            setattr(cfg, key, updates[key])
    cfg.updated_at = datetime.now(timezone.utc)
    await session.commit()
    await session.refresh(cfg)
    return cfg


async def admin_stats(session: AsyncSession) -> dict:
    base = await public_status(session)
    return {
        "total_registered": base["total_registered"],
        "count_by_category": base["count_by_category"],
        "countdown_target_utc": base["countdown_target_utc"],
        "registration_open": base["registration_open"],
        "public_message": base["public_message"],
        "total_checked_in": base.get("total_checked_in", 0),
    }

async def export_to_excel(session: AsyncSession) -> io.BytesIO:
    stmt = select(Registrant).order_by(Registrant.created_at.desc())
    r = await session.execute(stmt)
    registrants = r.scalars().all()

    wb = openpyxl.Workbook()
    ws = wb.active
    ws.title = "Registrants"
    
    headers = ["ID", "Full Name", "Company", "Email", "Phone", "Category", "Checked In", "Notes", "Registration Date"]
    ws.append(headers)
    
    for reg in registrants:
        ws.append([
            reg.id,
            reg.full_name,
            reg.company or "",
            reg.email or "",
            reg.phone or "",
            reg.category,
            "Yes" if reg.checked_in else "No",
            reg.notes or "",
            reg.created_at.strftime("%Y-%m-%d %H:%M:%S") if reg.created_at else ""
        ])
    
    stream = io.BytesIO()
    wb.save(stream)
    stream.seek(0)
    return stream
