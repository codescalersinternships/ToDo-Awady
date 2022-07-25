package docs

import (
	_ "embed"

	"github.com/swaggo/swag"
)

//go:embed swagger.json
var docTemplate []byte

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  string(docTemplate[:]),
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
