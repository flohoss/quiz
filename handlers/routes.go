package handlers

import (
	"html/template"
	"io"
	"net/http"
	"os"

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

func InitRouter() *echo.Echo {
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Gzip())
	e.Pre(middleware.RemoveTrailingSlash())

	e.Renderer = initTemplates()

	return e
}

func SetupRouter(e *echo.Echo) {
	h := huma.DefaultConfig("Quiz API", os.Getenv("APP_VERSION"))
	h.OpenAPIPath = "/api/openapi"
	h.DocsPath = "/api/docs"
	h.SchemasPath = "/api/schemas"
	humaAPI := humaecho.New(e, h)

	huma.Register(humaAPI, getAppOperation(), getAppHandler)

	humaAPI.UseMiddleware(languageValidationMiddleware(humaAPI))
	huma.Register(humaAPI, getQuestionsOperation(), getQuestionsHandler)
	huma.Register(humaAPI, validateAnswersOperation(), validateAnswersHandler)

	logo := config.GetApp().Logo
	e.File(logo, logo, longCacheLifetime)
	favicon := config.GetApp().Favicon
	e.File(favicon, favicon, longCacheLifetime)

	e.GET("/robots.txt", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "User-agent: *\nDisallow: /")
	})

	assets := e.Group("/assets", longCacheLifetime)
	assets.Static("/", "web/assets")

	e.RouteNotFound("*", func(ctx echo.Context) error {
		return ctx.Render(http.StatusOK, "index.html", nil)
	})
}
