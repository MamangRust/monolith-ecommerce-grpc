# Product domain E2E suite (monolith-ecommerce-grpc).
# NOTE: product create/update accept multipart form data with an image file.
# Self-contained: creates a category + merchant, then a product.
# Run with:
#   hurl --test --jobs 1 --variable baseUrl=http://127.0.0.1:5000 \
#     --variable testEmail=... --variable testPassword=... \
#     --variable testImage=... tests/hurl/product.hurl

# --- Auth bootstrap ---------------------------------------------------------
POST {{baseUrl}}/api/auth/register
Content-Type: application/json

{
    "firstname": "Hurl",
    "lastname": "Product",
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

# --- Create category + merchant (dependencies) ------------------------------
POST {{baseUrl}}/api/category-command/create
Authorization: Bearer {{accessToken}}
[MultipartFormData]
name: Hurl Product Category
description: Category for product E2E
slug_category: hurl-product-category-{{testEmail}}
image: file,assets/test.png; image/png
HTTP 201
[Captures]
categoryId: jsonpath "$.data.id"

POST {{baseUrl}}/api/merchant-command/create
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{
    "name": "Hurl Product Merchant",
    "address": "Jl. Produk No. 1",
    "contact_email": "prod.{{testEmail}}",
    "contact_phone": "081200000001",
    "description": "Merchant for product E2E",
    "status": "active",
    "user_id": {{userId}}
}
HTTP 201
[Captures]
merchantId: jsonpath "$.data.id"

# --- Create product (multipart + image) -------------------------------------
POST {{baseUrl}}/api/product-command/create
Authorization: Bearer {{accessToken}}
[MultipartFormData]
merchant_id: {{merchantId}}
category_id: {{categoryId}}
name: Hurl Product {{testEmail}}
description: Product created by hurl E2E test
price: 150000
count_in_stock: 100
brand: HurlBrand
weight: 500
rating: 5
slug_product: hurl-product-{{testEmail}}
image: file,assets/test.png; image/png
HTTP 201
[Captures]
productId: jsonpath "$.data.id"
[Asserts]
jsonpath "$.status" == "success"

# --- Queries ----------------------------------------------------------------
GET {{baseUrl}}/api/product-query?page=1&page_size=10
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/product-query/active?page=1&page_size=10
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/product-query/{{productId}}
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"
jsonpath "$.data.id" == {{productId}}

GET {{baseUrl}}/api/product-query/merchant/{{merchantId}}?page=1&page_size=10
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

# --- Update (multipart + image) ---------------------------------------------
POST {{baseUrl}}/api/product-command/update/{{productId}}
Authorization: Bearer {{accessToken}}
[MultipartFormData]
merchant_id: {{merchantId}}
category_id: {{categoryId}}
name: Hurl Product Updated
description: Product updated by hurl E2E test
price: 175000
count_in_stock: 80
brand: HurlBrand
weight: 550
rating: 4
slug_product: hurl-product-updated-{{testEmail}}
image: file,assets/test.png; image/png
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

# --- Trash / restore --------------------------------------------------------
POST {{baseUrl}}/api/product-command/trashed/{{productId}}
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/product-query/trashed?page=1&page_size=10
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

POST {{baseUrl}}/api/product-command/restore/{{productId}}
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"
