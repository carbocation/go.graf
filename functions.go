package graf

import (
	"net/url"
)

//We require all URLs to have a scheme (e.g., http) and a host (e.g., www.google.com)
func ValidUrl(u *url.URL) bool {
	if u.Scheme == "" || u.Host == "" {
		return false
	} else {
		return true
	}
}
