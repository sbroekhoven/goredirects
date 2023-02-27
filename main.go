package goredirects

import (
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/idna"

	"github.com/sbroekhoven/goresolve"
)

// Data struct
type Data struct {
	URL          string       `json:"url,omitempty"`
	Redirects    []*Redirects `json:"redirects,omitempty"`
	Error        bool         `json:"error,omitempty"`
	ErrorMessage string       `json:"errormessage,omitempty"`
}

// Redirects struct
type Redirects struct {
	Number     int             `json:"number"`
	StatusCode int             `json:"statuscode,omitempty"`
	URL        string          `json:"url,omitempty"`
	Protocol   string          `json:"protocol,omitempty"`
	DNS        *goresolve.Data `json:"dns,omitempty"`
}

// Get function
func Get(redirecturl string, nameserver string) *Data {
	r := new(Data)

	r.URL = redirecturl
	u, err := url.Parse(redirecturl)
	if err != nil {
		r.Error = true
		r.ErrorMessage = err.Error()
		return r
	}

	fqdn := u.Hostname()

	// Valid server name (ASCII or IDN)
	fqdn, err = idna.ToASCII(fqdn)
	if err != nil {
		r.Error = true
		r.ErrorMessage = err.Error()
		return r
	}

	var i int

	// max 20 times
	for i < 20 {
		// set client to CheckRedirect, not following the redirect
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}}

		// redirecturl prefix check for incomplete
		if caseInsenstiveContains(redirecturl, "http://") == false && caseInsenstiveContains(redirecturl, "https://") == false {
			// TODO: Set warning
			redirecturl = "http://" + redirecturl
		}

		// Repair the request
		req, err := http.NewRequest("GET", redirecturl, nil)
		if err != nil {
			r.Error = true
			r.ErrorMessage = err.Error()
			return r
		}

		// Set User-Agent
		req.Header.Set("User-Agent", "Golang_Research_Bot/3.0")

		// Do the request.
		resp, err := client.Do(req)
		if err != nil {
			r.Error = true
			r.ErrorMessage = err.Error()
			return r
		}
		defer resp.Body.Close()

		// Set soms vars.
		redirect := new(Redirects)
		redirect.Number = i
		redirect.StatusCode = resp.StatusCode
		redirect.URL = resp.Request.URL.String()
		redirect.Protocol = resp.Proto

		// Only unique hosts in hostlist
		hostname, _ := hostFromURL(resp.Request.URL.String())
		// TODO: Some error handeling

		// redirect.DNS = getHosts(hostname, nameserver)
		redirect.DNS = goresolve.Hostname(hostname, nameserver)
		r.Redirects = append(r.Redirects, redirect)

		if resp.StatusCode == 200 || resp.StatusCode > 303 {
			break
		} else {
			redirecturl = resp.Header.Get("Location")
			i++
		}
	}

	return r
}

func hostFromURL(geturl string) (string, error) {
	u, err := url.Parse(geturl)
	if err != nil {
		return "", err
	}
	return u.Hostname(), nil
}

func caseInsenstiveContains(a, b string) bool {
	return strings.Contains(strings.ToUpper(a), strings.ToUpper(b))
}
