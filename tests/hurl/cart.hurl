# Cart domain E2E suite (monolith-ecommerce-grpc).
# Self-contained: creates category + merchant + product, then adds to cart.
# Run with:
#   hurl --test --jobs 1 --variable baseUrl=http://127.0.0.1:5000 \
#     --variable testEmail=... --variable testPassword=... \
#     --variable testImage=... tests/hurl/cart.hurl

# --- Auth bootstrap ---------------------------------------------------------
POST {{baseUrl}}/api/auth/register
Content-Type: application/json

{
    "firstname": "Hurl",
    "lastname": "Cart",
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
name: Hurl Cart Category
description: Category for cart E2E
slug_category: hurl-cart-category-{{testEmail}}
image: file,assets/test.png; image/png
HTTP 201
[Captures]
categoryId: jsonpath "$.data.id"

POST {{baseUrl}}/api/merchant-command/create
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{
    "name": "Hurl Cart Merchant",
    "address": "Jl. Keranjang No. 1",
    "contact_email": "cart.{{testEmail}}",
    "contact_phone": "081200000002",
    "description": "Merchant for cart E2E",
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
name: Hurl Cart Product
description: Product for cart E2E
price: 50000
count_in_stock: 50
brand: HurlBrand
weight: 200
rating: 5
slug_product: hurl-cart-product-{{testEmail}}
image: file,assets/test.png; image/png
HTTP 201
[Captures]
productId: jsonpath "$.data.id"

# --- Add to cart ------------------------------------------------------------
POST {{baseUrl}}/api/cart-command/create
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{
    "product_id": {{productId}},
    "quantity": 2,
    "user_id": {{userId}}
}
HTTP 201
[Captures]
cartItemId: jsonpath "$.data.id"
[Asserts]
jsonpath "$.status" == "success"

# --- Query cart -------------------------------------------------------------
GET {{baseUrl}}/api/cart-query
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

# --- Delete item from cart (needs cart_id + user_id) ------------------------
DELETE {{baseUrl}}/api/cart-command/delete
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{
    "cart_id": {{cartItemId}},
    "user_id": {{userId}}
}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

# --- Delete all (needs cart_ids + user_id) ----------------------------------
POST {{baseUrl}}/api/cart-command/create
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{
    "product_id": {{productId}},
    "quantity": 1,
    "user_id": {{userId}}
}
HTTP 201
[Captures]
cartItemId2: jsonpath "$.data.id"

POST {{baseUrl}}/api/cart-command/delete-all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{
    "user_id": {{userId}},
    "cart_ids": [{{cartItemId2}}]
}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"
