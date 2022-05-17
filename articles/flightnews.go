package articles

import (
	"encoding/json"
	"fmt"
	errors "github.com/danielbonilha/flightnews/error"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

const flightNewsURL = "https://api.spaceflightnewsapi.net/v3"

func countArticles() (int64, error) {
	fmt.Println("Counting remote articles")
	c := http.Client{}
	resp, err := c.Get(fmt.Sprintf("%s/articles/count", flightNewsURL))
	if err != nil {
		fmt.Printf("Error %s", err)
		return 0, errors.Message{
			Msg:        err.Error(),
			StatusCode: 500,
		}
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, errors.Message{
			Msg:        err.Error(),
			StatusCode: 500,
		}
	}
	i64, err := strconv.ParseInt(string(body), 10, 0)
	return i64, nil
}

func getAsyncArticles(start int, limit int, response chan []FlightNews, wg *sync.WaitGroup) {
	fmt.Printf("Counting remote articles [offset=%d][limit=%d]\n", start, limit)
	defer wg.Done()
	var items []FlightNews

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/articles", flightNewsURL), nil)
	if err != nil {
		fmt.Printf("Error %s", err)
		response <- items
		return
	}

	q := req.URL.Query()
	q.Add("_sort", "id")
	q.Add("_start", strconv.Itoa(start))
	q.Add("_limit", strconv.Itoa(limit))
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		response <- items
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&items)
	if err != nil {
		response <- items
		return
	}

	response <- items
}
