# Category domain E2E suite (monolith-ecommerce-grpc).
# NOTE: category create/update accept multipart form data with an image file.
# Run with:
#   hurl --test --jobs 1 --variable baseUrl=http://127.0.0.1:5000 \
#     --variable testEmail=... --variable testPassword=... \
#     --variable testImage=... tests/hurl/category.hurl

# --- Auth bootstrap ---------------------------------------------------------
POST {{baseUrl}}/api/auth/register
Content-Type: application/json

{
    "firstname": "Hurl",
    "lastname": "Category",
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

# --- Create category (multipart + image) ------------------------------------
POST {{baseUrl}}/api/category-command/create
Authorization: Bearer {{accessToken}}
[MultipartFormData]
name: Hurl Category {{testEmail}}
description: Category created by hurl E2E test
slug_category: hurl-category-{{testEmail}}
image: file,assets/test.png; image/png
HTTP 201
[Captures]
categoryId: jsonpath "$.data.id"
[Asserts]
jsonpath "$.status" == "success"

# --- Queries ----------------------------------------------------------------
GET {{baseUrl}}/api/category-query?page=1&page_size=10
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/category-query/active?page=1&page_size=10
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/category-query/{{categoryId}}
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"
jsonpath "$.data.id" == {{categoryId}}

# --- Update (multipart + image) ---------------------------------------------
POST {{baseUrl}}/api/category-command/update/{{categoryId}}
Authorization: Bearer {{accessToken}}
[MultipartFormData]
name: Hurl Category Updated
description: Category updated by hurl E2E test
slug_category: hurl-category-updated-{{testEmail}}
image: file,assets/test.png; image/png
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

# --- Trash / restore --------------------------------------------------------
POST {{baseUrl}}/api/category-command/trashed/{{categoryId}}
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

GET {{baseUrl}}/api/category-query/trashed?page=1&page_size=10
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

POST {{baseUrl}}/api/category-command/restore/{{categoryId}}
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

# --- Permanent delete only applies to trashed records -----------------------
POST {{baseUrl}}/api/category-command/trashed/{{categoryId}}
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"

DELETE {{baseUrl}}/api/category-command/permanent/{{categoryId}}
Authorization: Bearer {{accessToken}}
HTTP 200
[Asserts]
jsonpath "$.status" == "success"
