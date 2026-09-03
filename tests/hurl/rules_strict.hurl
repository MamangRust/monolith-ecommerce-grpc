# Strict authorization/rules E2E suite (monolith-ecommerce-grpc).
# - A freshly registered user is auto-assigned ROLE_ADMIN, so role-query is allowed.
# - A seeded user (no role) must be denied with 403.
# Run with:
#   hurl --test --jobs 1 --variable baseUrl=http://127.0.0.1:5000 \
#     --variable adminEmail=... --variable adminPassword=... \
#     --variable userEmail=... --variable userPassword=... \
#     tests/hurl/rules_strict.hurl

# Unauthenticated access must be rejected.
GET {{baseUrl}}/api/role-query?page=1&page_size=10
HTTP 401

# Register + login the admin (register auto-assigns ROLE_ADMIN).
POST {{baseUrl}}/api/auth/register
Content-Type: application/json

{
    "firstname": "Rules",
    "lastname": "Admin",
    "email": "{{adminEmail}}",
    "password": "{{adminPassword}}",
    "confirm_password": "{{adminPassword}}"
}
HTTP 201

POST {{baseUrl}}/api/auth/login
Content-Type: application/json

{
    "email": "{{adminEmail}}",
    "password": "{{adminPassword}}"
}
HTTP 200
[Captures]
adminToken: jsonpath "$.data.access_token"

# Admin can read protected role data.
GET {{baseUrl}}/api/role-query?page=1&page_size=10
Authorization: Bearer {{adminToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/role-query/active?page=1&page_size=10
Authorization: Bearer {{adminToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/role-query/trashed?page=1&page_size=10
Authorization: Bearer {{adminToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/role-query/3
Authorization: Bearer {{adminToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/role-query/user/4
Authorization: Bearer {{adminToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

# A user without any role must be denied by the role validator (401).
POST {{baseUrl}}/api/auth/login
Content-Type: application/json

{
    "email": "{{userEmail}}",
    "password": "{{userPassword}}"
}
HTTP 200
[Captures]
userToken: jsonpath "$.data.access_token"

GET {{baseUrl}}/api/role-query?page=1&page_size=10
Authorization: Bearer {{userToken}}
HTTP 401
