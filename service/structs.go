package service

import (
	"encoding/json"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

type ErrorHandler struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type GetRequest struct {
	Url      string
	Response any
}

func (gr *GetRequest) Request(params []string) any {
	url := strings.Join(append([]string{gr.Url}, params...), "/")
	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("request sending error! err: %s", err)
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusOK {
		err = json.NewDecoder(response.Body).Decode(gr.Response)
		if err != nil {
			log.Fatalf("response decoding error! err: %s", err)
		}
		return gr.Response
	}
	if response.Status == "404 Not Found" {
		return &ErrorHandler{
			Code:    404,
			Message: "Not found",
		}
	}
	errHandler := &ErrorHandler{}
	err = json.NewDecoder(response.Body).Decode(errHandler)
	if err != nil {
		log.Fatalf("response decoding error! err: %s", err)
	}
	return errHandler
}

type PostRequest struct {
	Url      string
	Response any
}

func (gr *PostRequest) Request() any {
	return gr.Response
}
