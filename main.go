package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Book struct {
	ISBN      string `json:"pk"`
	Author    string `json:"sk"`
	Title     string `json:"title"`
	ItemCount int    `json:"itemCount"`
}

func router(req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	switch req.RequestContext.HTTP.Method {
	case "GET":
		return show(req)
	case "POST":
		return create(req)
	case "PATCH":
		return update(req)
	default:
		return clientError(http.StatusMethodNotAllowed)
	}
}

// GET request must use pk,sk as JSON
func show(req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// Unmarshal Request JSON to Book struct
	bookReq, err := unmarshalBookJson(req)
	if err != nil {
		return clientError(http.StatusUnprocessableEntity)
	}

	// Validate Request, see utils for details
	if !validateReadRequest(bookReq.ISBN, bookReq.Author) {
		return clientError(http.StatusBadRequest)
	}

	// Get the book response from DynamoDB based on the pk,sk pair
	bookRes, err := getItem(bookReq.ISBN, bookReq.Author)
	if err != nil {
		return serverError(err)
	}
	if bookRes == nil {
		return clientError(http.StatusNotFound)
	}

	// APIGateway Body needs to be JSON, so we convert here
	js, err := json.Marshal(bookRes)
	if err != nil {
		return serverError(err)
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       string(js),
	}, nil
}

// POST request must use Book fields as JSON
func create(req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// Validate header content-type is JSON
	if !validateJsonFormat(req) {
		return clientError(http.StatusNotAcceptable)
	}

	// Unmarshal Request JSON to Book struct
	bookReq, err := unmarshalBookJson(req)
	if err != nil {
		return clientError(http.StatusUnprocessableEntity)
	}

	// Validate Request, see utils for details
	if !validateWriteRequest(bookReq.ISBN, bookReq.Author, bookReq.Title) {
		return clientError(http.StatusBadRequest)
	}

	// putItem returns an error (normally will be nil)
	err = putItem(bookReq)
	if err != nil {
		return serverError(err)
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 201,
		Headers:    map[string]string{"Location": fmt.Sprintf("/books?pk=%s&sk=%s", bookReq.ISBN, bookReq.Author)},
		Body:       req.Body,
	}, nil
}

// PATCH request must use Book fields as JSON
func update(req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// Unmarshal Request JSON to Book struct
	bookReq, err := unmarshalBookJson(req)
	if err != nil {
		return clientError(http.StatusUnprocessableEntity)
	}

	// Validate Request, see utils for details
	if !validateReadRequest(bookReq.ISBN, bookReq.Author) {
		return clientError(http.StatusBadRequest)
	}

	// Get the book response from DynamoDB based on the pk,sk pair
	bookRes, err := incrementItem(bookReq.ISBN, bookReq.Author)
	if err != nil {
		return serverError(err)
	}
	if bookRes == nil {
		return clientError(http.StatusNotFound)
	}

	// APIGateway Body needs to be JSON, so we convert here
	js, err := json.Marshal(bookRes)
	if err != nil {
		return serverError(err)
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       string(js),
	}, nil
}

func main() {
	lambda.Start(router)
}
