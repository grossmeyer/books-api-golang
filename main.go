package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Book struct {
	ISBN   string `json:"pk"`
	Author string `json:"sk"`
	Title  string `json:"title"`
	Count  int    `json:"count"`
}

func router(req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	switch req.RequestContext.HTTP.Method {
	case "GET":
		return show(req)
	case "POST":
		return create(req)
	default:
		return clientError(http.StatusMethodNotAllowed)
	}
}

// GET request must use pk,sk as JSON
func show(req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// Convert JSON Body to format we can use
	var body map[string]interface{}
	json.Unmarshal([]byte(req.Body), &body)

	// Get the ISBN from the pk param and validate
	// Note: type assertion is mandatory for pk,sk
	pk := body["pk"].(string)
	if !isbnRegexp.MatchString(pk) {
		return clientError(http.StatusBadRequest)
	}

	// Get the Author from the sk param
	sk := body["sk"].(string)
	if sk == "" {
		return clientError(http.StatusBadRequest)
	}

	// Get the book from DynamoDB based on the pk,sk pair
	book, err := getItem(pk, sk)
	if err != nil {
		return serverError(err)
	}
	if book == nil {
		return clientError(http.StatusNotFound)
	}

	// APIGateway Body needs to be a string, so we convert here
	js, err := json.Marshal(book)
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
	if req.Headers["content-type"] != "application/json" && req.Headers["Content-Type"] != "application/json" {
		return clientError(http.StatusNotAcceptable)
	}

	book := new(Book)
	err := json.Unmarshal([]byte(req.Body), book)
	if err != nil {
		return clientError(http.StatusUnprocessableEntity)
	}

	if !isbnRegexp.MatchString(book.ISBN) {
		return clientError(http.StatusBadRequest)
	}
	if book.Title == "" || book.Author == "" {
		return clientError(http.StatusBadRequest)
	}

	err = putItem(book)
	if err != nil {
		return serverError(err)
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 201,
		Headers:    map[string]string{"Location": fmt.Sprintf("/books?pk=%s&sk=%s", book.ISBN, book.Author)},
	}, nil
}

func main() {
	lambda.Start(router)
}
