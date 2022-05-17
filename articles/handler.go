package articles

import (
	errors "coodesh/error"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/robfig/cron/v3"
	"time"
)

type (
	Handler struct {
		Service service
	}

	service interface {
		getArticles(offset string, limit string) ([]FlightNews, error)
		getArticle(id string) (*FlightNews, error)
		postArticle(*FlightNews) (*FlightNews, error)
		putArticle(string, *FlightNews) (*FlightNews, error)
		deleteArticle(string) error
		populateArticles()
	}
)

func (h *Handler) GetArticles(c echo.Context) error {
	offset := c.QueryParam("offset")
	limit := c.QueryParam("limit")

	result, err := h.Service.getArticles(offset, limit)
	if err != nil {
		err := err.(errors.Message)
		c.JSON(err.HTTPCode(), err)
		return nil
	}

	c.JSON(200, result)
	return nil
}

func (h *Handler) GetArticle(c echo.Context) error {
	id := c.Param("id")

	result, err := h.Service.getArticle(id)
	if err != nil {
		err := err.(errors.Message)
		c.JSON(err.HTTPCode(), err)
		return nil
	}

	c.JSON(200, result)
	return nil
}

func (h *Handler) PostArticle(c echo.Context) error {
	body := new(FlightNews)
	if err := c.Bind(body); err != nil {
		c.JSON(400, errors.Message{
			Msg: "invalid body",
		})
		return nil
	}

	result, err := h.Service.postArticle(body)
	if err != nil {
		err := err.(errors.Message)
		c.JSON(err.HTTPCode(), err)
		return nil
	}

	c.JSON(201, result)
	return nil
}

func (h *Handler) PutArticle(c echo.Context) error {
	id := c.Param("id")
	body := new(FlightNews)

	if err := c.Bind(body); err != nil {
		c.JSON(400, errors.Message{
			Msg: "invalid body",
		})
		return nil
	}

	result, err := h.Service.putArticle(id, body)
	if err != nil {
		err := err.(errors.Message)
		c.JSON(err.HTTPCode(), err)
		return nil
	}

	c.JSON(200, result)
	return nil
}

func (h *Handler) DeleteArticle(c echo.Context) error {
	id := c.Param("id")

	if err := h.Service.deleteArticle(id); err != nil {
		err := err.(errors.Message)
		c.JSON(err.HTTPCode(), err)
		return nil
	}

	c.NoContent(200)
	return nil
}

func (h *Handler) PopulateArticles(c echo.Context) error {
	go h.Service.populateArticles()
	c.NoContent(202)
	return nil
}

func (h *Handler) ScheduledPopulateArticles() {
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	c := cron.New(cron.WithLocation(loc))

	_, err := c.AddFunc("0 9 * * *", func() {
		go h.Service.populateArticles()
	})

	if err == nil {
		c.Start()
	} else {
		fmt.Printf("Invalid cron: %s", err.Error())
	}
}
