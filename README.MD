# A URL Shortener

## Quick Summary 

Takes a URL and transforms it to a shortcode that can be used to retrieve the original URL effectively shortening the original URL.

## Endpoints:

/ping: Health check  
Example: `localhost:8080/ping` -> returns if server online

/get/{code}: redirects user to the URL for the specified code  
Example: `localhost:8080/get/abc123` -> redirects to https://example.com/a/f8f9sio

/post: accepts JSON `{"url": "https://example.com/a/f8f9sio"}`  
Example: `localhost:8080/post` with body `{"url": "https://example.com/a/f8f9sio"}` -> returns `{"short_code":"abc123"}`


## Testing

Tests can be found in shorten_test.go

Unit tests cover the handler logic and helper methods.

To run test with coverage 

`
go test -cover ./...
`

To see an HTML report

`
go tool cover -html=coverage.out -o coverage.html
`

For full validation server needs to be hosted and endpoints tested manually end to end.

### Example Usage (Minimal)

Build binary file

`
go build -o urlshortener
`

Run application. This should host this on your local machine @ localhost:8080 or something similar.

`
./urlshorterner
`

Shorten a URL. This is an example curl to send a json to the server with an example URL.

`
curl -X POST http://localhost:8080/post \
     -H "Content-Type: application/json" \
     -d '{"url": "https://example.com/a/f8f9sio"}'
`

You will recieve a response with a shortcode.

`
{"short_code":"atkyt8"}
`

For testing purposes let's say this is hosted at localhost:8080. To see the redirect visit localhost:8080/atkyt8.

You should be redirected to https://example.com/a/f8f9sio



## Goals and Roadmap

Sprint 1: Basic URL shortening service with core API endpoints and unit tests. (Complete)

Sprint 2: Add persistent storage (e.g., database) instead of in-memory maps.

Sprint 3: Add analytics (click tracking, usage stats).

Sprint 4: Implement user auth and management.

Further goals include deployment to a cloud environment with CI/CD pipeline