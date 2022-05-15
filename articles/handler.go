package articles

import (
	errors "coodesh/error"
	"github.com/labstack/echo/v4"
)

type (
	Handler struct {
		Service service
	}

	service interface {
		getArticles() ([]*FlightNews, error)
		getArticle(string) (*FlightNews, error)
		postArticle(*FlightNews) (*FlightNews, error)
		putArticle(string, *FlightNews) (*FlightNews, error)
		deleteArticle(string) error
	}
)

func (h *Handler) GetArticles(c echo.Context) error {
	result, err := h.Service.getArticles()
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
