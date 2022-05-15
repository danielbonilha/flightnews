package errors

type Message struct {
	Msg        string `json:"message"`
	StatusCode int    `json:"-"`
}

func (e Message) Error() string {
	return e.Msg
}

func (e Message) HTTPCode() int {
	return e.StatusCode
}
