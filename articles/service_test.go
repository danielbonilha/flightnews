package articles

import (
	errors "github.com/danielbonilha/flightnews/error"
	"testing"
)

// Get Articles tests
func TestGetArticlesServiceInvalidOffset(t *testing.T) {
	svc := buildService()

	_, err := svc.getArticles("invalid", "10")
	if err == nil {
		t.Errorf("Expected not nil error")
	}
	if err.Error() != "Invalid offset" {
		t.Errorf("Expected 'Invalid offset' error")
	}
}

func TestGetArticlesServiceInvalidLimit(t *testing.T) {
	svc := buildService()

	_, err := svc.getArticles("10", "invalid")
	if err == nil {
		t.Errorf("Expected not nil error")
	}
	if err.Error() != "Invalid limit" {
		t.Errorf("Expected 'Invalid limit' error")
	}
}

func TestGetArticlesServiceZeroedLimit(t *testing.T) {
	svc := buildService()

	_, err := svc.getArticles("10", "0")
	if err == nil {
		t.Errorf("Expected not nil error")
	}
	if err.Error() != "Limit must be greater than zero" {
		t.Errorf("Expected 'Limit must be greater than zero' error")
	}
}

func TestGetArticlesServiceDatabaseError(t *testing.T) {
	svc := buildService()

	_, err := svc.getArticles("0", "99") //see mocking config - 99 return DB error
	if err == nil {
		t.Errorf("Expected not nil error")
	}
	if err.Error() != "Unexpected error" {
		t.Errorf("Expected 'Unexpected error' error")
	}
}

func TestGetArticlesServiceSuccess(t *testing.T) {
	svc := buildService()

	result, err := svc.getArticles("0", "10") //see mocking config
	if err != nil {
		t.Errorf("Expected nil error")
	}
	if len(result) != 3 {
		t.Errorf("Expected 3 items list")
	}
}

// Get Article tests
func TestGetArticleServiceInvalidId(t *testing.T) {
	svc := buildService()

	_, err := svc.getArticle("invalid")
	if err == nil {
		t.Errorf("Expected not nil error")
	}
	if err.Error() != "Invalid id" {
		t.Errorf("Expected 'Invalid id' error")
	}
}

func TestGetArticleServiceSuccess(t *testing.T) {
	svc := buildService()

	result, err := svc.getArticle("1") //see mocking config - 1 returns item
	if err != nil {
		t.Errorf("Expected nil error")
	}
	if result.Id != 1 {
		t.Errorf("Expected id == 1")
	}
}

func TestGetArticleServiceNotFound(t *testing.T) {
	svc := buildService()

	_, err := svc.getArticle("9") //see mocking config - 9 returns not found
	if err == nil {
		t.Errorf("Expected not nil error")
	}
	if err.Error() != "Item not found" {
		t.Errorf("Expected 'Item not found' error")
	}
}

func TestGetArticleServiceDatabaseError(t *testing.T) {
	svc := buildService()

	_, err := svc.getArticle("99") //see mocking config - 99 returns DB error
	if err == nil {
		t.Errorf("Expected not nil error")
	}
	if err.Error() != "Unexpected error" {
		t.Errorf("Expected 'Unexpected error' error")
	}
}

// Post Article tests
func TestPostArticleServiceInvalidId(t *testing.T) {
	svc := buildService()

	_, err := svc.postArticle(&FlightNews{Id: 0})
	if err == nil {
		t.Errorf("Expected not nil error")
	}
	if err.Error() != "Invalid id" {
		t.Errorf("Expected 'Invalid id' error")
	}
}

func TestPostArticleServiceExistingId(t *testing.T) {
	svc := buildService()

	_, err := svc.postArticle(&FlightNews{Id: 1}) //see mocking config (getArticle) - <> 9 returns existing id
	if err == nil {
		t.Errorf("Expected not nil error")
	}
	if err.Error() != "Id already exists" {
		t.Errorf("Expected 'Id already exists' error")
	}
}

func TestPostArticleServiceSuccess(t *testing.T) {
	svc := buildService()

	result, err := svc.postArticle(&FlightNews{Id: 9}) //see mocking config (getArticle) - 9 returns non-existing id
	if err != nil {
		t.Errorf("Expected nil error")
	}
	if result.Id != 9 {
		t.Errorf("Expected item with id == 9")
	}
}

