package articles

import (
	errors "coodesh/error"
	_ "errors"
	"fmt"
	"strconv"
	"sync"
)

type (
	Service struct {
		Repository repository
	}

	repository interface {
		getArticles(offset int64, limit int64) ([]FlightNews, error)
		getArticle(int) (*FlightNews, error)
		insertArticle(*FlightNews) (*FlightNews, error)
		updateArticle(int, *FlightNews) (*FlightNews, error)
		deleteArticle(int) error
		countArticles() (int64, error)
	}
)

func (s *Service) getArticles(offset string, limit string) ([]FlightNews, error) {
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
		return s.Repository.insertArticle(body)
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

	return s.Repository.updateArticle(intId, body)
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
	localCount, err := s.Repository.countArticles()
	if err != nil {
		return
	}

	remoteCount, err := countArticles()
	if err != nil {
		return
	}

	if remoteCount > localCount {
		s.syncArticles(localCount, remoteCount)
	} else {
		fmt.Println("No new articles")
	}
}

func (s *Service) syncArticles(localCount, remoteCount int64) {
	remoteItemsCh := make(chan []FlightNews)
	localItemsCh := make(chan []FlightNews)

	defer close(remoteItemsCh)

	go s.getAsyncRemoteItems(remoteCount, remoteItemsCh)
	go s.getAsyncLocalItems(localCount, localItemsCh)

	localItems := <-localItemsCh
	remoteItems := <-remoteItemsCh

	localItemsReferenceMap := buildLocalItemsReference(localItems)
	newItems := filterNonPersistedRemoteItems(remoteItems, localItemsReferenceMap)

	fmt.Printf("Persisting %d items", len(newItems))
	for _, item := range newItems {
		s.Repository.insertArticle(&item)
	}
}

func (s *Service) getAsyncRemoteItems(remoteCount int64, remoteItemsCh chan []FlightNews) {
	start := 0
	limit := 1000
	iterations := remoteCount/int64(limit) + 1

	itemsCh := make(chan []FlightNews, iterations)
	defer close(itemsCh)

	var wg sync.WaitGroup
	wg.Add(int(iterations))

	var items []FlightNews
	for i := int64(0); i < iterations; i++ {
		go getAsyncArticles(start, limit, itemsCh, &wg)
		start += limit
	}

	wg.Wait()

	for i := int64(0); i < iterations; i++ {
		items = append(items, <-itemsCh...)
	}

	remoteItemsCh <- items
}

func (s *Service) getAsyncLocalItems(localCount int64, localItemsCh chan []FlightNews) {
	start := int64(0)
	limit := int64(1000)
	iterations := localCount/limit + 1

	itemsCh := make(chan []FlightNews, iterations)
	defer close(itemsCh)

	var wg sync.WaitGroup
	wg.Add(int(iterations))

	var items []FlightNews
	for i := int64(0); i < iterations; i++ {
		go func(start int64, limit int64) {
			result, _ := s.Repository.getArticles(start, limit)
			itemsCh <- result
			wg.Done()
		}(start, limit)
		start += limit
	}

	wg.Wait()

	for i := int64(0); i < iterations; i++ {
		items = append(items, <-itemsCh...)
	}

	localItemsCh <- items
}

func buildLocalItemsReference(items []FlightNews) map[int]*FlightNews {
	itemsMap := make(map[int]*FlightNews)
	for _, item := range items {
		itemsMap[item.Id] = &item
	}
	return itemsMap
}

func filterNonPersistedRemoteItems(remoteItems []FlightNews, localItemsReference map[int]*FlightNews) []FlightNews {
	newItems := make([]FlightNews, 0)
	for _, remoteItem := range remoteItems {
		if localItemsReference[remoteItem.Id] == nil {
			newItems = append(newItems, remoteItem)
		}
	}
	return newItems
}
