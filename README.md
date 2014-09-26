chatter
=======

Chatter is an SSE demo application written in Go.


## Getting Started

To install chatter, simply install Go and run:

```sh
$ go get github.com/benbjohnson/chatter
```

Then run `chatterd`:

```sh
$ chatterd
Listening on http://localhost:9000
```

You can connect a client to chatter by simply going to
[http://localhost:9000](http://localhost:9000). There is not currently a message
creation UI so you'll need to use `curl` or `wget`:

```sh
$ curl -X POST http://localhost:9000/messages?body=hello
```
