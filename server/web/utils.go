package web

import "net/http"

func GetCookieByNameForResponse(r *http.Response, name string) *http.Cookie {
	var cookie *http.Cookie
	cookies := r.Cookies()
	// Search the cookie containing the token
	for _, tmpCookie := range cookies {
		if tmpCookie.Name == name && tmpCookie.Value != "" {
			cookie = tmpCookie
		}
	}
	return cookie
}

func GetCookieByNameForRequest(r *http.Request, name string) *http.Cookie {
	var cookie *http.Cookie
	cookies := r.Cookies()
	// Search the cookie containing the token
	for _, tmpCookie := range cookies {
		if tmpCookie.Name == name && tmpCookie.Value != "" {
			cookie = tmpCookie
		}
	}
	return cookie
}
