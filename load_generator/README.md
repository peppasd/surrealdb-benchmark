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

- Check the available options

```bash
./load_generator -help
```

Here is an example on how to run the load generator locally:

```bash
./load_generator -minutes 20 -threads 3 -url localhost:8000
```
