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

# Find All Users
GET {{baseUrl}}/api/user
Authorization: Bearer {{accessToken}}
HTTP 200

# Create User
POST {{baseUrl}}/api/user/create
Authorization: Bearer {{accessToken}}
{
    "firstname": "Alice",
    "lastname": "Smith",
    "email": "alice.smith@example.com",
    "password": "password123",
    "confirm_password": "password123"
}
HTTP 200
[Captures]
userId: jsonpath "$.data.id"

# Find User by ID
GET {{baseUrl}}/api/user/{{userId}}
Authorization: Bearer {{accessToken}}
HTTP 200

# Update User
POST {{baseUrl}}/api/user/update/{{userId}}
Authorization: Bearer {{accessToken}}
{
    "firstname": "Alicia",
    "lastname": "Smith",
    "email": "alice.smith@example.com",
    "password": "password123456",
    "confirm_password": "password123456"
}
HTTP 200

# Find Active Users
GET {{baseUrl}}/api/user/active
Authorization: Bearer {{accessToken}}
HTTP 200

# Trash User
POST {{baseUrl}}/api/user/trashed/{{userId}}
Authorization: Bearer {{accessToken}}
HTTP 200

# Find Trashed Users
GET {{baseUrl}}/api/user/trashed
Authorization: Bearer {{accessToken}}
HTTP 200

# Restore User
POST {{baseUrl}}/api/user/restore/{{userId}}
Authorization: Bearer {{accessToken}}
HTTP 200

# Permanent Delete User
DELETE {{baseUrl}}/api/user/permanent/{{userId}}
Authorization: Bearer {{accessToken}}
HTTP 200
