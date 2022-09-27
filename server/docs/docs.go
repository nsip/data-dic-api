// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "API Support"
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
        "/api/admin/users": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "summary": "get all users' info",
                "parameters": [
                    {
                        "type": "string",
                        "description": "user filter with uname wildcard(*)",
                        "name": "uname",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "user filter with name wildcard(*)",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "user filter with active status",
                        "name": "active",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK - list successfully"
                    },
                    "401": {
                        "description": "Fail - unauthorized error"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        },
        "/api/dictionary/auth/clear/{itemType}": {
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Dictionary"
                ],
                "summary": "delete all entities or collections (dangerous)",
                "parameters": [
                    {
                        "type": "string",
                        "description": "item type, only can be 'entity' or 'collection'",
                        "name": "itemType",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK - cleared successfully"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        },
        "/api/dictionary/auth/one": {
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Dictionary"
                ],
                "summary": "delete one entity or collection by its 'Entity' name",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Entity name for deleting",
                        "name": "name",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK - deleted successfully"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        },
        "/api/dictionary/auth/upsert": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Dictionary"
                ],
                "summary": "insert or update one entity or collection data by json payload",
                "parameters": [
                    {
                        "format": "binary",
                        "description": "entity or collection json data for uploading",
                        "name": "entity",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK - insert or update successfully"
                    },
                    "400": {
                        "description": "Fail - invalid parameters or request body"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        },
        "/api/dictionary/pub/colentities": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Dictionary"
                ],
                "summary": "get related entities' name of a collection",
                "parameters": [
                    {
                        "type": "string",
                        "description": "collection name",
                        "name": "colname",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK - got collection content successfully"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        },
        "/api/dictionary/pub/entclasses": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Dictionary"
                ],
                "summary": "get class info (derived path \u0026 children) of an entity",
                "parameters": [
                    {
                        "type": "string",
                        "description": "entity name",
                        "name": "entname",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK - got entity class info successfully"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        },
        "/api/dictionary/pub/items/{itemType}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Dictionary"
                ],
                "summary": "get all entities' or collections' full content",
                "parameters": [
                    {
                        "type": "string",
                        "description": "item type, only can be 'entity' or 'collection'",
                        "name": "itemType",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "entity/collection 'Entity' name for query. if empty, get all",
                        "name": "name",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK - get successfully"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        },
        "/api/dictionary/pub/kind": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Dictionary"
                ],
                "summary": "check item's kind ('entity' or 'collection') by its name",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Entity name for checking kind",
                        "name": "name",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK - got kind ('entity' or 'collection') successfully"
                    },
                    "404": {
                        "description": "Fail - neither 'entity' nor 'collection'"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        },
        "/api/dictionary/pub/list/{itemType}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Dictionary"
                ],
                "summary": "list all entities' or collections' name",
                "parameters": [
                    {
                        "type": "string",
                        "description": "item type, only can be 'entity' or 'collection'",
                        "name": "itemType",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "entity/collection 'Entity' name for query. if empty, get all",
                        "name": "name",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK - list successfully"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        },
        "/api/dictionary/pub/one": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Dictionary"
                ],
                "summary": "get one entity or collection by its 'Entity' name",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Entity name",
                        "name": "name",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "boolean",
                        "description": "regex applies?",
                        "name": "fuzzy",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK - got successfully"
                    },
                    "400": {
                        "description": "Fail - invalid parameters"
                    },
                    "404": {
                        "description": "Fail - not found"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        },
        "/api/dictionary/pub/search": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Dictionary"
                ],
                "summary": "get list of entity's \u0026 collection's name by searching. If not given, return all",
                "parameters": [
                    {
                        "type": "string",
                        "description": "search content from whole dictionary",
                        "name": "aim",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "case insensitive ?",
                        "name": "ignorecase",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK - got list of found item's name successfully"
                    },
                    "400": {
                        "description": "Fail - invalid parameters"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        },
        "/api/user/sign-in": {
            "post": {
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "sign in action. if ok, got token",
                "parameters": [
                    {
                        "type": "string",
                        "description": "user name or email",
                        "name": "uname",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "format": "password",
                        "description": "password",
                        "name": "pwd",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK - sign-in successfully"
                    },
                    "400": {
                        "description": "Fail - incorrect password"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        },
        "/api/user/sign-out": {
            "put": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "sign out action.",
                "responses": {
                    "200": {
                        "description": "OK - sign-out successfully"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        },
        "/api/user/sign-up": {
            "post": {
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "sign up action, send user's basic info for registry",
                "parameters": [
                    {
                        "type": "string",
                        "description": "unique user name",
                        "name": "uname",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "format": "email",
                        "description": "user's email",
                        "name": "email",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "user's password",
                        "name": "pwd",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK - then waiting for verification code"
                    },
                    "400": {
                        "description": "Fail - invalid registry fields"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "127.0.0.1:1323",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "National Education Data Dictionary API",
	Description:      "This is national education data dictionary backend-api server. Updated@ 2022-09-15T09:29:03+10:00",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
