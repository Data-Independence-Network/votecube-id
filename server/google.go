package server

import (
	"bytes"
	"database/sql"
	"fmt"
	"votecube-id/db"
	"votecube-id/models"
	"votecube-id/verify"

	"github.com/valyala/fasthttp"
)

const (
	Dev  = 0
	Prod = 1
)

const (
	InvalidUri = "0"
)

var GoogleLoginUri = []byte("/go/s")

type Environment = int

type MyHandler struct {
	foobar string
}

/**
 * Start a server on a given port
 */
func Start(port string, env Environment, dBase *sql.DB) {

	// myHandler := &MyHandler{foobar: "foobar"}

	connectString := ":" + port

	// fmt.Println("Start 2.0 " + connectString)
	// fasthttp.ListenAndServe(":8080", myHandler.HandleFastHTTP)

	var server *fasthttp.Server

	switch env {
	case Dev:
		server = &fasthttp.Server{
			Handler:               devHandler,
			MaxRequestBodySize:    2048,
			NoDefaultServerHeader: true,
			NoDefaultContentType:  true,
		}
	case Prod:
		server = &fasthttp.Server{
			// DisableKeepalive:      true,
			Handler:               prodHandler,
			MaxRequestBodySize:    2048,
			NoDefaultServerHeader: true,
			NoDefaultContentType:  true,
			// TCPKeepalive:          false,
		}
	}

	fmt.Println("Start 2.7")
	// pass plain function to fasthttp
	server.ListenAndServe(connectString)
}

// request handler in net/http style, i.e. method bound to MyHandler struct
func (h *MyHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hello, world! Requested path is %q. Foobar is %q",
		ctx.Path(), h.foobar)
}

func prodHandler(ctx *fasthttp.RequestCtx) {
	respond(ctx)
}

func devHandler(ctx *fasthttp.RequestCtx) {
	if ctx.Request.Header.IsOptions() {
		ctx.Response.Header.Set("Allow", "OPTIONS, PUT")
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "OPTIONS, PUT")
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "http://localhost:8000")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "content-type")
		return
	}

	if !ctx.Request.Header.IsPut() {
		return
	}

	ctx.Response.Header.Set("Access-Control-Allow-Origin", "http://localhost:8000")

	respond(ctx)
}

func respond(ctx *fasthttp.RequestCtx) {
	if !ctx.Request.Header.IsPut() {
		return
	}

	var requestType = ctx.Request.RequestURI()

	if bytes.Equal(GoogleLoginUri, requestType) {
		var claims, err = verify.VerifyToken(ctx.Request.Body())
		if err != nil {
			fmt.Fprintf(ctx, err.Error())
			return
		}
		u := &models.User{ID: 1, Email: "tester"}
		err = db.SaveUser(u)
		if err != nil {
			fmt.Fprintf(ctx, err.Error())
			return
		}
		fmt.Fprintf(ctx, claims.Email)
	} else {
		fmt.Fprintf(ctx, InvalidUri)
		return
	}

	fmt.Fprintf(ctx, "Token %q", ctx.Request.Body())
}
