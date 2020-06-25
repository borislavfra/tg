// GENERATED BY 'T'ransport 'G'enerator. DO NOT EDIT.
package transport

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/savsgio/gotils"
	"github.com/valyala/fasthttp"
)

func (http *httpJsonRPC) serveTest(ctx *fasthttp.RequestCtx) {
	http.serveMethod(ctx, "test", http.test)
}
func (http *httpJsonRPC) test(span opentracing.Span, ctx *fasthttp.RequestCtx, requestBase baseJsonRPC) (responseBase *baseJsonRPC) {

	var err error
	var request requestJsonRPCTest

	if err = json.Unmarshal(requestBase.Params, &request); err != nil {
		ext.Error.Set(span, true)
		span.SetTag("msg", "request body could not be decoded: "+err.Error())
		return makeErrorResponseJsonRPC(requestBase.ID, parseError, "request body could not be decoded: "+err.Error(), nil)
	}

	if requestBase.Version != Version {
		ext.Error.Set(span, true)
		span.SetTag("msg", "incorrect protocol version: "+requestBase.Version)
		return makeErrorResponseJsonRPC(requestBase.ID, parseError, "incorrect protocol version: "+requestBase.Version, nil)
	}

	methodContext := opentracing.ContextWithSpan(ctx, span)

	var response responseJsonRPCTest

	response.Ret1, response.Ret2, err = http.svc.Test(methodContext, request.Arg0, request.Arg1, request.Opts...)

	if err != nil {
		if http.errorHandler != nil {
			err = http.errorHandler(err)
		}
		ext.Error.Set(span, true)
		span.SetTag("msg", err.Error())
		span.SetTag("errData", toString(err))
		return makeErrorResponseJsonRPC(requestBase.ID, internalError, err.Error(), err)
	}

	responseBase = &baseJsonRPC{
		ID:      requestBase.ID,
		Version: Version,
	}

	if responseBase.Result, err = json.Marshal(response); err != nil {
		ext.Error.Set(span, true)
		span.SetTag("msg", "response body could not be encoded: "+err.Error())
		return makeErrorResponseJsonRPC(requestBase.ID, parseError, "response body could not be encoded: "+err.Error(), nil)
	}
	return
}

func (http *httpJsonRPC) serveBatch(ctx *fasthttp.RequestCtx) {

	batchSpan := extractSpan(http.log, fmt.Sprintf("jsonRPC:%s", gotils.B2S(ctx.URI().Path())), ctx)
	defer injectSpan(http.log, batchSpan, ctx)
	defer batchSpan.Finish()
	methodHTTP := gotils.B2S(ctx.Method())

	if methodHTTP != fasthttp.MethodPost {
		ext.Error.Set(batchSpan, true)
		batchSpan.SetTag("msg", "only POST method supported")
		ctx.Error("only POST method supported", fasthttp.StatusMethodNotAllowed)
		return
	}

	for _, handler := range http.httpBefore {
		handler(ctx)
	}

	if value := ctx.Value(CtxCancelRequest); value != nil {
		return
	}

	var err error
	var requests []baseJsonRPC

	if err = json.Unmarshal(ctx.PostBody(), &requests); err != nil {
		ext.Error.Set(batchSpan, true)
		batchSpan.SetTag("msg", "request body could not be decoded: "+err.Error())

		for _, handler := range http.httpAfter {
			handler(ctx)
		}
		sendResponse(http.log, ctx, makeErrorResponseJsonRPC([]byte("\"0\""), parseError, "request body could not be decoded: "+err.Error(), nil))
		return
	}

	responses := make(jsonrpcResponses, 0, len(requests))

	var wg sync.WaitGroup

	for _, request := range requests {

		methodNameOrigin := request.Method
		method := strings.ToLower(request.Method)

		span := opentracing.StartSpan(request.Method, opentracing.ChildOf(batchSpan.Context()))
		span.SetTag("batch", true)

		switch method {

		case "test":

			wg.Add(1)
			go func(request baseJsonRPC) {
				responses.append(http.test(span, ctx, request))
				wg.Done()
			}(request)

		default:
			ext.Error.Set(span, true)
			span.SetTag("msg", "invalid method '"+methodNameOrigin+"'")
			responses.append(makeErrorResponseJsonRPC(request.ID, methodNotFoundError, "invalid method '"+methodNameOrigin+"'", nil))
		}
		span.Finish()
	}
	wg.Wait()

	for _, handler := range http.httpAfter {
		handler(ctx)
	}
	sendResponse(http.log, ctx, responses)
}

func (http *httpJsonRPC) serveMethod(ctx *fasthttp.RequestCtx, methodName string, methodHandler methodJsonRPC) {

	span := extractSpan(http.log, fmt.Sprintf("jsonRPC:%s", gotils.B2S(ctx.URI().Path())), ctx)
	defer injectSpan(http.log, span, ctx)
	defer span.Finish()

	methodHTTP := gotils.B2S(ctx.Method())

	if methodHTTP != fasthttp.MethodPost {
		ext.Error.Set(span, true)
		span.SetTag("msg", "only POST method supported")
		ctx.Error("only POST method supported", fasthttp.StatusMethodNotAllowed)
	}

	for _, handler := range http.httpBefore {
		handler(ctx)
	}

	if value := ctx.Value(CtxCancelRequest); value != nil {
		ext.Error.Set(span, true)
		span.SetTag("msg", "request canceled")
		return
	}

	var err error
	var request baseJsonRPC
	var response *baseJsonRPC

	if err = json.Unmarshal(ctx.PostBody(), &request); err != nil {
		ext.Error.Set(span, true)
		span.SetTag("msg", "request body could not be decoded: "+err.Error())

		for _, handler := range http.httpAfter {
			handler(ctx)
		}
		sendResponse(http.log, ctx, makeErrorResponseJsonRPC([]byte("\"0\""), parseError, "request body could not be decoded: "+err.Error(), nil))
		return
	}

	methodNameOrigin := request.Method
	method := strings.ToLower(request.Method)

	if method != "" && method != methodName {
		ext.Error.Set(span, true)
		span.SetTag("msg", "invalid method "+methodNameOrigin)

		for _, handler := range http.httpAfter {
			handler(ctx)
		}
		sendResponse(http.log, ctx, makeErrorResponseJsonRPC(request.ID, methodNotFoundError, "invalid method "+methodNameOrigin, nil))
		return
	}

	response = methodHandler(span, ctx, request)

	for _, handler := range http.httpAfter {
		handler(ctx)
	}

	if response != nil {
		sendResponse(http.log, ctx, response)
	}
}
