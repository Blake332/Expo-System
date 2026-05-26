from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    model_config = SettingsConfigDict(env_file=".env", env_file_encoding="utf-8", extra="ignore")

    desk_api_key: str = "dev-desk-key"
    admin_api_key: str = "dev-admin-key"
    database_url: str = "sqlite+aiosqlite:///./expo_registration.db"


settings = Settings()
