package main

import (
	"encoding/json"

	"github.com/fulldump/golax"

	"net/http"

	"fmt"

	"io/ioutil"
)

type Auth struct{
	Token string `json:"token"`
}

func main() {

	my_api := golax.NewApi()

	my_api.Root.
		Interceptor(golax.InterceptorLog).
		Interceptor(golax.InterceptorError)

	auth := my_api.Root.Node("auth")

	// Create a http client for linkedin
	http_client_facebook := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: false,
		},
	}
	auth.
		Node("facebook").
		Method("POST", func(c *golax.Context) {

			a := &Auth{}
			json.NewDecoder(c.Request.Body).Decode(a)

			// TODO: urlencode/validate a.Token
			res, res_error := http_client_facebook.Get("https://graph.facebook.com/me?fields=id,name,email&access_token=" + a.Token)
			
			if nil != res_error {
				c.Error(http.StatusBadGateway, "something went wrong")
				return
			}			

			// SOme kind of validation
			if res.StatusCode != http.StatusOK {
				// do something
				c.Error(http.StatusBadGateway, "Facebook not available")
				return
			}

			body, _ := ioutil.ReadAll(res.Body) // This is []byte

			fmt.Fprintln(c.Response, string(body))	

			res.Body.Close() // Important to allow connection reusage (keep alive HTTP/1.1)

		})
	
	// Create a http client for google
	http_client_google := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: false,
		},
	}

	auth.
		Node("google").
		Method("POST", func(c *golax.Context) {
			
			a := &Auth{}
			json.NewDecoder(c.Request.Body).Decode(a)


			// TODO: urlencode/validate a.Token
			res, res_error := http_client_google.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + a.Token)
			
			if nil != res_error {
				c.Error(http.StatusBadGateway, "something went wrong")
				return
			}			

			// SOme kind of validation
			if res.StatusCode != http.StatusOK {
				// do something
				c.Error(http.StatusBadGateway, "Google not available")
				return
			}

			body, _ := ioutil.ReadAll(res.Body) // This is []byte

			fmt.Fprintln(c.Response, string(body))	

			res.Body.Close() // Important to allow connection reusage (keep alive HTTP/1.1)
		})


	// Create a http client for facebook
	http_client_linkedin := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: false,
		},
	}

	auth.
		Node("linkedin").
		Method("POST", func(c *golax.Context) {

			a := &Auth{}
			json.NewDecoder(c.Request.Body).Decode(a)

			// TODO: urlencode/validate a.Token
			res, res_error := http_client_linkedin.Get("https://api.linkedin.com/v1/people/~:(id,first-name,picture-url,email-address)?format=json&oauth2_access_token=" + a.Token)
			
			if nil != res_error {
				c.Error(http.StatusBadGateway, "something went wrong")
				return
			}			

			// SOme kind of validation
			if res.StatusCode != http.StatusOK {
				// do something
				c.Error(http.StatusBadGateway, "linkedin not available")
				return
			}

			body, _ := ioutil.ReadAll(res.Body) // This is []byte

			fmt.Fprintln(c.Response, string(body))	

			res.Body.Close() // Important to allow connection reusage (keep alive HTTP/1.1)
		})

	my_api.Serve()
}
