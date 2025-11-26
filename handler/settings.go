package handlers

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/flohoss/quiz/config"
)

func getUIOperation() huma.Operation {
	return huma.Operation{
		OperationID: "get-ui",
		Method:      http.MethodGet,
		Path:        "/api/ui",
		Summary:     "Get UI settings",
		Tags:        []string{"Settings"},
	}
}

func getUIHandler(ctx context.Context, input *struct {
}) (*struct {
	Body config.UI `json:"ui"`
}, error) {
	return &struct {
		Body config.UI `json:"ui"`
	}{Body: config.GetUI()}, nil
}
