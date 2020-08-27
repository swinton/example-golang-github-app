# `example-golang-github-app`
> An example GitHub App, written in Golang

## Usage

```shell
make

./bin/example-golang-github-app
```

## Development

Use [`reflex`](https://github.com/cespare/reflex) to auto-reload the webhook server when any file changes are detected:

```shell
reflex -s go run cmd/example-golang-github-app/main.go
```

Send an example webhook:

```shell
curl -X POST -d '{"hello":"world"}' -H "Content-type: application/json" http://localhost:8000/
```

## Resources

- [HTTP Server - Go Web Examples](https://gowebexamples.com/http-server/)
- [Routing (using gorilla/mux) - Go Web Examples](https://gowebexamples.com/routes-using-gorilla-mux/)
- [Accepting Github Webhooks with Go · groob.io](https://groob.io/tutorial/go-github-webhook/)
