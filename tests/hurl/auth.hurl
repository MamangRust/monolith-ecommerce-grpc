# Auth domain E2E suite (monolith-ecommerce-grpc).
# Run with:
#   hurl --test --jobs 1 --variable baseUrl=http://127.0.0.1:5000 \
#     --variable testEmail=... --variable testPassword=... tests/hurl/auth.hurl

# Hello (public liveness)
GET {{baseUrl}}/api/auth/hello
HTTP 200

# Register
POST {{baseUrl}}/api/auth/register
Content-Type: application/json

{
    "firstname": "John",
    "lastname": "Doe",
    "email": "{{testEmail}}",
    "password": "{{testPassword}}",
    "confirm_password": "{{testPassword}}"
}
HTTP 201

# Login
POST {{baseUrl}}/api/auth/login
Content-Type: application/json

{
    "email": "{{testEmail}}",
    "password": "{{testPassword}}"
}
HTTP 200
[Captures]
accessToken: jsonpath "$.data.access_token"
refreshToken: jsonpath "$.data.refresh_token"

# Get Me (JWT protected)
GET {{baseUrl}}/api/auth/me
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

# Refresh Token
POST {{baseUrl}}/api/auth/refresh-token
Content-Type: application/json

{
    "refresh_token": "{{refreshToken}}"
}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"
