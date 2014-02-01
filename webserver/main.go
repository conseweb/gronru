// Copyright 2013 gronru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/bmizerany/pat"
	"github.com/globocom/config"
	"github.com/xbee/gronru/api"
	"github.com/xbee/gronru/db"
	"log"
	"net/http"
)

const version = "0.2.0"

func main() {
	dry := flag.Bool("dry", false, "dry-run: does not start the server (for testing purpose)")
	configFile := flag.String("config", "/etc/gronru.conf", "Gronru configuration file")
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
	db.Connect()
	router := pat.New()
	router.Post("/user/:name/key", http.HandlerFunc(api.AddKey))
	router.Del("/user/:name/key/:keyname", http.HandlerFunc(api.RemoveKey))
	router.Get("/user/:name/keys", http.HandlerFunc(api.ListKeys))
	router.Post("/user", http.HandlerFunc(api.NewUser))
	router.Del("/user/:name", http.HandlerFunc(api.RemoveUser))
	router.Post("/repository", http.HandlerFunc(api.NewRepository))
	router.Post("/repository/grant", http.HandlerFunc(api.GrantAccess))
	router.Del("/repository/revoke", http.HandlerFunc(api.RevokeAccess))
	router.Del("/repository/:name", http.HandlerFunc(api.RemoveRepository))
	router.Get("/repository/:name", http.HandlerFunc(api.GetRepository))
	router.Put("/repository/:name", http.HandlerFunc(api.RenameRepository))
	router.Get("/healthcheck/", http.HandlerFunc(api.HealthCheck))

	bind, err := config.GetString("bind")
	if err != nil {
		var perr error
		bind, perr = config.GetString("webserver:port")
		if perr != nil {
			panic(err)
		}
	}
	if !*dry {
		log.Fatal(http.ListenAndServe(bind, router))
	}
}
