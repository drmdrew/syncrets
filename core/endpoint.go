package core

import (
	"net/url"
)

/*
// Endpoint describes the attributes of a target backend endpoint
type Endpoint struct {
	Name      string
	RawURL    *url.URL
	ServerURL *url.URL
	Path      string
	Walker    Walker
}
*/
// Endpoint describes the attributes of a target backend endpoint
type Endpoint interface {
	GetName() string
	GetRawURL() *url.URL
	GetURL() *url.URL
	GetPath() string
	Walk(visitor Visitor)
	Write(secret Secret) error
}
