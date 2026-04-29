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

# Find All Transfers
GET {{baseUrl}}/api/transfers
Authorization: Bearer {{accessToken}}
HTTP 200

# Create Transfer
POST {{baseUrl}}/api/transfers/create
Authorization: Bearer {{accessToken}}
{
    "transfer_from": "{{cardNumber}}",
    "transfer_to": "1234567890123456",
    "transfer_amount": 50000
}
HTTP 200
[Captures]
transferId: jsonpath "$.data.id"

# Find Transfer by ID
GET {{baseUrl}}/api/transfers/{{transferId}}
Authorization: Bearer {{accessToken}}
HTTP 200

# Update Transfer
POST {{baseUrl}}/api/transfers/update/{{transferId}}
Authorization: Bearer {{accessToken}}
{
    "transfer_id": {{transferId}},
    "transfer_from": "{{cardNumber}}",
    "transfer_to": "1234567890123456",
    "transfer_amount": 100000
}
HTTP 200

# Trash Transfer
POST {{baseUrl}}/api/transfers/trashed/{{transferId}}
Authorization: Bearer {{accessToken}}
HTTP 200

# Restore Transfer
POST {{baseUrl}}/api/transfers/restore/{{transferId}}
Authorization: Bearer {{accessToken}}
HTTP 200
