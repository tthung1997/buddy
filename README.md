# Buddy
[![Go](https://github.com/tthung1997/buddy/actions/workflows/go.yml/badge.svg)](https://github.com/tthung1997/buddy/actions/workflows/go.yml)
[![CodeQL](https://github.com/tthung1997/buddy/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/tthung1997/buddy/actions/workflows/github-code-scanning/codeql)
[![pages-build-deployment](https://github.com/tthung1997/buddy/actions/workflows/pages/pages-build-deployment/badge.svg)](https://github.com/tthung1997/buddy/actions/workflows/pages/pages-build-deployment)

Buddy is a Go project that represents a friend who can help you with your day-to-day simple needs such as: 

- Create a list of choices and then randomly select one of the choices based on their weights.
- Construct a ranking for a list of items
- Keep a recipe library and a pantry, then answer what can be cooked right now

## Project Structure

The project is structured into several packages:

- `app`: Contains the application logic.
- `backend`: Contains the main entry point for the backend service.
- `core`: Contains the core business logic and interfaces.
- `framework`: Contains the gRPC service definitions and generated code.
- `frontend`: Contains the main entry point for the frontend service.

## Modules

The frontend serves the following modules:

- **Board Games** (`/boardgames`): browse a BoardGameGeek collection, pick something to play, and track backed or preordered games.
- **Cooking** (`/cooking`): keep recipes with their ordered steps and timings, track pantry amounts, and let Buddy join the two to answer what can be cooked right now. Missing ingredients can be pushed straight to the shopping list.- **Shopping** (`/shopping`): keep the household inventory current and queue the next store run.

### Cooking

The cooking module is built from three layers:

- `core/cooking`: units and their kinds, durations, ingredients, recipes with steps, pantry items, and the `ICookEngine` interface.
- `app/cooking`: `SimpleUnitConverter` (conversion within a unit kind), `PantryCookEngine` (pantry to recipe matching, aggregated shopping lists, and pantry deduction), and JSON backed repositories.
- `frontend/cooking`: the dashboard, recipe library and editor, pantry workspace, and the "what can I cook" view.

Recipes, ingredients, and the pantry are stored as JSON under `frontend/cooking/.db/`, which is gitignored and created on first run.

#### Staples

Some ingredients are not worth measuring. Salt, oil, and rice are either in the cupboard or they are not, so an ingredient can be marked a **staple** and is then tracked by presence instead of by amount:

| | Measured ingredient | Staple |
| --- | --- | --- |
| Matching | compares amounts and converts units | satisfied while in stock or running low |
| Cooking | deducts the amount used | left untouched |
| Shopping list | `Flour (350 g)` | `Salt`, with no quantity |

A staple carries one of three levels: **in stock**, **running low**, or **out of stock**. Running low still cooks tonight but reaches the shopping list, while out of stock blocks any recipe that needs it. Recipes keep their amounts either way, so a recipe still tells you to add 5 g of salt.

## How to Run

To run the backend service, navigate to the `backend` directory and run `go run` for the microservice that you need:

```sh
go run choice.go
```

To run the frontend service:

```sh
go run frontend/main.go
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
