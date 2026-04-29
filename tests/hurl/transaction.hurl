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

# Find All Transactions
GET {{baseUrl}}/api/transactions
Authorization: Bearer {{accessToken}}
HTTP 200

# Create Transaction
POST {{baseUrl}}/api/transactions/create
Authorization: Bearer {{accessToken}}
{
    "card_number": "{{cardNumber}}",
    "amount": 50000,
    "payment_method": "visa",
    "merchant_id": {{merchantId}},
    "transaction_time": "2026-04-09T19:41:10Z"
}
HTTP 200
[Captures]
transactionId: jsonpath "$.data.id"

# Find Transaction by ID
GET {{baseUrl}}/api/transactions/{{transactionId}}
Authorization: Bearer {{accessToken}}
HTTP 200

# Find Transactions by Card Number
GET {{baseUrl}}/api/transactions/card-number/{{cardNumber}}
Authorization: Bearer {{accessToken}}
HTTP 200

# Update Transaction
POST {{baseUrl}}/api/transactions/update/{{transactionId}}
Authorization: Bearer {{accessToken}}
{
    "transaction_id": {{transactionId}},
    "card_number": "{{cardNumber}}",
    "amount": 100000,
    "payment_method": "visa",
    "merchant_id": {{merchantId}},
    "transaction_time": "2026-04-09T20:41:10Z"
}
HTTP 200

# Trash Transaction
POST {{baseUrl}}/api/transactions/trashed/{{transactionId}}
Authorization: Bearer {{accessToken}}
HTTP 200

# Restore Transaction
POST {{baseUrl}}/api/transactions/restore/{{transactionId}}
Authorization: Bearer {{accessToken}}
HTTP 200
