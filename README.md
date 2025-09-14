# Simpleblog

Simple Personal Blog with Basic Authentication

## Concepts

- Requirements and [Basic
  Authentication](https://www.youtube.com/watch?v=mwccHwUn7Gc&t=20s) from
  [roadmap.sh](https://roadmap.sh/projects/personal-blog)
- HTTP Router and URL Matcher with [mux](https://github.com/gorilla/mux)
- Slugify with [slug](https://github.com/gosimple/slug)
- Markdown Parser with [goldmark](https://github.com/yuin/goldmark)
- HTML Template [html/template](https://pkg.go.dev/html/template)
- Project template with [go-blueprint](https://github.com/Melkeydev/go-blueprint)
- Live reload with [air](https://github.com/air-verse/air)
- Etc

## Features

- CRUD articles
- Basic Authentication
- Wiki links style

## Todo

- [ ] Manage static files like images and videos to be used in our blog
- [ ] Filter, search, sort and paginate articles

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
