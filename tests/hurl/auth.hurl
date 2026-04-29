# Auth Hello
GET {{baseUrl}}/api/auth/hello
HTTP 200

# Register
POST {{baseUrl}}/api/auth/register
{
    "firstname": "John",
    "lastname": "Doe",
    "email": "john.doe@example.com",
    "password": "password123",
    "confirm_password": "password123"
}
HTTP 201

# Login
POST {{baseUrl}}/api/auth/login
{
    "email": "john.doe@example.com",
    "password": "password123"
}
HTTP 200
[Captures]
accessToken: jsonpath "$.data.access_token"
refreshToken: jsonpath "$.data.refresh_token"

# Get Me
GET {{baseUrl}}/api/auth/me
Authorization: Bearer {{accessToken}}
HTTP 200

# Refresh Token
POST {{baseUrl}}/api/auth/refresh-token
{
    "refresh_token": "{{refreshToken}}"
}
HTTP 200
