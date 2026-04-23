from approute.models._base import _Base


class ProductFieldOption(_Base):
    label: str
    value: str
    price: float | None = None


class ProductFieldValidation(_Base):
    min: float | None = None
    max: float | None = None
    pattern: str | None = None
    message: str | None = None


class ProductField(_Base):
    key: str
    type: str
    required: bool
    label: str | None = None
    options: list[ProductFieldOption] | None = None
    validation: ProductFieldValidation | None = None


class ProductItem(_Base):
    id: str
    name: str | None = None
    nominal: float | None = None
    price: float
    currency: str
    available: bool | None = None
    stock: int | None = None
    is_long_order: bool | None = None
    min_qty_to_long_order: int | None = None


class Product(_Base):
    id: str
    name: str | None = None
    type: str
    image_url: str | None = None
    country_code: str | None = None
    category_name: str | None = None
    subcategory_name: str | None = None
    items: list[ProductItem] = []
    fields: list[ProductField] | None = None


class ProductListResponse(_Base):
    items: list[Product]
    has_next: bool = False


class ProductStockItem(_Base):
    item_id: str
    stock: int | None = None


class ProductStockResponse(_Base):
    product_id: str
    items: list[ProductStockItem]
