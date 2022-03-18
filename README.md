# Golang AWS Lambda/DynamoDB/API Gateway Project

An adaptation of the [Golang API/Lambda Tutorial written by Alex Edwards](https://www.alexedwards.net/blog/serverless-api-with-go-and-aws-lambda)

## Changes

You can follow the link to Alex's blog article and read through the original project, but I made several changes to the code to accommodate my goals. Most of this directly relates to using the HTTP Gateway, aka APIGATEWAYV2. You can find the relevant structs in the [AWS Github Repo.](https://github.com/aws/aws-lambda-go/blob/0462b0000e7468bdc8a9c456273c1551fab284aa/events/apigw.go#L123)

To save some guesswork on your part of what I actually changed, I updated the ProxyRequest and ProxyResponse to instead use APIGatewayV2HTTPRequest and APIGatewayV2HTTPResponse. Alex's blog uses a query param for the GET /books route, but I was more interested in the pure API functionality so I changed the usage of accessing the QueryStringParameters from APIGW Request to instead process the query as JSON input.

Another notable breaking change between the REST Gateway (v1) and the HTTP Gateway (v2) is that the HTTP Method used in the request is not stored in HTTPMethod, instead those can be accessed from RequestContext.HTTP.Method.

Alex's blog only writes handlers for the GET and POST actions; I added a third action for PATCH. The PATCH handler in this case simply increments the itemCount by 1. Speaking of itemCount, fun fact, don't create a DynamoDB field called merely "count", as apparently that is a protected keyword. I ran into some inconsistent behavior that I finally tracked down by using the awscli tool to determine that my UpdateExpression was trying to modify fields that were named the same as the protected keyword.

Other changes I made were to rewrite the DynamoDB session initialization by following the pattern in the [AWS SDK.](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/using-dynamodb-with-go-sdk.html)

Finally, I refactored out some of the helper functions Alex wrote in main.go to instead be stored in utils.go, which I find handy since it keeps main.go focused being a method handler. The validation functions that appear in the handlers I also moved here to keep the handlers themselves to look neat and tidy.
