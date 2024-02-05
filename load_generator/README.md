# Load Generator

This is the source code for the load generator.

## Prerequisites

- [Go](https://go.dev/dl/) (v1.21.6 at the time of writing)

## Run locally

- Make sure you have SurrealDB running locally and the database is seeded. (See the [prepare_db](../prepare_db/README.md) directory for more information)
- Download the required packages

```bash
go mod download
```

- Build load generator

```bash
go build .
```

- Run the load generator

```bash
./load_generator
```
