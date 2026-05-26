from datetime import datetime

from sqlalchemy import Boolean, DateTime, Index, Integer, String, Text, func
from sqlalchemy.orm import Mapped, mapped_column

from app.database import Base


class EventConfig(Base):
    __tablename__ = "event_config"

    id: Mapped[int] = mapped_column(Integer, primary_key=True, autoincrement=True)
    countdown_target_utc: Mapped[datetime | None] = mapped_column(DateTime(timezone=True), nullable=True)
    registration_open: Mapped[bool] = mapped_column(Boolean, default=True, nullable=False)
    public_message: Mapped[str | None] = mapped_column(Text, nullable=True)
    updated_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), server_default=func.now(), nullable=False)


class Registrant(Base):
    __tablename__ = "registrants"

    id: Mapped[int] = mapped_column(Integer, primary_key=True, autoincrement=True)
    full_name: Mapped[str] = mapped_column(String(200), nullable=False)
    email: Mapped[str | None] = mapped_column(String(255), nullable=True)
    phone: Mapped[str | None] = mapped_column(String(64), nullable=True)
    company: Mapped[str | None] = mapped_column(String(255), nullable=True)
    # Three expo categories: general, trade, press (labels editable in UI copy only)
    category: Mapped[str] = mapped_column(String(32), nullable=False, index=True)
    notes: Mapped[str | None] = mapped_column(Text, nullable=True)
    checked_in: Mapped[bool] = mapped_column(Boolean, default=False, nullable=False)
    created_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), server_default=func.now(), nullable=False)


Index("ix_registrants_full_name", Registrant.full_name)
Index("ix_registrants_created", Registrant.created_at)
