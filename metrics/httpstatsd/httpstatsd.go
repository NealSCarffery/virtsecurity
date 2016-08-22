// Cosmos os system team support for Cosmosos - Cosmos os system for cloud
//
// Copyright 2010 The Cosmos Authors.  All rights reserved.
// https://github.com/cosmosos
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package httpstatsd

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// HTTPStats define the structure of status of HTTP
type HTTPStats struct {
	HTTPMethods map[string]int `json:"HTTPMethods"`
	MIMETypes   map[string]int `json:"MIMETypes"`
	HTTPStatus  map[string]int `json:"HTTPStatus"`
	RemoteAddrs map[string]int `json:"RemoteAddrs"`
	RouteURI    map[string]int `json:"RouteURI"`
}

type (
	// Config define the structure of status of HTTP
	Config struct {
		Skipper      middleware.Skipper
		RouteSkipper RouteSkipper
	}

	// RouteSkipper define skip map which will not be caculated
	RouteSkipper map[string]bool
)

var (
	// DefaultConfig is the default HTTPStats middleware config.
	DefaultConfig = Config{
		Skipper: func(c echo.Context) bool {
			return false
		},
	}

	// DefaultStats is the default HTTPStats middleware data.
	DefaultStats = HTTPStats{
		HTTPMethods: make(map[string]int),
		MIMETypes:   make(map[string]int),
		HTTPStatus:  make(map[string]int),
		RemoteAddrs: make(map[string]int),
		RouteURI:    make(map[string]int),
	}
)

// HTTPStatsd returns a middleware that stats HTTP requests.
func HTTPStatsd(rs RouteSkipper) echo.MiddlewareFunc {
	DefaultConfig.RouteSkipper = rs
	return WithConfig(DefaultConfig, DefaultStats)
}

// WithConfig returns a HTTPStatsd middleware from config.
// See: `Logger()`.
func WithConfig(config Config, stats HTTPStats) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()

			stats.RouteURI[req.URI()]++

			if config.RouteSkipper[req.URI()] {
				return next(c)
			}

			res := c.Response()

			stats.HTTPMethods[req.Method()]++
			stats.MIMETypes[req.Header().Get(echo.HeaderContentType)]++
			stats.RemoteAddrs[req.RemoteAddress()]++

			if err = next(c); err != nil {
				c.Error(err)
			}

			stats.HTTPStatus[http.StatusText(res.Status())]++

			return
		}
	}
}
