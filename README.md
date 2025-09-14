# Simpleblog

Simple Personal Blog with Basic Authentication

## Concepts

- Requirements and [Basic
  Authentication](https://www.youtube.com/watch?v=mwccHwUn7Gc&t=20s) from
  [roadmap.sh](https://roadmap.sh/projects/personal-blog)
- Wikilinks from [gowiki](https://go.dev/doc/articles/wiki/)
- HTTP Router and URL Matcher with [mux](https://github.com/gorilla/mux)
- Slugify with [slug](https://github.com/gosimple/slug)
- Markdown Parser with [goldmark](https://github.com/yuin/goldmark)
- HTML Template [html/template](https://pkg.go.dev/html/template)
- CSV Parser with [standard library](https://pkg.go.dev/encoding/csv)
- Project template with [go-blueprint](https://github.com/Melkeydev/go-blueprint)
- Live reload with [air](https://github.com/air-verse/air)
- More

## Features

- CRUD article
- Parse `[[link]]` syntax to be wikilinks

## Getting Started

These instructions will get you a copy of the project up and running on your
local machine for development and testing purposes. See deployment for notes on
how to deploy the project on a live system.

## Todo

- [ ] Manage static files like images and videos to be used in our blog

## MakeFile

Run build make command with tests

```bash
make all
```

Build the application

```bash
make build
```

Run the application

```bash
make run
```

Live reload the application:

```bash
make watch
```

Run the test suite:

```bash
make test
```

Clean up binary from the last build:

```bash
make clean
```
