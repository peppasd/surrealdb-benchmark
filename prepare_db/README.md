# Prepare the Database

This directory contains the script that creates the database schema and fake data from Faker.js. The resulting file is used to seed the database.

To use the same data used for the benchmark, you can use the `surrealdata.tar.gz` archive.

## Prerequisites

- [Bun](https://bun.sh/docs/installation) (v1.0.25 at the time of writing)

## Install dependencies

```bash
bun install
```

## Run the script

```bash
bun run index.ts --output db.surql
```

The flag `--output` specifies the path and name of the output file.

## Run the database locally

- Install [SurrealDB](https://docs.surrealdb.com/docs/installation/overview) (v1.1.1 at the time of writing)
- Start the database

```bash
surreal start --allow-net file://surreal.db
```

- Generate the required seed data using the script in this directory
- Import the seed data into the database

```bash
surreal import --conn http://localhost:8000 --ns benchmark --db benchmark db.surql
```
