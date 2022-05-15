package articles

type (
	FlightNews struct {
		Id          int        `json:"id"`
		Featured    bool       `json:"featured"`
		Title       string     `json:"title"`
		Url         string     `json:"url"`
		ImageUrl    string     `json:"imageUrl"`
		NewsSite    string     `json:"newsSite"`
		Summary     string     `json:"summary"`
		PublishedAt string     `json:"publishedAt"`
		Launches    []Provider `json:"launches"`
		Events      []Provider `json:"events"`
	}

	Provider struct {
		Id       string `json:"id"`
		Provider string `json:"provider"`
	}
)
