from pydantic import BaseModel, ConfigDict, Field


class FieldErrorDetail(BaseModel):
    field: str
    code: str
    message: str


class Envelope(BaseModel):
    status: str
    code: str | None = None
    message: str | None = None
    trace_id: str | None = Field(default=None, alias="traceId")
    data: dict[str, object] | list[object] | None = None
    errors: list[FieldErrorDetail] | None = None
    status_code: int | None = Field(default=None, alias="statusCode")
    status_message: str | None = Field(default=None, alias="statusMessage")
    error_code: str | None = Field(default=None, alias="errorCode")

    model_config = ConfigDict(alias_generator=lambda s: s, populate_by_name=True)
