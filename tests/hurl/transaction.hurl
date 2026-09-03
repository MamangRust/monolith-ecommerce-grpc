# Transaction domain E2E suite (monolith-ecommerce-grpc).
# Self-contained: creates category + merchant + product + order, then records
# a transaction. NOTE: the transaction amount must cover merchandise total +
# shipping cost + 11% tax, otherwise the API rejects it as "Insufficient balance".
# Run with:
#   hurl --test --jobs 1 --variable baseUrl=http://127.0.0.1:5000 \
#     --variable testEmail=... --variable testPassword=... \
#     --variable testImage=... tests/hurl/transaction.hurl

# --- Auth bootstrap ---------------------------------------------------------
POST {{baseUrl}}/api/auth/register
Content-Type: application/json

{
    "firstname": "Hurl",
    "lastname": "Transaction",
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
name: Hurl Txn Category
description: Category for transaction E2E
slug_category: hurl-txn-category-{{testEmail}}
image: file,assets/test.png; image/png
HTTP 201
[Captures]
categoryId: jsonpath "$.data.id"

POST {{baseUrl}}/api/merchant-command/create
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{
    "name": "Hurl Txn Merchant",
    "address": "Jl. Transaksi No. 1",
    "contact_email": "txn.{{testEmail}}",
    "contact_phone": "081200000004",
    "description": "Merchant for transaction E2E",
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
name: Hurl Txn Product
description: Product for transaction E2E
price: 50000
count_in_stock: 25
brand: HurlBrand
weight: 150
rating: 5
slug_product: hurl-txn-product-{{testEmail}}
image: file,assets/test.png; image/png
HTTP 201
[Captures]
productId: jsonpath "$.data.id"

# --- Create order (merchandise 100000 + shipping 10000 + 11% tax) -----------
POST {{baseUrl}}/api/order-command/create
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{
    "merchant_id": {{merchantId}},
    "user_id": {{userId}},
    "total_price": 100000,
    "items": [
        {
            "product_id": {{productId}},
            "quantity": 2,
            "price": 50000
        }
    ],
    "shipping_address": {
        "alamat": "Jl. Pembayaran No. 1",
        "provinsi": "Jawa Barat",
        "kota": "Bandung",
        "courier": "JNE",
        "shipping_method": "REG",
        "shipping_cost": 10000,
        "negara": "Indonesia"
    }
}
HTTP 200
[Captures]
orderId: jsonpath "$.data.id"

# --- Create transaction (amount must be >= 122100) --------------------------
POST {{baseUrl}}/api/transaction-command/create
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{
    "user_id": {{userId}},
    "order_id": {{orderId}},
    "merchant_id": {{merchantId}},
    "payment_method": "visa",
    "amount": 200000,
    "payment_status": "success"
}
HTTP 201
[Captures]
transactionId: jsonpath "$.data.id"
[Asserts]
jsonpath "$.status" == "success"

# --- Queries ----------------------------------------------------------------
GET {{baseUrl}}/api/transaction-query?page=1&page_size=10
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/transaction-query/{{transactionId}}
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"
jsonpath "$.data.id" == {{transactionId}}

GET {{baseUrl}}/api/transaction-query/merchant/{{merchantId}}?page=1&page_size=10
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

# --- Trash / restore --------------------------------------------------------
POST {{baseUrl}}/api/transaction-command/trashed/{{transactionId}}
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

POST {{baseUrl}}/api/transaction-command/restore/{{transactionId}}
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"
