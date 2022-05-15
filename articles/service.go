package articles

import (
	errors "coodesh/error"
	_ "errors"
	"strconv"
)

type (
	Service struct {
		Repository repository
	}

	repository interface {
		getArticles() ([]*FlightNews, error)
		getArticle(int) (*FlightNews, error)
		postArticle(*FlightNews) (*FlightNews, error)
		putArticle(int, *FlightNews) (*FlightNews, error)
		deleteArticle(int) error
	}
)

func (s *Service) getArticles() ([]*FlightNews, error) {
	return s.Repository.getArticles()
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
			Msg:        "Invalid Id",
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