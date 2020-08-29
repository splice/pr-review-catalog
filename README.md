## Welcome to the Catalog Interview!

## Package structure

Everything you'll need will be located in:

- model/
  - Where model types/fields are defined
- server/
  - controller.go
    - Where the routes are defined
  - error.go
    - Where error types are defined
  - handler.go
    - Define new route handlers are defined
  - helper.go
    - Where handler helpers for things like parsing params are defined
  - payload.go
    - Where the responses from the handlers are defined
- storage/
  - Where the repositories for finding data from storage are defined
  - seeds/
    - Where you can set up new seed data so that you don't need to setup a local database

## Running the service

Run:

`$ go run main.go`

Navigate in your browser to: `localhost`

This will show a healthcheck to make sure everything is running correctly.

From there, you can navigate to other routes like:
`localhost/products` or `localhost/products?limit=2` or any routes that you define

Optionally, you can start the server with a specified port:

`$ go run main.go -port=8080`

Then navigate to `localhost:8080`

## Current functionality
There is only one route `/products` that fetches all products in storage up to the provided or default max limit.
Products are currently stored in `./storage/seeds/products.json` in lieu of having to set up a database.
