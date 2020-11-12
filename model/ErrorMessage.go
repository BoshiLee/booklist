package model

type ErrorMessage struct {
	Message string `json:"message"`
}

func (e ErrorMessage) Error() string {
	return e.Message
}
