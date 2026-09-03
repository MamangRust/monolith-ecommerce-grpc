# Shipping Address domain E2E suite (monolith-ecommerce-grpc).
# The API exposes no create/update route for shipping addresses (they are
# created implicitly via order checkout), so this suite exercises the
# query + trash/restore lifecycle on an existing (seeded) record.
# Run with:
#   hurl --test --jobs 1 --variable baseUrl=http://127.0.0.1:5000 \
#     --variable testEmail=... --variable testPassword=... tests/hurl/shipping_address.hurl

# --- Auth bootstrap ---------------------------------------------------------
POST {{baseUrl}}/api/auth/register
Content-Type: application/json

{
    "firstname": "Hurl",
    "lastname": "Shipping",
    "email": "{{testEmail}}",
    "password": "{{testPassword}}",
    "confirm_password": "{{testPassword}}"
}
HTTP 201

POST {{baseUrl}}/api/auth/login
Content-Type: application/json

{
    "email": "{{testEmail}}",
    "password": "{{testPassword}}"
}
HTTP 200
[Captures]
accessToken: jsonpath "$.data.access_token"

# --- Queries ----------------------------------------------------------------
GET {{baseUrl}}/api/shipping-address-query?page=1&page_size=10
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"
[Captures]
shippingId: jsonpath "$.data[0].id"

GET {{baseUrl}}/api/shipping-address-query/{{shippingId}}
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"
jsonpath "$.data.id" == {{shippingId}}

GET {{baseUrl}}/api/shipping-address-query/active?page=1&page_size=10
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

# --- Trash / restore --------------------------------------------------------
POST {{baseUrl}}/api/shipping-address-command/trashed/{{shippingId}}
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/shipping-address-query/trashed?page=1&page_size=10
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

POST {{baseUrl}}/api/shipping-address-command/restore/{{shippingId}}
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"
