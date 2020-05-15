package main

// Application serves as the base for all the API endpoints
// exposed by the web server. Some of the associated methods
// write the content of a dynamic template and others are used
// to process a HTTP request and return either a valid JSON
// object or a HTTP status code.
type Application struct{}
