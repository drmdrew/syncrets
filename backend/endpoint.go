package backend

import "net/url"

// Endpoint describes the attributes of a target backend endpoint
type Endpoint struct {
	Name      string
	RawURL    *url.URL
	ServerURL *url.URL
	Path      string
}
