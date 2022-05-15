package main

import (
	"coodesh/articles"
	"coodesh/mongodb"
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	e := echo.New()
	configRoutes(e)
	e.Logger.Fatal(e.Start(":1323"))
}

func configRoutes(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Back-end Challenge 2021 🏅 - Space Flight News")
	})

	handler := buildHandler()

	e.GET("/articles", handler.GetArticles)
	e.GET("/articles/:id", handler.GetArticle)
	e.POST("/articles", handler.PostArticle)
	e.PUT("/articles/:id", handler.PutArticle)
	e.DELETE("/articles/:id", handler.DeleteArticle)
}

func buildHandler() articles.Handler {
	return articles.Handler{
		Service: &articles.Service{
			Repository: &articles.NoSqlRepository{
				Conn: mongodb.Connect(),
			},
		},
	}
}
