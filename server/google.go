package server

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

type MyHandler struct {
	foobar string
}

// request handler in net/http style, i.e. method bound to MyHandler struct
func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hello, world! Requested path is %q. Foobar is %q",
		ctx.Path(), h.foobar)
}

// func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
// 	fmt.Fprintf(ctx, "Hi there! Request URI is %q", ctx.RequestURI())
// }

/**
 * Start a server on a given port
 */
func Start(port string) {

	myHandler := &MyHandler{foobar: "foobar"}

	fasthttp.ListenAndServe(":"+port, myHandler.HandleFastHTTP)

	// pass plain function to fasthttp
	// fasthttp.ListenAndServe(":8081", fastHTTPHandler)
}
