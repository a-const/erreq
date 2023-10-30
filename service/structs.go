package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

type ErrorHandler struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func mustWrite(builder *strings.Builder, str string) {
	_, err := builder.WriteString(str)
	if err != nil {
		log.Fatalf("error writing to string builder! err: %s", err)
	}
}

type GetRequest struct {
	Url      string
	Endpoint string
	Response any
}

func (gr *GetRequest) Request(params []string, port string) any {
	urlBuilder := new(strings.Builder)
	mustWrite(urlBuilder, gr.Url)
	mustWrite(urlBuilder, port)
	mustWrite(urlBuilder, gr.Endpoint)
	url := urlBuilder.String()
	url = strings.Join(append([]string{url}, params...), "/")
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
	Endpoint string
	Response any
}

func (pr *PostRequest) Request(body any, port string) any {
	urlBuilder := new(strings.Builder)
	mustWrite(urlBuilder, pr.Url)
	mustWrite(urlBuilder, port)
	mustWrite(urlBuilder, pr.Endpoint)
	url := urlBuilder.String()
	marshaledBody, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("body encoding error! err: %s", err)
	}
	response, err := http.Post(url, "application/json", bytes.NewBuffer(marshaledBody))
	if err != nil {
		log.Fatalf("request sending error! err: %s", err)
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusOK {
		err = json.NewDecoder(response.Body).Decode(pr.Response)
		if err != nil {
			log.Fatalf("response decoding error! err: %s", err)
		}
		return pr.Response
	}

	return &ErrorHandler{
		Code:    response.StatusCode,
		Message: "Error",
	}
}
