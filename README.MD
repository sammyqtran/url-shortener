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
