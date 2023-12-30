# Buddy

Buddy is a Go project that provides a simple randomizer service. It allows you to create a list of choices and then randomly select one of the choices based on their weights.

## Project Structure

The project is structured into several packages:

- `app`: Contains the application logic, including the randomizer and the local repository for storing choice lists.
- `backend`: Contains the main entry point for the backend service.
- `core`: Contains the core business logic and interfaces for the randomizer and the choice list repository.
- `framework`: Contains the gRPC service definitions and generated code.
- `frontend`: Contains the main entry point for the frontend service.

## How to Run

To run the backend service, navigate to the `backend` directory and run:

```sh
go run main.go
```

To run the frontend service, navigate to the `frontend` directory and run:
```sh
go run main.go
```

## Dependencies
The project uses the following dependencies:

- `github.com/google/uuid` for generating unique identifiers.
- `google.golang.org/grpc` for the gRPC service.
- `google.golang.org/protobuf` for protobuf support.

## License
This project is licensed under the MIT License. See the `LICENSE` file for details.
