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
        "/api/dictionary/all/{itemType}": {
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
        "/api/dictionary/clear/{itemType}": {
            "delete": {
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
        "/api/dictionary/list/{itemType}": {
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
        "/api/dictionary/one": {
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
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK - got successfully"
                    },
                    "404": {
                        "description": "Fail - not found"
                    },
                    "500": {
                        "description": "Fail - internal error"
                    }
                }
            },
            "delete": {
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
        "/api/dictionary/upsert": {
            "post": {
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
	Description:      "This is national education data dictionary backend-api server. Updated@ 2022-09-03T21:10:00+10:00",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
