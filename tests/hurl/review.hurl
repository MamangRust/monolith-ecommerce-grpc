# Review domain E2E suite (monolith-ecommerce-grpc).
# Self-contained: creates category + merchant + product, then submits a review.
# Run with:
#   hurl --test --jobs 1 --variable baseUrl=http://127.0.0.1:5000 \
#     --variable testEmail=... --variable testPassword=... tests/hurl/review.hurl

# --- Auth bootstrap ---------------------------------------------------------
POST {{baseUrl}}/api/auth/register
Content-Type: application/json

{
    "firstname": "Hurl",
    "lastname": "Review",
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
name: Hurl Review Category
description: Category for review E2E
slug_category: hurl-review-category-{{testEmail}}
image: file,assets/test.png; image/png
HTTP 201
[Captures]
categoryId: jsonpath "$.data.id"

POST {{baseUrl}}/api/merchant-command/create
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{
    "name": "Hurl Review Merchant",
    "address": "Jl. Ulasan No. 1",
    "contact_email": "review.{{testEmail}}",
    "contact_phone": "081200000005",
    "description": "Merchant for review E2E",
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
name: Hurl Review Product
description: Product for review E2E
price: 75000
count_in_stock: 20
brand: HurlBrand
weight: 300
rating: 5
slug_product: hurl-review-product-{{testEmail}}
image: file,assets/test.png; image/png
HTTP 201
[Captures]
productId: jsonpath "$.data.id"

# --- Create review ----------------------------------------------------------
POST {{baseUrl}}/api/review-command/create
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{
    "product_id": {{productId}},
    "rating": 5,
    "comment": "Produk bagus sekali!",
    "user_id": {{userId}}
}
HTTP 201
[Captures]
reviewId: jsonpath "$.data.id"
[Asserts]
jsonpath "$.status" == "success"

# --- Queries ----------------------------------------------------------------
GET {{baseUrl}}/api/review-query?page=1&page_size=10
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/review-query/product/{{productId}}?page=1&page_size=10
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/review-query/merchant/{{merchantId}}?page=1&page_size=10
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

# --- Update -----------------------------------------------------------------
POST {{baseUrl}}/api/review-command/update/{{reviewId}}
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{
    "product_id": {{productId}},
    "rating": 4,
    "comment": "Ulasan diperbarui",
    "user_id": {{userId}}
}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

# --- Trash / restore --------------------------------------------------------
POST {{baseUrl}}/api/review-command/trashed/{{reviewId}}
Authorization: Bearer {{accessToken}}
HTTP 200

GET {{baseUrl}}/api/review-query/trashed?page=1&page_size=10
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

POST {{baseUrl}}/api/review-command/restore/{{reviewId}}
Authorization: Bearer {{accessToken}}
HTTP 200
