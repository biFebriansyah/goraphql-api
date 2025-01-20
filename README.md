# goraphql-api

A backend service for managing simple data using graphql. Designed and developed for research and development purposes, with a focus on scalability, efficiency, and modern communication protocols.

## Features

-   Integration with Grahpql for client communication.
-   CRUD operations for resource.
-   Secure and optimized for production.

## Tech Stack

-   **Language:** Golang
-   **Framework:** Go-fiber, Graphql
-   **Database:** Mongodb
-   **Other Tools:** Docker, Make

## Getting Started

### Prerequisites

-   Go (v1.24+)
-   Docker (optional)
-   Mongodb instance (local or cloud)

### Project setup

```bash
# Clone the repository
$ git clone https://github.com/biFebriansyah/goraphql-api.git

# Install Package
$ go mod download

# using Make
$ make install
```

### Compile and run the project

```bash
# development
$ go run *.go

# generate
$ make generate

# watch mode
$ make run

# production mode
$ make build
```

## Authors

-   [@biFebriansyah](https://www.github.com/biFebriansyah)

## License

[MIT](https://choosealicense.com/licenses/mit/)
