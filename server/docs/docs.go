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
        "/api/entity/db": {
            "put": {
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Entity"
                ],
                "summary": "set mongodb database and collection for entity storage",
                "parameters": [
                    {
                        "type": "string",
                        "default": "dictionary",
                        "description": "database name",
                        "name": "database",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "entity",
                        "description": "collection name",
                        "name": "collection",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK - set db successfully"
                    },
                    "400": {
                        "description": "Fail - invalid fields"
                    }
                }
            }
        },
        "/api/entity/find": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Entity"
                ],
                "summary": "find entities json content by pass a json query string via payload",
                "parameters": [
                    {
                        "format": "binary",
                        "description": "json data for query",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK - find successfully"
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
        "/api/entity/insert/{entity}": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Entity"
                ],
                "summary": "insert or update one entity data by a json file",
                "parameters": [
                    {
                        "type": "string",
                        "description": "entity name for incoming entity data",
                        "name": "entity",
                        "in": "path",
                        "required": true
                    },
                    {
                        "format": "binary",
                        "description": "entity json data for uploading",
                        "name": "data",
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
