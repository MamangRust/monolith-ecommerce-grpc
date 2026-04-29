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

# Find All Topups
GET {{baseUrl}}/api/topups
Authorization: Bearer {{accessToken}}
HTTP 200

# Create Topup
POST {{baseUrl}}/api/topups/create
Authorization: Bearer {{accessToken}}
{
    "card_number": "{{cardNumber}}",
    "topup_amount": 50000,
    "topup_method": "visa"
}
HTTP 200
[Captures]
topupId: jsonpath "$.data.topup_id"

# Find Topup by ID
GET {{baseUrl}}/api/topups/{{topupId}}
Authorization: Bearer {{accessToken}}
HTTP 200

# Find Topups by Card Number
GET {{baseUrl}}/api/topups/card-number/{{cardNumber}}
Authorization: Bearer {{accessToken}}
HTTP 200

# Update Topup
POST {{baseUrl}}/api/topups/update/{{topupId}}
Authorization: Bearer {{accessToken}}
{
    "topup_id": {{topupId}},
    "card_number": "{{cardNumber}}",
    "topup_amount": 100000,
    "topup_method": "visa"
}
HTTP 200

# Trash Topup
POST {{baseUrl}}/api/topups/trashed/{{topupId}}
Authorization: Bearer {{accessToken}}
HTTP 200

# Restore Topup
POST {{baseUrl}}/api/topups/restore/{{topupId}}
Authorization: Bearer {{accessToken}}
HTTP 200
