package api

import (
	_ "github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen" // to avoid go mod tidy clean it up
)

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=codegen-config.yaml api.yaml
