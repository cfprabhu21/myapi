package main

import (
	"encoding/json"

	"github.com/fulldump/golax"

	fb "github.com/huandu/facebook"

	"net/http"
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

	auth.
		Node("facebook").
		Method("POST", func(c *golax.Context) {
			type User struct {
				FirstName string
				ID string
				Email string
			}

			a := &Auth{}
			json.NewDecoder(c.Request.Body).Decode(a)

			res, _ := fb.Get("/me", fb.Params{
				"fields": "first_name,id,email",
				"access_token": a.Token,
			})


			user := &User{}
			res.Decode(user)

			
			json.NewEncoder(c.Response).Encode(map[string]interface{}{
				"id": user.ID,
				"email": user.Email,
				"name": user.FirstName,
			})

		})

	auth.
		Node("google").
		Method("POST", func(c *golax.Context) {
			type User struct {
				Given_Name string
				ID string
				Email string
				Picture string
				Gender string
				Locale string
			}
			
			a := &Auth{}
			json.NewDecoder(c.Request.Body).Decode(a)

			res, _ := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + a.Token)

			user := &User{}

			json.NewDecoder(res.Body).Decode(user)

			json.NewEncoder(c.Response).Encode(map[string]interface{}{
				"id": user.ID,
				"email": user.Email,
				"name": user.Given_Name,
			})
		})

	auth.
		Node("linkedin").
		Method("POST", func(c *golax.Context) {
			
			type User struct {
				FirstName string
				ID string
				EmailAddress string
			}

			a := &Auth{}
			json.NewDecoder(c.Request.Body).Decode(a)

			res, _ := http.Get("https://api.linkedin.com/v1/people/~:(id,first-name,picture-url,email-address)?format=json&oauth2_access_token=" + a.Token)
			user := &User{}
			json.NewDecoder(res.Body).Decode(user)
			json.NewEncoder(c.Response).Encode(map[string]interface{}{
				"id": user.ID,
				"email": user.EmailAddress,
				"name": user.FirstName,
			})

		})

	my_api.Serve()
}
