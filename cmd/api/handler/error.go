package handler

import "fmt"

type clientError struct {
	status int
	body   string
}

func (e clientError) Error() string {
	return fmt.Sprintf("status : %v, body : %v", e.status, e.body)
}
