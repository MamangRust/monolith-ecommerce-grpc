# Order domain E2E suite (monolith-ecommerce-grpc).
# Self-contained: creates category + merchant + product, then places an order.
# NOTE: order create requires a full shipping_address object.
# Run with:
#   hurl --test --jobs 1 --variable baseUrl=http://127.0.0.1:5000 \
#     --variable testEmail=... --variable testPassword=... \
#     --variable testImage=... tests/hurl/order.hurl

# --- Auth bootstrap ---------------------------------------------------------
POST {{baseUrl}}/api/auth/register
Content-Type: application/json

{
    "firstname": "Hurl",
    "lastname": "Order",
    "email": "{{testEmail}}",
    "password": "{{testPassword}}",
    "confirm_password": "{{testPassword}}"
}
HTTP 201
[Captures]
userId: jsonpath "$.data.id"

POST {{baseUrl}}/api/auth/login
Content-Type: application/json

{
    "email": "{{testEmail}}",
    "password": "{{testPassword}}"
}
HTTP 200
[Captures]
accessToken: jsonpath "$.data.access_token"

# --- Create dependencies: category + merchant + product ---------------------
POST {{baseUrl}}/api/category-command/create
Authorization: Bearer {{accessToken}}
[MultipartFormData]
name: Hurl Order Category
description: Category for order E2E
slug_category: hurl-order-category-{{testEmail}}
image: file,assets/test.png; image/png
HTTP 201
[Captures]
categoryId: jsonpath "$.data.id"

POST {{baseUrl}}/api/merchant-command/create
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{
    "name": "Hurl Order Merchant",
    "address": "Jl. Pesanan No. 1",
    "contact_email": "order.{{testEmail}}",
    "contact_phone": "081200000003",
    "description": "Merchant for order E2E",
    "status": "active",
    "user_id": {{userId}}
}
HTTP 201
[Captures]
merchantId: jsonpath "$.data.id"

POST {{baseUrl}}/api/product-command/create
Authorization: Bearer {{accessToken}}
[MultipartFormData]
merchant_id: {{merchantId}}
category_id: {{categoryId}}
name: Hurl Order Product
description: Product for order E2E
price: 25000
count_in_stock: 30
brand: HurlBrand
weight: 100
rating: 5
slug_product: hurl-order-product-{{testEmail}}
image: file,assets/test.png; image/png
HTTP 201
[Captures]
productId: jsonpath "$.data.id"

# --- Create order -----------------------------------------------------------
POST {{baseUrl}}/api/order-command/create
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{
    "merchant_id": {{merchantId}},
    "user_id": {{userId}},
    "total_price": 50000,
    "items": [
        {
            "product_id": {{productId}},
            "quantity": 2,
            "price": 25000
        }
    ],
    "shipping_address": {
        "alamat": "Jl. Antah Berantah No. 9",
        "provinsi": "DKI Jakarta",
        "kota": "Jakarta Selatan",
        "courier": "JNE",
        "shipping_method": "REG",
        "shipping_cost": 10000,
        "negara": "Indonesia"
    }
}
HTTP 200
[Captures]
orderId: jsonpath "$.data.id"
[Asserts]
jsonpath "$.status" == "success"

# --- Queries ----------------------------------------------------------------
GET {{baseUrl}}/api/order-query?page=1&page_size=10
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/order-query/{{orderId}}
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"
jsonpath "$.data.id" == {{orderId}}

GET {{baseUrl}}/api/order-query/active?page=1&page_size=10
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

# --- Order items ------------------------------------------------------------
GET {{baseUrl}}/api/order-item/order/{{orderId}}
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

# --- Trash / restore --------------------------------------------------------
POST {{baseUrl}}/api/order-command/trashed/{{orderId}}
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

POST {{baseUrl}}/api/order-command/restore/{{orderId}}
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"
