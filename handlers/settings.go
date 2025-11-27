package handlers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/flohoss/quiz/config"
)

func getAppOperation() huma.Operation {
	return huma.Operation{
		OperationID: "get-app",
		Method:      http.MethodGet,
		Path:        "/api/app",
		Summary:     "Get App settings",
		Tags:        []string{"Settings"},
	}
}

func getAppHandler(ctx context.Context, input *struct {
}) (*struct {
	Body config.AppSettings `json:"app"`
}, error) {
	return &struct {
		Body config.AppSettings `json:"app"`
	}{Body: config.GetApp()}, nil
}
