// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

// SwaggerJSON embedded version of the swagger document used at generation time
var SwaggerJSON json.RawMessage

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "produces": [
    "application/json"
  ],
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "title": "oauth2 debug",
    "version": "0.3.0"
  },
  "basePath": "/api",
  "paths": {
    "/auth/callback": {
      "get": {
        "summary": "return access_token",
        "responses": {
          "200": {
            "description": "login",
            "schema": {
              "properties": {
                "access_token": {
                  "type": "string",
                  "format": "string"
                }
              }
            }
          },
          "default": {
            "description": "error",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    },
    "/customers": {
      "get": {
        "tags": [
          "customers"
        ],
        "summary": "Get a customerId given an SSN",
        "operationId": "getId",
        "parameters": [
          {
            "name": "info",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/social_id"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/customer"
            }
          },
          "401": {
            "description": "unauthorized",
            "schema": {
              "$ref": "#/definitions/error"
            }
          },
          "404": {
            "description": "resource not found",
            "schema": {
              "$ref": "#/definitions/error"
            }
          },
          "default": {
            "description": "error",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      },
      "post": {
        "tags": [
          "customers"
        ],
        "summary": "Create a new customer to track",
        "operationId": "create",
        "parameters": [
          {
            "name": "info",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/customer"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "created",
            "schema": {
              "$ref": "#/definitions/customer"
            }
          },
          "default": {
            "description": "error",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    },
    "/login": {
      "get": {
        "summary": "login through oauth2 server",
        "responses": {
          "200": {
            "description": "login",
            "schema": {
              "properties": {
                "access_token": {
                  "type": "string",
                  "format": "string"
                }
              }
            }
          },
          "default": {
            "description": "error",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "customer": {
      "type": "object",
      "required": [
        "customerId",
        "name",
        "surname",
        "ssn",
        "fipsCode"
      ],
      "properties": {
        "agentId": {
          "description": "agent associated with this customer",
          "type": "integer",
          "format": "int32"
        },
        "customerId": {
          "description": "internal identifier of a customer",
          "type": "integer",
          "format": "int64",
          "readOnly": true
        },
        "fipsCode": {
          "type": "string",
          "format": "string",
          "minLength": 1
        },
        "name": {
          "type": "string",
          "format": "string",
          "minLength": 1
        },
        "ssn": {
          "description": "Lookup identifier to find a customer in the system",
          "type": "string",
          "format": "string",
          "minLength": 11
        },
        "surname": {
          "type": "string",
          "format": "string",
          "minLength": 1
        }
      }
    },
    "error": {
      "type": "object",
      "required": [
        "message"
      ],
      "properties": {
        "code": {
          "type": "integer",
          "format": "int64"
        },
        "fields": {
          "type": "string"
        },
        "message": {
          "type": "string"
        }
      }
    },
    "principal": {
      "type": "string"
    },
    "social_id": {
      "type": "object",
      "required": [
        "ssn"
      ],
      "properties": {
        "ssn": {
          "type": "string",
          "format": "string",
          "minLength": 11
        }
      }
    }
  },
  "securityDefinitions": {
    "OauthSecurity": {
      "type": "oauth2",
      "flow": "accessCode",
      "authorizationUrl": "https://accounts.google.com/o/oauth2/v2/auth",
      "tokenUrl": "hhttps://www.googleapis.com/oauth2/v4/token",
      "scopes": {
        "admin": "Admin scope",
        "user": "User scope"
      }
    }
  },
  "security": [
    {
      "OauthSecurity": [
        "user"
      ]
    }
  ]
}`))
}
