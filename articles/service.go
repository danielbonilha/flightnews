package articles

import (
	errors "coodesh/error"
	_ "errors"
	"fmt"
	"strconv"
)

type (
	Service struct {
		Repository repository
	}

	repository interface {
		getArticles(offset int64, limit int64) ([]*FlightNews, error)
		getArticle(int) (*FlightNews, error)
		postArticle(*FlightNews) (*FlightNews, error)
		putArticle(int, *FlightNews) (*FlightNews, error)
		deleteArticle(int) error
	}
)

func (s *Service) getArticles(offset string, limit string) ([]*FlightNews, error) {
	if offset == "" {
		offset = "0"
	}

	if limit == "" {
		limit = "10"
	}

	intOffset, err := strconv.ParseInt(offset, 10, 64)
	if err != nil || intOffset < 0 {
		return nil, errors.Message{
			Msg:        "Invalid offset",
			StatusCode: 400,
		}
	}

	intLimit, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return nil, errors.Message{
			Msg:        "Invalid limit",
			StatusCode: 400,
		}
	}
	if intLimit < 1 {
		return nil, errors.Message{
			Msg:        "Limit must be greater than zero",
			StatusCode: 400,
		}
	}

	return s.Repository.getArticles(intOffset, intLimit)
}

func (s *Service) getArticle(id string) (*FlightNews, error) {
	intId, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.Message{
			Msg:        "Invalid id",
			StatusCode: 400,
		}
	}

	return s.Repository.getArticle(intId)
}

func (s *Service) postArticle(body *FlightNews) (*FlightNews, error) {
	if body.Id == 0 {
		return nil, errors.Message{
			Msg:        "Invalid id",
			StatusCode: 400,
		}
	}

	_, err := s.Repository.getArticle(body.Id)
	if err != nil && err.Error() == "Item not found" {
		return s.Repository.postArticle(body)
	}

	return nil, errors.Message{
		Msg:        "Id already exists",
		StatusCode: 400,
	}
}

func (s *Service) putArticle(id string, body *FlightNews) (*FlightNews, error) {
	intId, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.Message{
			Msg:        "Invalid id",
			StatusCode: 400,
		}
	}

	_, err = s.Repository.getArticle(intId)
	if err != nil {
		return nil, err
	}

	if body.Id != intId {
		return nil, errors.Message{
			Msg:        "Id mismatch",
			StatusCode: 400,
		}
	}

	return s.Repository.putArticle(intId, body)
}

func (s *Service) deleteArticle(id string) error {
	intId, err := strconv.Atoi(id)
	if err != nil {
		return errors.Message{
			Msg:        "Invalid id",
			StatusCode: 400,
		}
	}

	_, err = s.Repository.getArticle(intId)
	if err != nil {
		return err
	}

	return s.Repository.deleteArticle(intId)
}

func (s *Service) populateArticles() {
	fmt.Println("Will populate all articles at 9am BRT or after GET /articles/populate")
}
