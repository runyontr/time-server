package main

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/sirupsen/logrus"
	"net/http"
	t "time"
)

//V1AppInfoPath contains the path for where the GetAppInfo call will listen
const V1GetTime = "/v1/time"

//MakeInfoServiceHandler creates an http.Handler object hosting endpoints for the
// functions defined in service.AppInfoService
func MakeInfoServiceHandler(srv TimeService, logger log.Logger) http.Handler {
	//uses Go-Kit for handling

	options := []httptransport.ServerOption{
		//How should errors get encoded
		httptransport.ServerErrorEncoder(errorEncoder),
		//The server will use the provided log.Logger object to log errors
		httptransport.ServerErrorLogger(logger),
	}

	m := http.NewServeMux()

	var timeEndpoint endpoint.Endpoint
	{
		timeEndpoint = MakeAppInfoEndpoint(srv)
		timeEndpoint = LoggingMiddleware(log.With(logger, "method", "GetTime"))(timeEndpoint)
	}

	m.Handle(V1GetTime, httptransport.NewServer(
		timeEndpoint,
		decodeEmptyRequest,
		encodeGenericResponse,
		options...,
	))
	return m
}

/****************
	Endpoints
 *************/

func MakeAppInfoEndpoint(srv TimeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return srv.GetTime()
	}
}

/**************
	Errors
 *************/

type errorWrapper struct {
	Error string `json:"error"`
}

//errorEncoder describes how to write errors back to the client.
func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(err2code(err))
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}

//err2code allows the customization of response codes based on the error from
// the service.  Currently no custom errors defined and everything is a 500
func err2code(err error) int {
	switch err {

	}
	return http.StatusInternalServerError
}

/*********
	Decoders
 ************/

// decodeEmptyRequest can be used when there is no data to extract from the http request.
// this function will log the request
func decodeEmptyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	logrus.WithFields(logrus.Fields{
		"Method":     r.Method,
		"Path":       r.URL.Path,
		"Host":       r.Host,
		"RemoteAddr": r.RemoteAddr,
		"RequestURI": r.RequestURI}).Infof("Received Message at %v", t.Now())
	return nil, nil
}

/***********
Encoders
*/

// encodeGenericResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer.
func encodeGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
