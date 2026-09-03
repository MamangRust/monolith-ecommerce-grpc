# Strict statistics E2E suite.
# Run with:
# hurl --test --jobs 1 --variable baseUrl=http://127.0.0.1:5000 --variable testEmail=... --variable testPassword=... tests/hurl/stats_strict.hurl

POST {{baseUrl}}/api/auth/register
Content-Type: application/json

{
  "firstname": "Stats",
  "lastname": "Tester",
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

# Category statistics
GET {{baseUrl}}/api/category-stats/monthly-total-pricing?year=2024&month=1
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/category-stats/yearly-total-pricing?year=2024
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/category-stats/merchant/monthly-total-pricing?year=2024&month=1&merchant_id=1
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/category-stats/merchant/yearly-total-pricing?year=2024&merchant_id=1
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/category-stats/mycategory/monthly-total-pricing?year=2024&month=1&category_id=1
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/category-stats/mycategory/yearly-total-pricing?year=2024&category_id=1
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/category-stats/monthly-pricing?year=2024&month=1
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/category-stats/yearly-pricing?year=2024
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/category-stats/merchant/monthly-pricing?year=2024&month=1&merchant_id=1
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/category-stats/merchant/yearly-pricing?year=2024&merchant_id=1
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/category-stats/mycategory/monthly-pricing?year=2024&month=1&category_id=1
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/category-stats/mycategory/yearly-pricing?year=2024&category_id=1
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

# Order statistics
GET {{baseUrl}}/api/order/monthly-total-revenue?year=2024&month=1
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/order/yearly-total-revenue?year=2024
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/order/merchant/monthly-total-revenue?year=2024&month=1&merchant_id=1
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/order/merchant/yearly-total-revenue?year=2024&merchant_id=1
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/order/monthly-revenue?year=2024
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/order/yearly-revenue?year=2024
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/order/merchant/monthly-revenue?year=2024&merchant_id=1
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/order/merchant/yearly-revenue?year=2024&merchant_id=1
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

# Transaction statistics
GET {{baseUrl}}/api/transaction-stats/monthly-success?year=2024&month=1
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/transaction-stats/yearly-success?year=2024
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/transaction-stats/monthly-failed?year=2024&month=1
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/transaction-stats/yearly-failed?year=2024
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/transaction-stats/merchant/monthly-success?year=2024&month=1&merchant_id=1
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/transaction-stats/merchant/yearly-success?year=2024&merchant_id=1
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/transaction-stats/merchant/monthly-failed?year=2024&month=1&merchant_id=1
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/transaction-stats/merchant/yearly-failed?year=2024&merchant_id=1
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/transaction-stats/monthly-method-success?year=2024&month=1
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/transaction-stats/yearly-method-success?year=2024
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/transaction-stats/merchant/monthly-method-success?year=2024&month=1&merchant_id=1
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/transaction-stats/merchant/yearly-method-success?year=2024&merchant_id=1
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/transaction-stats/monthly-method-failed?year=2024&month=1
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/transaction-stats/yearly-method-failed?year=2024
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/transaction-stats/merchant/monthly-method-failed?year=2024&month=1&merchant_id=1
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/transaction-stats/merchant/yearly-method-failed?year=2024&merchant_id=1
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"
