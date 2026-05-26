from datetime import datetime
from enum import Enum

from pydantic import BaseModel, EmailStr, Field


class RegCategory(str, Enum):
    general = "general"
    trade = "trade"
    press = "press"


class PublicStatus(BaseModel):
    countdown_target_utc: datetime | None
    registration_open: bool
    public_message: str | None
    total_registered: int
    count_by_category: dict[str, int]


class RegisterIn(BaseModel):
    full_name: str = Field(..., min_length=1, max_length=200)
    email: EmailStr | None = None
    phone: str | None = Field(None, max_length=64)
    company: str | None = Field(None, max_length=255)
    category: RegCategory
    notes: str | None = None


class RegistrantOut(BaseModel):
    id: int
    full_name: str
    email: str | None
    phone: str | None
    company: str | None
    category: str
    notes: str | None
    checked_in: bool
    created_at: datetime

    model_config = {"from_attributes": True}


class AdminConfigUpdate(BaseModel):
    countdown_target_utc: datetime | None = None
    registration_open: bool | None = None
    public_message: str | None = Field(None, max_length=2000)


class AdminStats(BaseModel):
    total_registered: int
    total_checked_in: int
    count_by_category: dict[str, int]
    countdown_target_utc: datetime | None
    registration_open: bool
    public_message: str | None
