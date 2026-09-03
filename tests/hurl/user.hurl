# User domain E2E suite (monolith-ecommerce-grpc).
# Run with:
#   hurl --test --jobs 1 --variable baseUrl=http://127.0.0.1:5000 \
#     --variable testEmail=... --variable testPassword=... tests/hurl/user.hurl

# --- Auth bootstrap ---------------------------------------------------------
POST {{baseUrl}}/api/auth/register
Content-Type: application/json

{
    "firstname": "Hurl",
    "lastname": "User",
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

# --- Create user ------------------------------------------------------------
POST {{baseUrl}}/api/user-command/create
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{
    "firstname": "Alice",
    "lastname": "Smith",
    "email": "alice.{{testEmail}}",
    "password": "{{testPassword}}",
    "confirm_password": "{{testPassword}}"
}
HTTP 200
[Captures]
newUserId: jsonpath "$.data.id"
[Asserts]
jsonpath "$.status" == "success"

# --- Queries ----------------------------------------------------------------
GET {{baseUrl}}/api/user-query?page=1&page_size=10
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/user-query/active?page=1&page_size=10
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/user-query/{{newUserId}}
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"
jsonpath "$.data.id" == {{newUserId}}

# --- Update -----------------------------------------------------------------
POST {{baseUrl}}/api/user-command/update/{{newUserId}}
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{
    "firstname": "Alicia",
    "lastname": "Smith",
    "email": "alice.{{testEmail}}",
    "password": "{{testPassword}}",
    "confirm_password": "{{testPassword}}"
}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

# --- Trash / restore --------------------------------------------------------
POST {{baseUrl}}/api/user-command/trashed/{{newUserId}}
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/user-query/trashed?page=1&page_size=10
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

POST {{baseUrl}}/api/user-command/restore/{{newUserId}}
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

# --- Permanent delete -------------------------------------------------------
DELETE {{baseUrl}}/api/user-command/permanent/{{newUserId}}
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"
