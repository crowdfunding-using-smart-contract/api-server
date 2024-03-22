// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/login": {
            "post": {
                "description": "Authenticate user with email and password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Authenticate User",
                "operationId": "Login",
                "parameters": [
                    {
                        "description": "User data to be authenticated",
                        "name": "User",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/UserLoginPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/ResultResponse-UserLoginResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "description": "Create user with specific user data and role",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Register User",
                "operationId": "Register",
                "parameters": [
                    {
                        "description": "User data to be created",
                        "name": "User",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/UserCreatePayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/ResultResponse-User"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/auth/renew-token": {
            "post": {
                "description": "Renew access token with refresh token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Renew Access Token",
                "operationId": "RenewAccessToken",
                "parameters": [
                    {
                        "description": "Refresh token to be renewed",
                        "name": "User",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_http_handler.RenewAccessTokenPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/ResultResponse-internal_http_handler_RenewAccessTokenResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/auth/send-verify-email": {
            "post": {
                "description": "Send verify email to user email",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Send Verify Email",
                "operationId": "SendVerifyEmail",
                "parameters": [
                    {
                        "description": "User email to be verified",
                        "name": "User",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internal_http_handler.SendVerifyEmailPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/MessageResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/auth/verify-email": {
            "get": {
                "description": "Verify email with email id and secret code",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Verify Email",
                "operationId": "VerifyEmail",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Email ID to be verified",
                        "name": "email_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Secret Code to be verified",
                        "name": "secret_code",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/ResultResponse-User"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/hello": {
            "get": {
                "description": "Health checking for the service",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "healthcheck"
                ],
                "summary": "Health Check",
                "operationId": "GetHelloMessageHandler",
                "parameters": [
                    {
                        "type": "string",
                        "description": "name of the active user",
                        "name": "name",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/MessageResponse"
                        }
                    }
                }
            }
        },
        "/projects": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Create project with required data",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "projects"
                ],
                "summary": "Create Project",
                "operationId": "CreateProject",
                "parameters": [
                    {
                        "description": "Project data to be created",
                        "name": "Project",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/fund-o_api-server_internal_entity.ProjectCreatePayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/ResultResponse-Project"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/projects/categories": {
            "get": {
                "description": "List project categories for selection",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "projects"
                ],
                "summary": "List Project Categories",
                "operationId": "ListProjectCategories",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/ResultResponse-array_ProjectCategory"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/projects/own": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get own projects with authenticate creator",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "projects"
                ],
                "summary": "Get own Projects",
                "operationId": "GetOwnProjects",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/ResultResponse-array_Project"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/projects/{id}/ratings": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Create project rating with required data",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "projects"
                ],
                "summary": "Create Project Rating",
                "operationId": "CreateProjectRating",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Project ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Project rating data to be created",
                        "name": "ProjectRating",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/fund-o_api-server_internal_entity.ProjectRatingCreatePayload"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/ResultResponse-Project"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/projects/{id}/ratings/verify": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Verify project rating by user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "projects"
                ],
                "summary": "Verify Project Rating",
                "operationId": "VerifyProjectRating",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Project ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/ResultResponse-bool"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/transactions": {
            "get": {
                "description": "Get list of transactions",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "transasctions"
                ],
                "summary": "List Transaction",
                "operationId": "ListTransactions",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "number of page",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "size of data per page",
                        "name": "size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/ResultResponse-PaginateResult-Transaction"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Create transaction with reference code",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "transasctions"
                ],
                "summary": "Create Transaction",
                "operationId": "CreateTransaction",
                "parameters": [
                    {
                        "description": "Transaction data to be created",
                        "name": "Transaction",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/TransactionCreatePayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/ResultResponse-Transaction"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/transactions/{id}": {
            "get": {
                "description": "Get transaction by id",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "transasctions"
                ],
                "summary": "Get Transaction",
                "operationId": "GetTransaction",
                "parameters": [
                    {
                        "type": "string",
                        "description": "reference code of transaction",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/ResultResponse-Transaction"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/users/me": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get current user by validating authorization token",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Get current user",
                "operationId": "GetMe",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/ResultResponse-User"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "status_code": {
                    "type": "integer"
                }
            }
        },
        "MessageResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "status_code": {
                    "type": "integer"
                }
            }
        },
        "PaginateResult-Transaction": {
            "type": "object",
            "properties": {
                "current_page": {
                    "type": "integer"
                },
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/Transaction"
                    }
                },
                "from": {
                    "type": "integer"
                },
                "last_page": {
                    "type": "integer"
                },
                "per_page": {
                    "type": "integer"
                },
                "to": {
                    "type": "integer"
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "Project": {
            "type": "object",
            "properties": {
                "category": {
                    "$ref": "#/definitions/ProjectCategory"
                },
                "created_at": {
                    "type": "string"
                },
                "current_amount": {
                    "type": "number"
                },
                "description": {
                    "type": "string"
                },
                "end_date": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "image": {
                    "type": "string"
                },
                "launch_date": {
                    "type": "string"
                },
                "location": {
                    "type": "string"
                },
                "monetary_unit": {
                    "type": "string"
                },
                "owner": {
                    "$ref": "#/definitions/User"
                },
                "rating": {
                    "type": "number"
                },
                "start_date": {
                    "type": "string"
                },
                "sub_category": {
                    "$ref": "#/definitions/ProjectSubCategory"
                },
                "sub_title": {
                    "type": "string"
                },
                "target_amount": {
                    "type": "number"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "ProjectCategory": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "subcategories": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/ProjectSubCategory"
                    }
                }
            }
        },
        "ProjectSubCategory": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "ResultResponse-PaginateResult-Transaction": {
            "type": "object",
            "properties": {
                "result": {
                    "$ref": "#/definitions/PaginateResult-Transaction"
                },
                "status": {
                    "type": "string"
                },
                "status_code": {
                    "type": "integer"
                }
            }
        },
        "ResultResponse-Project": {
            "type": "object",
            "properties": {
                "result": {
                    "$ref": "#/definitions/Project"
                },
                "status": {
                    "type": "string"
                },
                "status_code": {
                    "type": "integer"
                }
            }
        },
        "ResultResponse-Transaction": {
            "type": "object",
            "properties": {
                "result": {
                    "$ref": "#/definitions/Transaction"
                },
                "status": {
                    "type": "string"
                },
                "status_code": {
                    "type": "integer"
                }
            }
        },
        "ResultResponse-User": {
            "type": "object",
            "properties": {
                "result": {
                    "$ref": "#/definitions/User"
                },
                "status": {
                    "type": "string"
                },
                "status_code": {
                    "type": "integer"
                }
            }
        },
        "ResultResponse-UserLoginResponse": {
            "type": "object",
            "properties": {
                "result": {
                    "$ref": "#/definitions/UserLoginResponse"
                },
                "status": {
                    "type": "string"
                },
                "status_code": {
                    "type": "integer"
                }
            }
        },
        "ResultResponse-array_Project": {
            "type": "object",
            "properties": {
                "result": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/Project"
                    }
                },
                "status": {
                    "type": "string"
                },
                "status_code": {
                    "type": "integer"
                }
            }
        },
        "ResultResponse-array_ProjectCategory": {
            "type": "object",
            "properties": {
                "result": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/ProjectCategory"
                    }
                },
                "status": {
                    "type": "string"
                },
                "status_code": {
                    "type": "integer"
                }
            }
        },
        "ResultResponse-bool": {
            "type": "object",
            "properties": {
                "result": {
                    "type": "boolean"
                },
                "status": {
                    "type": "string"
                },
                "status_code": {
                    "type": "integer"
                }
            }
        },
        "ResultResponse-internal_http_handler_RenewAccessTokenResponse": {
            "type": "object",
            "properties": {
                "result": {
                    "$ref": "#/definitions/internal_http_handler.RenewAccessTokenResponse"
                },
                "status": {
                    "type": "string"
                },
                "status_code": {
                    "type": "integer"
                }
            }
        },
        "Transaction": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "ref_code": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "TransactionCreatePayload": {
            "type": "object",
            "required": [
                "ref_code"
            ],
            "properties": {
                "ref_code": {
                    "type": "string"
                }
            }
        },
        "User": {
            "type": "object",
            "properties": {
                "birthdate": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "full_name": {
                    "type": "string"
                },
                "gender": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "is_email_verified": {
                    "type": "boolean"
                },
                "profile_image": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "UserCreatePayload": {
            "type": "object",
            "required": [
                "birthdate",
                "email",
                "firstname",
                "gender",
                "lastname",
                "password",
                "password_confirmation"
            ],
            "properties": {
                "birthdate": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "firstname": {
                    "type": "string"
                },
                "gender": {
                    "type": "string"
                },
                "lastname": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "password_confirmation": {
                    "type": "string"
                }
            }
        },
        "UserLoginPayload": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "UserLoginResponse": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "access_token_expired_at": {
                    "type": "string"
                },
                "refresh_token": {
                    "type": "string"
                },
                "refresh_token_expired_at": {
                    "type": "string"
                },
                "session_id": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/User"
                }
            }
        },
        "fund-o_api-server_internal_entity.ProjectCreatePayload": {
            "type": "object"
        },
        "fund-o_api-server_internal_entity.ProjectRatingCreatePayload": {
            "type": "object",
            "required": [
                "rating"
            ],
            "properties": {
                "rating": {
                    "type": "number",
                    "maximum": 5,
                    "minimum": 0
                }
            }
        },
        "internal_http_handler.RenewAccessTokenPayload": {
            "type": "object",
            "required": [
                "refresh_token"
            ],
            "properties": {
                "refresh_token": {
                    "type": "string"
                }
            }
        },
        "internal_http_handler.RenewAccessTokenResponse": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "access_token_expired_at": {
                    "type": "string"
                }
            }
        },
        "internal_http_handler.SendVerifyEmailPayload": {
            "type": "object",
            "required": [
                "email"
            ],
            "properties": {
                "email": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{"http", "https"},
	Title:            "FundO API",
	Description:      "This is a sample server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
