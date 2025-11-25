package handlers

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/flohoss/christmas/config"
	"github.com/labstack/echo/v4"
)

func languageValidationMiddleware(api huma.API) func(ctx huma.Context, next func(ctx huma.Context)) {
	return func(ctx huma.Context, next func(ctx huma.Context)) {
		lang := ctx.Query("language")
		if err := config.ValidateLanguage(lang); err != nil {
			huma.WriteErr(api, ctx, http.StatusBadRequest, "Invalid language", err)
			return
		}
		next(ctx)
	}
}

func SetupRouter(e *echo.Echo) {
	config := huma.DefaultConfig("Christmas API", "1.0.0")
	config.OpenAPIPath = "/api/openapi"
	config.SchemasPath = "/api/schemas"
	h := humaecho.New(e, config)

	e.GET("/api/docs", func(ctx echo.Context) error {
		return ctx.HTML(http.StatusOK, `<!doctype html>
			<html>
				<head>
					<title>API Reference</title>
					<meta charset="utf-8" />
					<meta name="viewport" content="width=device-width, initial-scale=1" />
				</head>
				<body>
					<script id="api-reference" data-url="/api/openapi.json"></script>
					<script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
				</body>
			</html>`,
		)
	})

	h.UseMiddleware(languageValidationMiddleware(h))
	huma.Register(h, getQuestionsOperation(), getQuestionsHandler)
	huma.Register(h, validateAnswersOperation(), validateAnswersHandler)
}