// Put Article tests
func TestPutArticleServiceInvalidId(t *testing.T) {
	svc := buildService()

	_, err := svc.putArticle("invalid", &FlightNews{Id: 0})
	if err == nil {
		t.Errorf("Expected not nil error")
	}
	if err.Error() != "Invalid id" {
		t.Errorf("Expected 'Invalid id' error")
	}
}

func TestPutArticleServiceNonExistingId(t *testing.T) {
	svc := buildService()

	_, err := svc.putArticle("9", &FlightNews{Id: 9}) //see mocking config (getArticle) - 9 returns non-existing id
	if err == nil {
		t.Errorf("Expected not nil error")
	}
	if err.Error() != "Item not found" {
		t.Errorf("Expected 'Item not found' error")
	}
}

func TestPutArticleServiceMismatchId(t *testing.T) {
	svc := buildService()

	_, err := svc.putArticle("1", &FlightNews{Id: 9}) //see mocking config (getArticle) - 9 returns non-existing id
	if err == nil {
		t.Errorf("Expected not nil error")
	}
	if err.Error() != "Id mismatch" {
		t.Errorf("Expected 'Id mismatch' error")
	}
}

func TestPutArticleServiceSuccess(t *testing.T) {
	svc := buildService()

	result, err := svc.putArticle("1", &FlightNews{Id: 1}) //see mocking config (getArticle) - 9 returns non-existing id
	if err != nil {
		t.Errorf("Expected nil error")
	}
	if result.Id != 1 {
		t.Errorf("Expected item with id == 1")
	}
}

// Delete article tests
func TestDeleteArticleServiceInvalidId(t *testing.T) {
	svc := buildService()

	err := svc.deleteArticle("invalid")
	if err == nil {
		t.Errorf("Expected not nil error")
	}
	if err.Error() != "Invalid id" {
		t.Errorf("Expected 'Invalid id' error")
	}
}

func TestDeleteArticleServiceExistingId(t *testing.T) {
	svc := buildService()

	err := svc.deleteArticle("9") //see mocking config (getArticle) - <> 9 returns existing id
	if err == nil {
		t.Errorf("Expected not nil error")
	}
	if err.Error() != "Item not found" {
		t.Errorf("Expected 'Item not found' error")
	}
}

func TestDeleteArticleServiceSuccess(t *testing.T) {
	svc := buildService()

	err := svc.deleteArticle("1") //see mocking config (getArticle) - 9 returns non-existing id
	if err != nil {
		t.Errorf("Expected nil error")
	}
}

// repository mock
type mockRepository struct{}

func buildService() *Service {
	return &Service{
		Repository: &mockRepository{},
	}
}

// for mocking purposes,
// limit == 99 will return error case
// any other limit will return a success case
func (m *mockRepository) getArticles(offset int64, limit int64) ([]FlightNews, error) {
	if limit == 99 {
		return nil, errors.Message{
			Msg:        "Unexpected error",
			StatusCode: 500,
		}
	}
	return []FlightNews{{Id: 1}, {Id: 2}, {Id: 3}}, nil
}

// for mocking purposes,
// id == 9 will return 'not found' case
// id == 99 will return error case
// any other id will return a success case
func (m *mockRepository) getArticle(id int) (*FlightNews, error) {
	if id == 9 {
		return nil, errors.Message{
			Msg:        "Item not found",
			StatusCode: 404,
		}

	} else if id == 99 {
		return nil, errors.Message{
			Msg:        "Unexpected error",
			StatusCode: 500,
		}
	}

	return &FlightNews{Id: 1}, nil
}

// for mocking purposes,
// id == 99 will return error case
// any other id will return a success case
func (m *mockRepository) insertArticle(body *FlightNews) (*FlightNews, error) {
	if body.Id == 99 {
		return nil, errors.Message{
			Msg:        "Unexpected error",
			StatusCode: 500,
		}
	}

	return body, nil
}

// for mocking purposes,
// id == 99 will return error case
// any other id will return a success case
func (m *mockRepository) updateArticle(id int, body *FlightNews) (*FlightNews, error) {
	if id == 99 {
		return nil, errors.Message{
			Msg:        "Unexpected error",
			StatusCode: 500,
		}
	}

	return body, nil
}

// for mocking purposes,
// id == 99 will return error case
// any other id will return a success case
func (m *mockRepository) deleteArticle(id int) error {
	if id == 99 {
		return errors.Message{
			Msg:        "Unexpected error",
			StatusCode: 500,
		}
	}

	return nil
}

func (m *mockRepository) countArticles() (int64, error) {
	return 0, nil
}
