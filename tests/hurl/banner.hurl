# Banner domain E2E suite (monolith-ecommerce-grpc).
# Run with:
#   hurl --test --jobs 1 --variable baseUrl=http://127.0.0.1:5000 \
#     --variable testEmail=... --variable testPassword=... tests/hurl/banner.hurl

# --- Auth bootstrap ---------------------------------------------------------
POST {{baseUrl}}/api/auth/register
Content-Type: application/json

{
    "firstname": "Hurl",
    "lastname": "Banner",
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

# --- Create banner ----------------------------------------------------------
POST {{baseUrl}}/api/banner-command/create
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{
    "name": "Hurl Banner {{testEmail}}",
    "start_date": "2026-01-01",
    "end_date": "2026-12-31",
    "start_time": "00:00:00",
    "end_time": "23:59:59",
    "is_active": true
}
HTTP 200
[Captures]
bannerId: jsonpath "$.data.id"
[Asserts]
jsonpath "$.status" == "success"

# --- Queries ----------------------------------------------------------------
GET {{baseUrl}}/api/banner-query?page=1&page_size=10
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/banner-query/active?page=1&page_size=10
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/banner-query/{{bannerId}}
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"
jsonpath "$.data.id" == {{bannerId}}

# --- Update -----------------------------------------------------------------
POST {{baseUrl}}/api/banner-command/update/{{bannerId}}
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{
    "banner_id": {{bannerId}},
    "name": "Hurl Banner Updated",
    "start_date": "2026-02-01",
    "end_date": "2026-11-30",
    "start_time": "08:00:00",
    "end_time": "22:00:00",
    "is_active": true
}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

# --- Trash / restore --------------------------------------------------------
POST {{baseUrl}}/api/banner-command/trashed/{{bannerId}}
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/banner-query/trashed?page=1&page_size=10
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

POST {{baseUrl}}/api/banner-command/restore/{{bannerId}}
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"
