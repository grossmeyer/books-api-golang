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
	// Use unmarshal to map JSON to book struct
	bookReq := new(Book)
	err := json.Unmarshal([]byte(req.Body), bookReq)
	if err != nil {
		return clientError(http.StatusUnprocessableEntity)
	}

	// Get the ISBN from the pk param and validate
	pk := bookReq.ISBN
	if !isbnRegexp.MatchString(pk) {
		return clientError(http.StatusBadRequest)
	}

	// Get the Author from the sk param
	sk := bookReq.Author
	if sk == "" {
		return clientError(http.StatusBadRequest)
	}

	// Get the book response from DynamoDB based on the pk,sk pair
	bookRes, err := getItem(pk, sk)
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
	if !checkJsonFormat(req) {
		return clientError(http.StatusNotAcceptable)
	}

	// Use unmarshal to map JSON to book struct
	book := new(Book)
	err := json.Unmarshal([]byte(req.Body), book)
	if err != nil {
		return clientError(http.StatusUnprocessableEntity)
	}

	// Validate input data
	if !isbnRegexp.MatchString(book.ISBN) {
		return clientError(http.StatusBadRequest)
	}
	if book.Title == "" || book.Author == "" {
		return clientError(http.StatusBadRequest)
	}

	// putItem returns an error (normally will be nil)
	err = putItem(book)
	if err != nil {
		return serverError(err)
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 201,
		Headers:    map[string]string{"Location": fmt.Sprintf("/books?pk=%s&sk=%s", book.ISBN, book.Author)},
		Body:       req.Body,
	}, nil
}

// PATCH request must use Book fields as JSON
func update(req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// Use unmarshal to map JSON to book struct
	bookReq := new(Book)
	err := json.Unmarshal([]byte(req.Body), bookReq)
	if err != nil {
		return clientError(http.StatusUnprocessableEntity)
	}

	// Get the ISBN from the pk param and validate
	pk := bookReq.ISBN
	if !isbnRegexp.MatchString(pk) {
		return clientError(http.StatusBadRequest)
	}

	// Get the Author from the sk param
	sk := bookReq.Author
	if sk == "" {
		return clientError(http.StatusBadRequest)
	}

	// Get the book response from DynamoDB based on the pk,sk pair
	bookRes, err := incrementItem(pk, sk)
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
