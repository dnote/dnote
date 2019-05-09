/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

var apiProxy *httputil.ReverseProxy

func appHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) >= 2 && parts[1] == "dist" {
		fs := http.StripPrefix("/dist/", http.FileServer(http.Dir("./public/dist")))

		fs.ServeHTTP(w, r)
		return
	}

	// All other requests should go to index.html
	http.ServeFile(w, r, "./public/index.html")
}

func init() {
	if os.Getenv("GO_ENV") != "PRODUCTION" {
		err := godotenv.Load(".env.dev")
		if err != nil {
			panic(errors.Wrap(err, "loading env vars"))
		}
	}

	apiHostURL, err := url.Parse(os.Getenv("API_HOST"))
	if err != nil {
		panic(errors.Wrap(err, "parsing api host url"))
	}

	apiProxy = httputil.NewSingleHostReverseProxy(apiHostURL)
}

func main() {
	http.HandleFunc("/", appHandler)
	http.Handle("/api/", http.StripPrefix("/api/", apiProxy))

	port := os.Getenv("PORT")
	log.Printf("Web listening on port %s", port)

	log.Println(http.ListenAndServe(":"+port, nil))
}
