# Let's Go Pro — Web Application

A production-style web application built while working through **"Let's Go Pro"** (Alex Edwards) in **Go**.

## What it does

- A "snippetbox"-style app: create, view and share text snippets
- Full user auth: signup, login, logout with **session management** (`scs`) and MySQL-backed sessions
- Secure middleware chain: panic recovery, request logging, common headers, session loading
- Server-side rendered HTML templates (base layout, partials, pages) with a custom `humanDate` template function
- Form decoding + validation via `go-playground/form` and a custom validator package
- TLS and sensible HTTP server timeouts configured in `main.go`

## Tech Stack

`Go` `net/http` `MySQL` `html/template` `scs` (sessions) `alice` (middleware chaining) `Docker` (dev DB)

## How to run

```bash
docker compose -f doc/docker/compose.yml up -d   # MySQL
go run ./cmd/web -addr=:4000
# open http://localhost:4000
```

SQL schema and seed data live in `doc/queries.sql`; API collection in `doc/bruno`.
