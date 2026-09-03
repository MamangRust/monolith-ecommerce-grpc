# Generated API Gateway route sweep from service/apigateway/docs/swagger.json
# Run with:
# hurl --test --jobs 1 --variable baseUrl=http://127.0.0.1:5000 --variable testEmail=... --variable testPassword=... tests/hurl/all-endpoints.hurl

# Establish a valid JWT for protected endpoints.
POST {{baseUrl}}/api/auth/register
Content-Type: application/json

{
    "firstname": "Hurl",
    "lastname": "Test",
    "email": "{{testEmail}}",
    "password": "{{testPassword}}",
    "confirm_password": "{{testPassword}}"
}
HTTP *
[Asserts]
status < 500

POST {{baseUrl}}/api/auth/login
Content-Type: application/json

{
    "email": "{{testEmail}}",
    "password": "{{testPassword}}"
}
HTTP *
[Captures]
accessToken: jsonpath "$.data.access_token"
[Asserts]
status < 500

# 001 GET /api/auth/hello
GET {{baseUrl}}/api/auth/hello
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 002 POST /api/auth/login
POST {{baseUrl}}/api/auth/login
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 003 GET /api/auth/me
GET {{baseUrl}}/api/auth/me
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 004 POST /api/auth/refresh-token
POST {{baseUrl}}/api/auth/refresh-token
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 005 POST /api/auth/register
POST {{baseUrl}}/api/auth/register
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 006 POST /api/banner-command/create
POST {{baseUrl}}/api/banner-command/create
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 007 POST /api/banner-command/permanent/all
POST {{baseUrl}}/api/banner-command/permanent/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 008 DELETE /api/banner-command/permanent/{id}
DELETE {{baseUrl}}/api/banner-command/permanent/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 009 POST /api/banner-command/restore/all
POST {{baseUrl}}/api/banner-command/restore/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 010 POST /api/banner-command/restore/{id}
POST {{baseUrl}}/api/banner-command/restore/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 011 POST /api/banner-command/trashed/{id}
POST {{baseUrl}}/api/banner-command/trashed/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 012 POST /api/banner-command/update/{id}
POST {{baseUrl}}/api/banner-command/update/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 013 GET /api/banner-query
GET {{baseUrl}}/api/banner-query
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 014 GET /api/banner-query/active
GET {{baseUrl}}/api/banner-query/active
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 015 GET /api/banner-query/trashed
GET {{baseUrl}}/api/banner-query/trashed
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 016 GET /api/banner-query/{id}
GET {{baseUrl}}/api/banner-query/1
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 017 POST /api/cart-command/create
POST {{baseUrl}}/api/cart-command/create
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 018 DELETE /api/cart-command/delete
DELETE {{baseUrl}}/api/cart-command/delete
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 019 POST /api/cart-command/delete-all
POST {{baseUrl}}/api/cart-command/delete-all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 020 GET /api/cart-query
GET {{baseUrl}}/api/cart-query
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 021 POST /api/category-command/create
POST {{baseUrl}}/api/category-command/create
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 022 DELETE /api/category-command/permanent/all
DELETE {{baseUrl}}/api/category-command/permanent/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 023 DELETE /api/category-command/permanent/{id}
DELETE {{baseUrl}}/api/category-command/permanent/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 024 POST /api/category-command/restore/all
POST {{baseUrl}}/api/category-command/restore/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 025 POST /api/category-command/restore/{id}
POST {{baseUrl}}/api/category-command/restore/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 026 POST /api/category-command/trashed/{id}
POST {{baseUrl}}/api/category-command/trashed/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 027 POST /api/category-command/update/{id}
POST {{baseUrl}}/api/category-command/update/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 028 GET /api/category-query
GET {{baseUrl}}/api/category-query
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 029 GET /api/category-query/active
GET {{baseUrl}}/api/category-query/active
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 030 GET /api/category-query/trashed
GET {{baseUrl}}/api/category-query/trashed
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 031 GET /api/category-query/{id}
GET {{baseUrl}}/api/category-query/1
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 032 GET /api/category-stats/merchant/monthly-pricing
GET {{baseUrl}}/api/category-stats/merchant/monthly-pricing
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 033 GET /api/category-stats/merchant/monthly-total-pricing
GET {{baseUrl}}/api/category-stats/merchant/monthly-total-pricing
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 034 GET /api/category-stats/merchant/yearly-pricing
GET {{baseUrl}}/api/category-stats/merchant/yearly-pricing
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 035 GET /api/category-stats/merchant/yearly-total-pricing
GET {{baseUrl}}/api/category-stats/merchant/yearly-total-pricing
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 036 GET /api/category-stats/monthly-pricing
GET {{baseUrl}}/api/category-stats/monthly-pricing
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 037 GET /api/category-stats/monthly-total-pricing
GET {{baseUrl}}/api/category-stats/monthly-total-pricing
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 038 GET /api/category-stats/mycategory/monthly-pricing
GET {{baseUrl}}/api/category-stats/mycategory/monthly-pricing
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 039 GET /api/category-stats/mycategory/monthly-total-pricing
GET {{baseUrl}}/api/category-stats/mycategory/monthly-total-pricing
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 040 GET /api/category-stats/mycategory/yearly-pricing
GET {{baseUrl}}/api/category-stats/mycategory/yearly-pricing
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 041 GET /api/category-stats/mycategory/yearly-total-pricing
GET {{baseUrl}}/api/category-stats/mycategory/yearly-total-pricing
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 042 GET /api/category-stats/yearly-pricing
GET {{baseUrl}}/api/category-stats/yearly-pricing
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 043 GET /api/category-stats/yearly-total-pricing
GET {{baseUrl}}/api/category-stats/yearly-total-pricing
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 044 POST /api/merchant-award-command/create
POST {{baseUrl}}/api/merchant-award-command/create
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 045 POST /api/merchant-award-command/permanent/all
POST {{baseUrl}}/api/merchant-award-command/permanent/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 046 DELETE /api/merchant-award-command/permanent/{id}
DELETE {{baseUrl}}/api/merchant-award-command/permanent/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 047 POST /api/merchant-award-command/restore/all
POST {{baseUrl}}/api/merchant-award-command/restore/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 048 POST /api/merchant-award-command/restore/{id}
POST {{baseUrl}}/api/merchant-award-command/restore/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 049 POST /api/merchant-award-command/trashed/{id}
POST {{baseUrl}}/api/merchant-award-command/trashed/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 050 POST /api/merchant-award-command/update/{id}
POST {{baseUrl}}/api/merchant-award-command/update/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 051 GET /api/merchant-award-query
GET {{baseUrl}}/api/merchant-award-query
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 052 GET /api/merchant-award-query/active
GET {{baseUrl}}/api/merchant-award-query/active
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 053 GET /api/merchant-award-query/trashed
GET {{baseUrl}}/api/merchant-award-query/trashed
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 054 GET /api/merchant-award-query/{id}
GET {{baseUrl}}/api/merchant-award-query/1
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 055 POST /api/merchant-business-command/create
POST {{baseUrl}}/api/merchant-business-command/create
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 056 POST /api/merchant-business-command/permanent/all
POST {{baseUrl}}/api/merchant-business-command/permanent/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 057 DELETE /api/merchant-business-command/permanent/{id}
DELETE {{baseUrl}}/api/merchant-business-command/permanent/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 058 POST /api/merchant-business-command/restore/all
POST {{baseUrl}}/api/merchant-business-command/restore/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 059 POST /api/merchant-business-command/restore/{id}
POST {{baseUrl}}/api/merchant-business-command/restore/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 060 POST /api/merchant-business-command/trashed/{id}
POST {{baseUrl}}/api/merchant-business-command/trashed/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 061 POST /api/merchant-business-command/update/{id}
POST {{baseUrl}}/api/merchant-business-command/update/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 062 GET /api/merchant-business-query
GET {{baseUrl}}/api/merchant-business-query
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 063 GET /api/merchant-business-query/active
GET {{baseUrl}}/api/merchant-business-query/active
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 064 GET /api/merchant-business-query/trashed
GET {{baseUrl}}/api/merchant-business-query/trashed
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 065 GET /api/merchant-business-query/{id}
GET {{baseUrl}}/api/merchant-business-query/1
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 066 POST /api/merchant-command/create
POST {{baseUrl}}/api/merchant-command/create
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 067 POST /api/merchant-command/permanent/all
POST {{baseUrl}}/api/merchant-command/permanent/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 068 DELETE /api/merchant-command/permanent/{id}
DELETE {{baseUrl}}/api/merchant-command/permanent/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 069 POST /api/merchant-command/restore/all
POST {{baseUrl}}/api/merchant-command/restore/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 070 POST /api/merchant-command/restore/{id}
POST {{baseUrl}}/api/merchant-command/restore/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 071 POST /api/merchant-command/trashed/{id}
POST {{baseUrl}}/api/merchant-command/trashed/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 072 POST /api/merchant-command/update/{id}
POST {{baseUrl}}/api/merchant-command/update/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 073 POST /api/merchant-detail-command/create
POST {{baseUrl}}/api/merchant-detail-command/create
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 074 POST /api/merchant-detail-command/permanent/all
POST {{baseUrl}}/api/merchant-detail-command/permanent/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 075 DELETE /api/merchant-detail-command/permanent/{id}
DELETE {{baseUrl}}/api/merchant-detail-command/permanent/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 076 POST /api/merchant-detail-command/restore/all
POST {{baseUrl}}/api/merchant-detail-command/restore/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 077 POST /api/merchant-detail-command/restore/{id}
POST {{baseUrl}}/api/merchant-detail-command/restore/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 078 POST /api/merchant-detail-command/trashed/{id}
POST {{baseUrl}}/api/merchant-detail-command/trashed/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 079 POST /api/merchant-detail-command/update/{id}
POST {{baseUrl}}/api/merchant-detail-command/update/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 080 GET /api/merchant-detail-query
GET {{baseUrl}}/api/merchant-detail-query
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 081 GET /api/merchant-detail-query/active
GET {{baseUrl}}/api/merchant-detail-query/active
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 082 GET /api/merchant-detail-query/trashed
GET {{baseUrl}}/api/merchant-detail-query/trashed
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 083 GET /api/merchant-detail-query/{id}
GET {{baseUrl}}/api/merchant-detail-query/1
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 084 POST /api/merchant-document-command/create
POST {{baseUrl}}/api/merchant-document-command/create
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 085 POST /api/merchant-document-command/permanent/all
POST {{baseUrl}}/api/merchant-document-command/permanent/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 086 DELETE /api/merchant-document-command/permanent/{id}
DELETE {{baseUrl}}/api/merchant-document-command/permanent/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 087 POST /api/merchant-document-command/restore/all
POST {{baseUrl}}/api/merchant-document-command/restore/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 088 POST /api/merchant-document-command/restore/{id}
POST {{baseUrl}}/api/merchant-document-command/restore/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 089 POST /api/merchant-document-command/trashed/{id}
POST {{baseUrl}}/api/merchant-document-command/trashed/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 090 POST /api/merchant-document-command/update-status/{id}
POST {{baseUrl}}/api/merchant-document-command/update-status/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 091 POST /api/merchant-document-command/update/{id}
POST {{baseUrl}}/api/merchant-document-command/update/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 092 GET /api/merchant-document-query
GET {{baseUrl}}/api/merchant-document-query
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 093 GET /api/merchant-document-query/active
GET {{baseUrl}}/api/merchant-document-query/active
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 094 GET /api/merchant-document-query/trashed
GET {{baseUrl}}/api/merchant-document-query/trashed
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 095 GET /api/merchant-document-query/{id}
GET {{baseUrl}}/api/merchant-document-query/1
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 096 POST /api/merchant-policy-command/create
POST {{baseUrl}}/api/merchant-policy-command/create
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 097 POST /api/merchant-policy-command/permanent/all
POST {{baseUrl}}/api/merchant-policy-command/permanent/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 098 DELETE /api/merchant-policy-command/permanent/{id}
DELETE {{baseUrl}}/api/merchant-policy-command/permanent/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 099 POST /api/merchant-policy-command/restore/all
POST {{baseUrl}}/api/merchant-policy-command/restore/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 100 POST /api/merchant-policy-command/restore/{id}
POST {{baseUrl}}/api/merchant-policy-command/restore/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 101 POST /api/merchant-policy-command/trashed/{id}
POST {{baseUrl}}/api/merchant-policy-command/trashed/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 102 POST /api/merchant-policy-command/update/{id}
POST {{baseUrl}}/api/merchant-policy-command/update/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 103 GET /api/merchant-policy-query
GET {{baseUrl}}/api/merchant-policy-query
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 104 GET /api/merchant-policy-query/active
GET {{baseUrl}}/api/merchant-policy-query/active
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 105 GET /api/merchant-policy-query/trashed
GET {{baseUrl}}/api/merchant-policy-query/trashed
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 106 GET /api/merchant-policy-query/{id}
GET {{baseUrl}}/api/merchant-policy-query/1
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 107 GET /api/merchant-query
GET {{baseUrl}}/api/merchant-query
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 108 GET /api/merchant-query/active
GET {{baseUrl}}/api/merchant-query/active
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 109 GET /api/merchant-query/trashed
GET {{baseUrl}}/api/merchant-query/trashed
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 110 GET /api/merchant-query/{id}
GET {{baseUrl}}/api/merchant-query/1
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 111 POST /api/merchant-social-link/create
POST {{baseUrl}}/api/merchant-social-link/create
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 112 POST /api/merchant-social-link/update/{id}
POST {{baseUrl}}/api/merchant-social-link/update/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 113 POST /api/order-command/create
POST {{baseUrl}}/api/order-command/create
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 114 POST /api/order-command/permanent/all
POST {{baseUrl}}/api/order-command/permanent/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 115 DELETE /api/order-command/permanent/{id}
DELETE {{baseUrl}}/api/order-command/permanent/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 116 POST /api/order-command/restore/all
POST {{baseUrl}}/api/order-command/restore/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 117 POST /api/order-command/restore/{id}
POST {{baseUrl}}/api/order-command/restore/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 118 POST /api/order-command/trashed/{id}
POST {{baseUrl}}/api/order-command/trashed/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 119 POST /api/order-command/update/{id}
POST {{baseUrl}}/api/order-command/update/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 120 GET /api/order-item
GET {{baseUrl}}/api/order-item
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 121 GET /api/order-item/active
GET {{baseUrl}}/api/order-item/active
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 122 GET /api/order-item/order/{order_id}
GET {{baseUrl}}/api/order-item/order/1
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 123 POST /api/order-item/permanent/all
POST {{baseUrl}}/api/order-item/permanent/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 124 DELETE /api/order-item/permanent/{id}
DELETE {{baseUrl}}/api/order-item/permanent/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 125 POST /api/order-item/restore/all
POST {{baseUrl}}/api/order-item/restore/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 126 POST /api/order-item/restore/{id}
POST {{baseUrl}}/api/order-item/restore/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 127 POST /api/order-item/trash/{id}
POST {{baseUrl}}/api/order-item/trash/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 128 GET /api/order-item/trashed
GET {{baseUrl}}/api/order-item/trashed
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 129 GET /api/order-query
GET {{baseUrl}}/api/order-query
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 130 GET /api/order-query/active
GET {{baseUrl}}/api/order-query/active
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 131 GET /api/order-query/trashed
GET {{baseUrl}}/api/order-query/trashed
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 132 GET /api/order-query/{id}
GET {{baseUrl}}/api/order-query/1
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 133 GET /api/order/merchant/monthly-revenue
GET {{baseUrl}}/api/order/merchant/monthly-revenue
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 134 GET /api/order/merchant/monthly-total-revenue
GET {{baseUrl}}/api/order/merchant/monthly-total-revenue
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 135 GET /api/order/merchant/yearly-revenue
GET {{baseUrl}}/api/order/merchant/yearly-revenue
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 136 GET /api/order/merchant/yearly-total-revenue
GET {{baseUrl}}/api/order/merchant/yearly-total-revenue
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 137 GET /api/order/monthly-revenue
GET {{baseUrl}}/api/order/monthly-revenue
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 138 GET /api/order/monthly-total-revenue
GET {{baseUrl}}/api/order/monthly-total-revenue
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 139 GET /api/order/yearly-revenue
GET {{baseUrl}}/api/order/yearly-revenue
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 140 GET /api/order/yearly-total-revenue
GET {{baseUrl}}/api/order/yearly-total-revenue
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 141 POST /api/product-command/create
POST {{baseUrl}}/api/product-command/create
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 142 POST /api/product-command/permanent/all
POST {{baseUrl}}/api/product-command/permanent/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 143 DELETE /api/product-command/permanent/{id}
DELETE {{baseUrl}}/api/product-command/permanent/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 144 POST /api/product-command/restore/all
POST {{baseUrl}}/api/product-command/restore/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 145 POST /api/product-command/restore/{id}
POST {{baseUrl}}/api/product-command/restore/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 146 POST /api/product-command/trashed/{id}
POST {{baseUrl}}/api/product-command/trashed/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 147 POST /api/product-command/update/{id}
POST {{baseUrl}}/api/product-command/update/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 148 GET /api/product-query
GET {{baseUrl}}/api/product-query
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 149 GET /api/product-query/active
GET {{baseUrl}}/api/product-query/active
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 150 GET /api/product-query/category/{category_name}
GET {{baseUrl}}/api/product-query/category/1
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 151 GET /api/product-query/merchant/{merchant_id}
GET {{baseUrl}}/api/product-query/merchant/1
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 152 GET /api/product-query/trashed
GET {{baseUrl}}/api/product-query/trashed
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 153 GET /api/product-query/{id}
GET {{baseUrl}}/api/product-query/1
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 154 POST /api/review-command/create
POST {{baseUrl}}/api/review-command/create
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 155 POST /api/review-command/permanent/all
POST {{baseUrl}}/api/review-command/permanent/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 156 DELETE /api/review-command/permanent/{id}
DELETE {{baseUrl}}/api/review-command/permanent/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 157 POST /api/review-command/restore/all
POST {{baseUrl}}/api/review-command/restore/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 158 POST /api/review-command/restore/{id}
POST {{baseUrl}}/api/review-command/restore/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 159 POST /api/review-command/trashed/{id}
POST {{baseUrl}}/api/review-command/trashed/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 160 POST /api/review-command/update/{id}
POST {{baseUrl}}/api/review-command/update/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 161 POST /api/review-detail-command/create
POST {{baseUrl}}/api/review-detail-command/create
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 162 POST /api/review-detail-command/permanent/all
POST {{baseUrl}}/api/review-detail-command/permanent/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 163 DELETE /api/review-detail-command/permanent/{id}
DELETE {{baseUrl}}/api/review-detail-command/permanent/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 164 POST /api/review-detail-command/restore/all
POST {{baseUrl}}/api/review-detail-command/restore/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 165 POST /api/review-detail-command/restore/{id}
POST {{baseUrl}}/api/review-detail-command/restore/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 166 POST /api/review-detail-command/trashed/{id}
POST {{baseUrl}}/api/review-detail-command/trashed/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 167 POST /api/review-detail-command/update/{id}
POST {{baseUrl}}/api/review-detail-command/update/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 168 GET /api/review-detail-query
GET {{baseUrl}}/api/review-detail-query
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 169 GET /api/review-detail-query/active
GET {{baseUrl}}/api/review-detail-query/active
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 170 GET /api/review-detail-query/trashed
GET {{baseUrl}}/api/review-detail-query/trashed
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 171 GET /api/review-detail-query/{id}
GET {{baseUrl}}/api/review-detail-query/1
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 172 GET /api/review-query
GET {{baseUrl}}/api/review-query
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 173 GET /api/review-query/active
GET {{baseUrl}}/api/review-query/active
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 174 GET /api/review-query/merchant/{id}
GET {{baseUrl}}/api/review-query/merchant/1
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 175 GET /api/review-query/product/{id}
GET {{baseUrl}}/api/review-query/product/1
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 176 GET /api/review-query/trashed
GET {{baseUrl}}/api/review-query/trashed
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 177 POST /api/role-command/create
POST {{baseUrl}}/api/role-command/create
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 178 DELETE /api/role-command/permanent/all
DELETE {{baseUrl}}/api/role-command/permanent/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 179 DELETE /api/role-command/permanent/{id}
DELETE {{baseUrl}}/api/role-command/permanent/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 180 POST /api/role-command/restore/all
POST {{baseUrl}}/api/role-command/restore/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 181 PUT /api/role-command/restore/{id}
PUT {{baseUrl}}/api/role-command/restore/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 182 POST /api/role-command/trashed/{id}
POST {{baseUrl}}/api/role-command/trashed/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 183 POST /api/role-command/update/{id}
POST {{baseUrl}}/api/role-command/update/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 184 GET /api/role-query
GET {{baseUrl}}/api/role-query
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 185 GET /api/role-query/active
GET {{baseUrl}}/api/role-query/active
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 186 GET /api/role-query/trashed
GET {{baseUrl}}/api/role-query/trashed
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 187 GET /api/role-query/user/{user_id}
GET {{baseUrl}}/api/role-query/user/1
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 188 GET /api/role-query/{id}
GET {{baseUrl}}/api/role-query/1
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 189 POST /api/shipping-address-command/permanent/all
POST {{baseUrl}}/api/shipping-address-command/permanent/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 190 DELETE /api/shipping-address-command/permanent/{id}
DELETE {{baseUrl}}/api/shipping-address-command/permanent/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 191 POST /api/shipping-address-command/restore/all
POST {{baseUrl}}/api/shipping-address-command/restore/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 192 POST /api/shipping-address-command/restore/{id}
POST {{baseUrl}}/api/shipping-address-command/restore/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 193 POST /api/shipping-address-command/trashed/{id}
POST {{baseUrl}}/api/shipping-address-command/trashed/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 194 GET /api/shipping-address-query
GET {{baseUrl}}/api/shipping-address-query
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 195 GET /api/shipping-address-query/active
GET {{baseUrl}}/api/shipping-address-query/active
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 196 GET /api/shipping-address-query/order/{id}
GET {{baseUrl}}/api/shipping-address-query/order/1
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 197 GET /api/shipping-address-query/trashed
GET {{baseUrl}}/api/shipping-address-query/trashed
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 198 GET /api/shipping-address-query/{id}
GET {{baseUrl}}/api/shipping-address-query/1
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 199 POST /api/slider-command/create
POST {{baseUrl}}/api/slider-command/create
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 200 POST /api/slider-command/permanent/all
POST {{baseUrl}}/api/slider-command/permanent/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 201 DELETE /api/slider-command/permanent/{id}
DELETE {{baseUrl}}/api/slider-command/permanent/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 202 POST /api/slider-command/restore/all
POST {{baseUrl}}/api/slider-command/restore/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 203 POST /api/slider-command/restore/{id}
POST {{baseUrl}}/api/slider-command/restore/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 204 POST /api/slider-command/trashed/{id}
POST {{baseUrl}}/api/slider-command/trashed/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 205 POST /api/slider-command/update/{id}
POST {{baseUrl}}/api/slider-command/update/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 206 GET /api/slider-query
GET {{baseUrl}}/api/slider-query
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 207 GET /api/slider-query/active
GET {{baseUrl}}/api/slider-query/active
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 208 GET /api/slider-query/trashed
GET {{baseUrl}}/api/slider-query/trashed
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 209 GET /api/slider-query/{id}
GET {{baseUrl}}/api/slider-query/1
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 210 POST /api/transaction-command/create
POST {{baseUrl}}/api/transaction-command/create
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 211 POST /api/transaction-command/permanent/all
POST {{baseUrl}}/api/transaction-command/permanent/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 212 DELETE /api/transaction-command/permanent/{id}
DELETE {{baseUrl}}/api/transaction-command/permanent/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 213 POST /api/transaction-command/restore/all
POST {{baseUrl}}/api/transaction-command/restore/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 214 POST /api/transaction-command/restore/{id}
POST {{baseUrl}}/api/transaction-command/restore/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 215 POST /api/transaction-command/trashed/{id}
POST {{baseUrl}}/api/transaction-command/trashed/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 216 POST /api/transaction-command/update/{id}
POST {{baseUrl}}/api/transaction-command/update/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 217 GET /api/transaction-query
GET {{baseUrl}}/api/transaction-query
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 218 GET /api/transaction-query/active
GET {{baseUrl}}/api/transaction-query/active
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 219 GET /api/transaction-query/merchant/{merchant_id}
GET {{baseUrl}}/api/transaction-query/merchant/1
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 220 GET /api/transaction-query/trashed
GET {{baseUrl}}/api/transaction-query/trashed
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 221 GET /api/transaction-query/{id}
GET {{baseUrl}}/api/transaction-query/1
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 222 GET /api/transaction-stats/merchant/monthly-failed
GET {{baseUrl}}/api/transaction-stats/merchant/monthly-failed
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 223 GET /api/transaction-stats/merchant/monthly-method-failed
GET {{baseUrl}}/api/transaction-stats/merchant/monthly-method-failed
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 224 GET /api/transaction-stats/merchant/monthly-method-success
GET {{baseUrl}}/api/transaction-stats/merchant/monthly-method-success
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 225 GET /api/transaction-stats/merchant/monthly-success
GET {{baseUrl}}/api/transaction-stats/merchant/monthly-success
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 226 GET /api/transaction-stats/merchant/yearly-failed
GET {{baseUrl}}/api/transaction-stats/merchant/yearly-failed
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 227 GET /api/transaction-stats/merchant/yearly-method-failed
GET {{baseUrl}}/api/transaction-stats/merchant/yearly-method-failed
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 228 GET /api/transaction-stats/merchant/yearly-method-success
GET {{baseUrl}}/api/transaction-stats/merchant/yearly-method-success
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 229 GET /api/transaction-stats/merchant/yearly-success
GET {{baseUrl}}/api/transaction-stats/merchant/yearly-success
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 230 GET /api/transaction-stats/monthly-failed
GET {{baseUrl}}/api/transaction-stats/monthly-failed
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 231 GET /api/transaction-stats/monthly-method-failed
GET {{baseUrl}}/api/transaction-stats/monthly-method-failed
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 232 GET /api/transaction-stats/monthly-method-success
GET {{baseUrl}}/api/transaction-stats/monthly-method-success
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 233 GET /api/transaction-stats/monthly-success
GET {{baseUrl}}/api/transaction-stats/monthly-success
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 234 GET /api/transaction-stats/yearly-failed
GET {{baseUrl}}/api/transaction-stats/yearly-failed
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 235 GET /api/transaction-stats/yearly-method-failed
GET {{baseUrl}}/api/transaction-stats/yearly-method-failed
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 236 GET /api/transaction-stats/yearly-method-success
GET {{baseUrl}}/api/transaction-stats/yearly-method-success
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 237 GET /api/transaction-stats/yearly-success
GET {{baseUrl}}/api/transaction-stats/yearly-success
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 238 POST /api/user-command/create
POST {{baseUrl}}/api/user-command/create
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 239 POST /api/user-command/permanent/all
POST {{baseUrl}}/api/user-command/permanent/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 240 DELETE /api/user-command/permanent/{id}
DELETE {{baseUrl}}/api/user-command/permanent/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 241 POST /api/user-command/restore/all
POST {{baseUrl}}/api/user-command/restore/all
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 242 POST /api/user-command/restore/{id}
POST {{baseUrl}}/api/user-command/restore/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 243 GET /api/user-command/trashed/{id}
GET {{baseUrl}}/api/user-command/trashed/1
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 244 POST /api/user-command/update/{id}
POST {{baseUrl}}/api/user-command/update/1
Authorization: Bearer {{accessToken}}
Content-Type: application/json

{}
HTTP *
[Asserts]
status < 500

# 245 GET /api/user-query
GET {{baseUrl}}/api/user-query
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 246 GET /api/user-query/active
GET {{baseUrl}}/api/user-query/active
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 247 GET /api/user-query/trashed
GET {{baseUrl}}/api/user-query/trashed
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500

# 248 GET /api/user-query/{id}
GET {{baseUrl}}/api/user-query/1
Authorization: Bearer {{accessToken}}
HTTP *
[Asserts]
status < 500
