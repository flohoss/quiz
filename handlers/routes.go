package handlers

import (
	"html/template"
	"io"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/flohoss/quiz/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func initTemplates() *Template {
	return &Template{
		templates: template.Must(template.ParseGlob("web/index.html")),
	}
}

func longCacheLifetime(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderCacheControl, "public, max-age=31536000")
		return next(c)
	}
}

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

func SetupRouter() *echo.Echo {
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Path(), "events")
		},
	}))

	h := huma.DefaultConfig("Quiz API", "1.0.0")
	h.OpenAPIPath = "/api/openapi"
	h.SchemasPath = "/api/schemas"
	he := humaecho.New(e, h)

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
	e.Renderer = initTemplates()

	huma.Register(he, getAppOperation(), getAppHandler)

	he.UseMiddleware(languageValidationMiddleware(he))
	huma.Register(he, getQuestionsOperation(), getQuestionsHandler)
	huma.Register(he, validateAnswersOperation(), validateAnswersHandler)

	logo := config.GetApp().Logo
	e.File(logo, logo, longCacheLifetime)
	icons := config.GetApp().Icons
	for _, path := range icons {
		e.File(path, path, longCacheLifetime)
	}

	e.GET("/robots.txt", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "User-agent: *\nDisallow: /")
	})

	assets := e.Group("/assets", longCacheLifetime)
	assets.Static("/", "web/assets")

	favicon := e.Group("/static", longCacheLifetime)
	favicon.Static("/", "web/static")

	e.RouteNotFound("*", func(ctx echo.Context) error {
		return ctx.Render(http.StatusOK, "index.html", nil)
	})

	return e
}
