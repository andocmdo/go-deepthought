package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	Header      string
	ContentType string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		"",
		"",
		Index,
	},
	Route{
		"IndexFile",
		"GET",
		"/{file}",
		"",
		"",
		Index,
	},
	Route{
		"JobIndex",
		"GET",
		"/api/v1/jobs",
		"",
		"",
		JobIndex,
	},
	Route{
		"JobShow",
		"GET",
		"/api/v1/jobs/{jobID}",
		"",
		"",
		JobShow,
	},
	Route{
		"JobCreateURLEnc",
		"POST",
		"/api/v1/jobs",
		"Content-Type",
		//"application/x-www-form-urlencoded",
		"application/x-www-form-urlencoded.*",
		JobCreateURLEnc,
	},
	Route{
		"JobCreateJSON",
		"POST",
		"/api/v1/jobs",
		"",
		"",
		JobCreateJSON,
	},
}
