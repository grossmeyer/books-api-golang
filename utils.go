package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
)

var isbnRegexp = regexp.MustCompile(`[0-9]{3}\-[0-9]{10}`)
var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

// Helper for Unmarshaling JSON to Book struct
func unmarshalBookJson(req events.APIGatewayV2HTTPRequest) (*Book, error) {
	bookReq := new(Book)
	err := json.Unmarshal([]byte(req.Body), bookReq)
	return bookReq, err
}

// Helpers for error handling; logs to os.Stderr
func serverError(err error) (events.APIGatewayV2HTTPResponse, error) {
	errorLogger.Println(err.Error())

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

func clientError(status int) (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

// Returns true if formatted correctly
func validateJsonFormat(req events.APIGatewayV2HTTPRequest) bool {
	if req.Headers["content-type"] != "application/json" && req.Headers["Content-Type"] != "application/json" {
		return false
	}
	return true
}

// Returns true if ISBN is formatted correctly
func validateIsbnFormat(isbn string) bool {
	return isbnRegexp.MatchString(isbn)
}

// Returns true is field is not blank
func validateFieldLength(field string) bool {
	if field == "" {
		return false
	}

	return true
}

// Returns false if any validation fails
func validateReadRequest(isbn, author string) bool {
	if !validateIsbnFormat(isbn) {
		return false
	}

	if !validateFieldLength(author) {
		return false
	}

	return true
}

// Returns false if any validation fails
func validateWriteRequest(isbn, author, title string) bool {
	if !validateReadRequest(isbn, author) {
		return false
	}

	if !validateFieldLength(title) {
		return false
	}

	return true
}
