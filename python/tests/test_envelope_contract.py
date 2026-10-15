"""Literal Public API wire responses, shared with the other three SDK suites."""
from __future__ import annotations

import json
from dataclasses import asdict
from pathlib import Path
from typing import Any

import httpx
import pytest

from approute import errors
from approute._transport._base import handle_response
from approute.models import Envelope

CASES = json.loads((Path(__file__).parents[2] / "testdata/envelopes.json").read_text())


def test_exported_envelope_model_accepts_current_wire_shape() -> None:
    model = Envelope.model_validate(CASES[0]["body"])
    assert model.status_code == 0
    assert model.status_message == "OK"
    assert model.trace_id == "wire-0"


@pytest.mark.parametrize("case", CASES, ids=[case["name"] for case in CASES])
def test_public_api_wire_envelope(case: dict[str, Any]) -> None:
    response = httpx.Response(
        case["httpStatus"], content=json.dumps(case["body"]),
        headers={"Content-Type": "application/json"},
    )

    if case.get("invalid"):
        with pytest.raises(errors.AppRouteError, match="Invalid.*envelope"):
            handle_response(response)
    elif "data" in case:
        assert handle_response(response) == case["data"]
    else:
        error_type = getattr(errors, case["errorType"] + "Error")
        with pytest.raises(error_type) as caught:
            handle_response(response)
        error = caught.value
        assert error.code == case["code"]
        assert error.message == case["message"]
        assert error.trace_id == case["body"]["traceId"]
        assert error.status_code == case["httpStatus"]
        assert error.api_status_code == case.get("apiStatusCode")
        assert error.api_error_code == case.get("apiErrorCode")
        assert [asdict(detail) for detail in error.errors] == case.get("errors", [])
