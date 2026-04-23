from __future__ import annotations

import pytest

from approute.models import Product, ProductListResponse, ProductStockResponse
from approute.resources.services import ServicesResource
from tests.conftest import MockTransport
from tests.factories import make_product, make_product_list_response, make_product_stock_response


@pytest.fixture
def resource(transport: MockTransport) -> ServicesResource:
    return ServicesResource(transport)  # type: ignore[arg-type]


class TestServicesList:
    def test_returns_product_list(
        self, resource: ServicesResource, transport: MockTransport
    ) -> None:
        transport.set_response(make_product_list_response())
        result = resource.list()
        assert isinstance(result, ProductListResponse)
        assert len(result.items) == 1
        assert result.items[0].id == "prod-001"

    def test_calls_correct_endpoint(
        self, resource: ServicesResource, transport: MockTransport
    ) -> None:
        transport.set_response(make_product_list_response())
        resource.list()
        assert transport.last_call.method == "GET"
        assert transport.last_call.path == "/services"


class TestServicesGet:
    def test_returns_product(self, resource: ServicesResource, transport: MockTransport) -> None:
        transport.set_response(make_product())
        result = resource.get("prod-001")
        assert isinstance(result, Product)
        assert result.id == "prod-001"

    def test_calls_correct_endpoint(
        self, resource: ServicesResource, transport: MockTransport
    ) -> None:
        transport.set_response(make_product())
        resource.get("prod-001")
        assert transport.last_call.method == "GET"
        assert transport.last_call.path == "/services/prod-001"


class TestServicesStock:
    def test_returns_stock_response(
        self, resource: ServicesResource, transport: MockTransport
    ) -> None:
        transport.set_response(make_product_stock_response())
        result = resource.stock("prod-001")
        assert isinstance(result, ProductStockResponse)
        assert result.product_id == "prod-001"
        assert len(result.items) == 2

    def test_calls_correct_endpoint(
        self, resource: ServicesResource, transport: MockTransport
    ) -> None:
        transport.set_response(make_product_stock_response())
        resource.stock("prod-001")
        assert transport.last_call.method == "GET"
        assert transport.last_call.path == "/services/prod-001/stock"
