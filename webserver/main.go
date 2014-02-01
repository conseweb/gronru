// Copyright 2013 gronru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	// "github.com/bmizerany/pat"
	"github.com/globocom/config"
	"github.com/xbee/gronru/api"
	"github.com/xbee/gronru/db"
	"log"
	"net/http"

	"github.com/codegangsta/martini"
	// "github.com/codegangsta/martini-contrib/auth"
)

const version = "0.2.0"

var bind string
var dry *bool

// The one and only access token! In real-life scenarios, a more complex authentication
// middleware than auth.Basic should be used, obviously.
const AuthToken = "token"

// The one and only martini instance.
var m *martini.Martini

// func init() {
// 	m = martini.New()
// 	// Setup middleware
// 	m.Use(martini.Recovery())
// 	m.Use(martini.Logger())
// 	// m.Use(auth.Basic(AuthToken, ""))
// 	// m.Use(MapEncoder)
// 	// Setup routes
// 	r := martini.NewRouter()

// 	r.Get(`/albums`, GetAlbums)
// 	r.Get(`/albums/:id`, GetAlbum)
// 	r.Post(`/albums`, AddAlbum)
// 	r.Put(`/albums/:id`, UpdateAlbum)
// 	r.Delete(`/albums/:id`, DeleteAlbum)

// 	r.Post(`/users`, CreateUser)

// 	// Inject database
// 	m.MapTo(db, (*DB)(nil))
// 	// Add the router action
// 	m.Action(r.Handle)
// }

func init() {
	dry = flag.Bool("dry", false, "dry-run: does not start the server (for testing purpose)")
	configFile := flag.String("config", "./etc/gronru.conf", "Gronru configuration file")
	gVersion := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *gVersion {
		fmt.Printf("gronru-webserver version %s\n", version)
		return
	}

	err := config.ReadAndWatchConfigFile(*configFile)
	if err != nil {
		msg := `Could not find gronru config file. Searched on %s.
For an example conf check gronru/etc/gronru.conf file.\n %s`
		log.Panicf(msg, *configFile, err)
	}

	bind, err = config.GetString("bind")
	if err != nil {
		var perr error
		bind, perr = config.GetString("webserver:port")
		if perr != nil {
			panic(err)
		}
	}

	db.Connect()

	m = martini.New()
	// Setup middleware
	m.Use(martini.Recovery())
	m.Use(martini.Logger())
	// m.Use(auth.Basic(AuthToken, ""))
	// m.Use(MapEncoder)
	// Setup routes
	router := martini.NewRouter()

	// router := pat.New()
	// router.Post("/user/:name/key", http.HandlerFunc(api.AddKey))
	// router.Del("/user/:name/key/:keyname", http.HandlerFunc(api.RemoveKey))
	// router.Get("/user/:name/keys", http.HandlerFunc(api.ListKeys))
	// router.Post("/user", http.HandlerFunc(api.NewUser))
	// router.Del("/user/:name", http.HandlerFunc(api.RemoveUser))
	// router.Post("/repository", http.HandlerFunc(api.NewRepository))
	// router.Post("/repository/grant", http.HandlerFunc(api.GrantAccess))
	// router.Del("/repository/revoke", http.HandlerFunc(api.RevokeAccess))
	// router.Del("/repository/:name", http.HandlerFunc(api.RemoveRepository))
	// router.Get("/repository/:name", http.HandlerFunc(api.GetRepository))
	// router.Put("/repository/:name", http.HandlerFunc(api.RenameRepository))
	// router.Get("/healthcheck/", http.HandlerFunc(api.HealthCheck))

	router.Post("/user/:name/key", api.AddKey)
	router.Delete("/user/:name/key/:keyname", api.RemoveKey)
	router.Get("/user/:name/keys", api.ListKeys)
	router.Post("/user", api.NewUser)
	router.Delete("/user/:name", api.RemoveUser)
	router.Post("/repository", api.NewRepository)
	router.Post("/repository/grant", api.GrantAccess)
	router.Delete("/repository/revoke", api.RevokeAccess)
	router.Delete("/repository/:name", api.RemoveRepository)
	router.Get("/repository/:name", api.GetRepository)
	router.Put("/repository/:name", api.RenameRepository)
	router.Get("/healthcheck/", api.HealthCheck)

	// r.Get(`/albums`, GetAlbums)
	// r.Get(`/albums/:id`, GetAlbum)
	// r.Post(`/albums`, AddAlbum)
	// r.Put(`/albums/:id`, UpdateAlbum)
	// r.Delete(`/albums/:id`, DeleteAlbum)

	// r.Post(`/users`, CreateUser)

	// Inject database
	// m.MapTo(db, (*DB)(nil))
	// Add the router action
	m.Action(router.Handle)

}

func main() {

	// go func() {
	// 	// Listen on http: to raise an error and indicate that https: is required.
	// 	//
	// 	// This could also be achieved by passing the same `m` martini instance as
	// 	// used by the https server, and by using a middleware that checks for https
	// 	// and returns an error if it is not a secure connection. This would have the benefit
	// 	// of handling only the defined routes. However, it is common practice to define
	// 	// APIs on separate web servers from the web (html) pages, for maintenance and
	// 	// scalability purposes, so it's not like it will block otherwise valid routes.
	// 	//
	// 	// It is also common practice to use a different subdomain so that cookies are
	// 	// not transfered with every API request.
	// 	// So with that in mind, it seems reasonable to refuse each and every request
	// 	// on the non-https server, regardless of the route. This could of course be done
	// 	// on a reverse-proxy in front of this web server.
	// 	//
	// 	if !*dry {
	// 		if err := http.ListenAndServe(bind, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 			http.Error(w, "https scheme is required", http.StatusBadRequest)
	// 		})); err != nil {
	// 			log.Fatal(err)
	// 		}
	// 	}
	// }()

	// Listen on https: with the preconfigured martini instance. The certificate files
	// can be created using this command in this repository's root directory:
	//
	// go run /path/to/goroot/src/pkg/crypto/tls/generate_cert.go --host="localhost"
	//
	if !*dry {
		// log.Fatal(http.ListenAndServe(bind, router))
		// if err := http.ListenAndServeTLS(bind, "cert.pem", "key.pem", m); err != nil {
		if err := http.ListenAndServe(bind, m); err != nil {
			log.Fatal(err)
		}
	}

}
