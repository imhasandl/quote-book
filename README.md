## Quote Book

## Prerequisites

- Go 1.20 or later
- PostgreSQL database

## Setup project

Clone the remote repository

```bash
git clone https://github.com/imhasandl/quote-book
cd quote-book
```

Install all dependency

```bash
go mod tidy
```

And run the server

```bash
go run main.go
```

## Setup Postgresql database

Install Postgresql in Linux / WSL environment.

```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
```

You can check if it's working using this command.

```bash
psql --version
```

(Linux only) Update postgres password:

```bash
sudo passwd postgres
```

Start Postgresql server

```bash
sudo service postgresql start
```

Enter the psql shell:

```bash
sudo -u postgres psql
```

Then you will be in postgres shell, you need to create database to configure it in .env file. After creation of database you can go inside it using \c and then the name of database.

```sql
CREATE DATABASE name_of_your_database;
```

## Configuration

Create a `.env` file in the root directory with the following variables:

```env
PORT="YOUR_PORT"
DB_URL="postgres://username:password@host:port/name_of_your_database?sslmode=disable"
```

This service uses Goose for database migrations, before setting up migration cd into directory that is stored and then run the script:

```bash
# Install Goose
go install github.com/pressly/goose/v3/cmd/goose@latest

# Go inside database/migrations where the files located
cd database/migrations

# Run migrations
goose postgres "YOUR_DB_CONNECTION_STRING" up
```


## HTTP methods

---
### CreateQuote - /quotes

Stores implemented request in database.

```json
{
  "Text": "Some quote",
  "Author": "author of the qoute"
}
```
---

### GetQuotes - /quotes or /quotes?author=NAME OF THE AUTHOR

Gets all quotes if the filter is not implemented. if request sent in /quotes it will return all quotes from database. If use some filter like /quotes?author=AUTHOR it will return the quotes from that specific author.

---

### RandomQuote - /quotes/random

Returns random quote from database

---

### DeleteQuote - /quotes/{id}

Deletes quote who has this id.

---

## Run tests for http handlers

```bash
go test ./...
```

## Test http request using curl commands

## Create a quote

```bash
curl -X POST http://localhost:8080/quotes \
  -H "Content-Type: application/json" \
  -d '{"author":"Confucius", "quote":"Life is simple, but we insist on making it complicated."}'
```

**Example response:**
```json
{
  "id": 1,
  "author": "Confucius",
  "quote": "Life is simple, but we insist on making it complicated."
}
```

---

### Get all quotes

```bash
curl http://localhost:8080/quotes
```

---

### Get a random quote

```bash
curl http://localhost:8080/quotes/random
```

---

### Get quotes by author

```bash
curl http://localhost:8080/quotes?author=Confucius
```

---

### Delete a quote by ID

```bash
curl -X DELETE http://localhost:8080/quotes/1
```

**Example response:**
```json
{
  "status": true
}
```










