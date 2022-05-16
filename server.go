package main

import (
	"coodesh/articles"
	"coodesh/mongodb"
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	e := echo.New()
	handler := buildHandler()

	configRoutes(e, handler)
	handler.ScheduledPopulateArticles()

	e.Logger.Fatal(e.Start(":1323"))
}

func configRoutes(e *echo.Echo, handler articles.Handler) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Back-end Challenge 2021 🏅 - Space Flight News")
	})

	e.GET("/articles", handler.GetArticles)
	e.GET("/articles/:id", handler.GetArticle)
	e.POST("/articles", handler.PostArticle)
	e.PUT("/articles/:id", handler.PutArticle)
	e.DELETE("/articles/:id", handler.DeleteArticle)

	// this endpoint will run the "populate articles" method.
	// The same method that will run via CRON every day at 9am
	e.GET("/articles/populate", handler.PopulateArticles)
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
