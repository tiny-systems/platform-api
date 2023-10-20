package api

import (
	_ "github.com/deepmap/oapi-codegen/pkg/codegen" // to avoid go mod tidy clean it up
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=codegen-config.yaml api.yaml
