# books-api-golang
An adaptation of the [Golang API/Lambda Tutorial written by Alex Edwards](https://www.alexedwards.net/blog/serverless-api-with-go-and-aws-lambda)

## Changes
You can follow the link to Alex's blog article, but I made several changes to the code to accommodate my goals. Most of this directly relates to using the HTTP Gateway, aka APIGATEWAYV2. You can find the relevant structs in the [AWS Github Repo.](https://github.com/aws/aws-lambda-go/blob/0462b0000e7468bdc8a9c456273c1551fab284aa/events/apigw.go#L123)

To save some guesswork on your part, I updated the ProxyRequest and ProxyResponse to instead use APIGatewayV2HTTPRequest and APIGatewayV2HTTPResponse. Alex's blog uses a query param for the GET /books route, but I am more interested in the pure API functionality so I changed the usage of accessing the QueryStringParameters from APIGW Request to instead process the query as JSON input.

Another notable breaking change between the REST Gateway (v1) and the HTTP Gateway (v2) is that the HTTP Method used in the request is not stored in HTTPMethod, instead those can be accessed from RequestContext.HTTP.Method.

Finally, I refactored out some of the helper functions Alex wrote in main.go to instead be stored in utils.go, which I find handy since it keeps main.go focused being a method handler.
