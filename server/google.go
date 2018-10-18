package server

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

const (
	Dev  = 0
	Prod = 1
)

type Environment = int

type MyHandler struct {
	foobar string
}

/**
 * Start a server on a given port
 */
func Start(port string, env Environment) {

	// myHandler := &MyHandler{foobar: "foobar"}

	connectString := ":" + port

	// fmt.Println("Start 2.0 " + connectString)
	// fasthttp.ListenAndServe(":8080", myHandler.HandleFastHTTP)

	var httpHandler fasthttp.RequestHandler

	switch env {
	case Dev:
		httpHandler = devHandler
	case Prod:
		httpHandler = prodHandler
	}

	fmt.Println("Start 2.2")
	// pass plain function to fasthttp
	fasthttp.ListenAndServe(connectString, httpHandler)

	fmt.Println("Start 3")
}

// request handler in net/http style, i.e. method bound to MyHandler struct
func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hello, world! Requested path is %q. Foobar is %q",
		ctx.Path(), h.foobar)
}

func prodHandler(ctx *fasthttp.RequestCtx) {
	if ctx.IsOptions() {
		return
	}

	fmt.Fprintf(ctx, "Hi there! Request URI is %q", ctx.RequestURI())
}

func devHandler(ctx *fasthttp.RequestCtx) {
	if ctx.IsOptions() {
		ctx.Response.Header.Set("Allow", "POST")
	}
	fmt.Fprintf(ctx, "Hi there! Request URI is %q", ctx.RequestURI())
}
