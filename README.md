# Buddy
[![Go](https://github.com/tthung1997/buddy/actions/workflows/go.yml/badge.svg)](https://github.com/tthung1997/buddy/actions/workflows/go.yml)

Buddy is a Go project that represents a friend who can help you with your day-to-day simple needs such as: 

- Create a list of choices and then randomly select one of the choices based on their weights.
- Construct a ranking for a list of items

## Project Structure

The project is structured into several packages:

- `app`: Contains the application logic.
- `backend`: Contains the main entry point for the backend service.
- `core`: Contains the core business logic and interfaces.
- `framework`: Contains the gRPC service definitions and generated code.
- `frontend`: Contains the main entry point for the frontend service.

## How to Run

To run the backend service, navigate to the `backend` directory and run `go run` for the microservice that you need:

```sh
go run choice.go
```

To run the frontend service, navigate to the `frontend` directory and run `go run` for the microservice that you need:

```sh
go run choice.go
```

## Dependencies
The project uses the following dependencies:

- `google.golang.org/grpc` for the gRPC service.
- `google.golang.org/protobuf` for protobuf support.

## Tests
Unit tests are located in the tests directory. To run the tests:

```sh
go test ./tests/...
```

## License
This project is licensed under the MIT License. See the `LICENSE` file for details.
