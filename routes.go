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
		"FrontEnd",
		"GET",
		"/frontend/{file}",
		"",
		"",
		FrontEnd,
	},
	Route{
		"Index",
		"GET",
		"/",
		"",
		"",
		Index,
	},
	Route{
		"JobIndex",
		"GET",
		"/jobs",
		"",
		"",
		JobIndex,
	},
	Route{
		"JobShow",
		"GET",
		"/jobs/{jobID}",
		"",
		"",
		JobShow,
	},
	Route{
		"JobCreateURLEnc",
		"POST",
		"/jobs",
		"Content-Type",
		//"application/x-www-form-urlencoded",
		"application/x-www-form-urlencoded.*",
		JobCreateURLEnc,
	},
	Route{
		"JobCreateJSON",
		"POST",
		"/jobs",
		"",
		"",
		JobCreateJSON,
	},
	Route{
		"JobDestroy",
		"DELETE",
		"/jobs/{jobID}",
		"",
		"",
		JobDestroy,
	},
}
