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

# Find All Merchants
GET {{baseUrl}}/api/merchants
Authorization: Bearer {{accessToken}}
HTTP 200

# Create Merchant
POST {{baseUrl}}/api/merchants/create
Authorization: Bearer {{accessToken}}
{
    "name": "Global Store",
    "user_id": {{userId}}
}
HTTP 200
[Captures]
merchantId: jsonpath "$.data.id"
apiKey: jsonpath "$.data.api_key"

# Find Merchant by ID
GET {{baseUrl}}/api/merchants/{{merchantId}}
Authorization: Bearer {{accessToken}}
HTTP 200

# Find Merchant by API Key
GET {{baseUrl}}/api/merchants/api-key?api_key={{apiKey}}
Authorization: Bearer {{accessToken}}
HTTP 200

# Update Merchant
POST {{baseUrl}}/api/merchants/updates/{{merchantId}}
Authorization: Bearer {{accessToken}}
{
    "name": "Global Store Updated",
    "user_id": {{userId}},
    "status": "active"
}
HTTP 200

# Find Active Merchants
GET {{baseUrl}}/api/merchants/active
Authorization: Bearer {{accessToken}}
HTTP 200

# Trash Merchant
POST {{baseUrl}}/api/merchants/trashed/{{merchantId}}
Authorization: Bearer {{accessToken}}
HTTP 200

# Restore Merchant
POST {{baseUrl}}/api/merchants/restore/{{merchantId}}
Authorization: Bearer {{accessToken}}
HTTP 200

# Find Monthly Amount
GET {{baseUrl}}/api/merchants/monthly-amount?year=2024
Authorization: Bearer {{accessToken}}
HTTP 200
