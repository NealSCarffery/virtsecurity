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

package cmd

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/go-resty/resty"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	vsmiddleware "github.com/virtsecurity/middleware"

	"github.com/urfave/cli"
	"github.com/virtsecurity/daemon"
	httpecho "github.com/virtsecurity/http/echo"
	"github.com/virtsecurity/logger"
	"github.com/virtsecurity/metrics/httpstatsd"
	"github.com/virtsecurity/seovs/endpoint"
	"github.com/virtsecurity/seovs/server"
	udsecho "github.com/virtsecurity/uds/echo"
)

// Bind the security openvswitch to addr and port
var Bind string

// Logdir is the directory where is to save the log
var Logdir string

const (
	// LogFileName define the name of log file
	LogFileName = "seovs.log"
)

func setEchoLoggerOutput(e *echo.Echo, f *os.File) {
	if f != nil {
		e.Logger().SetOutput(f)
	}
}

// RunAction perform start the security openvswitch daemon
func RunAction(c *cli.Context) error {
	var err error

	logger.Logger.Info("init echo server")
	e := echo.New()

	logger.Logger.Info("init logger middleware")
	mf, f := vsmiddleware.NewLoggerMiddleware(Logdir, LogFileName)

	logger.Logger.Info("init echo server logger")
	setEchoLoggerOutput(e, f)

	if f != nil {
		defer f.Close()
	}

	logger.Logger.Info("init echo server logger middleware")
	e.Use(mf)
	logger.Logger.Info("init echo server recover middleware")
	e.Use(middleware.Recover())
	e.Use(httpstatsd.HTTPStatsd(endpoint.NewRouteSkipper()))

	logger.Logger.Info("init echo server router")
	endpoint.ConfigRoute(e)

	logger.Logger.Info("Start echo server")

	server.HTTPServer.Server, server.HTTPServer.LN, err = httpecho.NewHTTPEchoServer(Bind)
	if err != nil {
		return err
	}

	server.UDSServer.Server, server.UDSServer.LN, err = udsecho.NewUDSEchoServer("/var/run/seovs.sock")
	if err != nil {
		return err
	}

	go e.Run(server.HTTPServer.Server)
	go e.Run(server.UDSServer.Server)

	daemon.Run()

	logger.Logger.Info("Shutdown echo server")
	return nil
}

// HTTPStatsdAction get the status from httpstatsd
func HTTPStatsdAction(c *cli.Context) error {
	transport := http.Transport{
		Dial: func(_, _ string) (net.Conn, error) {
			return net.Dial("unix", "/var/run/seovs.sock")
		},
	}

	r := resty.New().
		SetTransport(&transport).
		SetScheme("http")

	resp, err := r.R().
		SetHeader(echo.HeaderContentType, echo.MIMEApplicationJSON).
		Get("/seovs/v1/ovs/http/status")

	if err != nil {
		return err
	}

	fmt.Printf("%s\n", resp.String())

	return nil
}

// ShutdownAction perform shutdown the security openvswitch daemon
func ShutdownAction(c *cli.Context) error {
	transport := http.Transport{
		Dial: func(_, _ string) (net.Conn, error) {
			return net.Dial("unix", "/var/run/seovs.sock")
		},
	}

	r := resty.New().
		SetTransport(&transport).
		SetScheme("http")

	resp, err := r.R().
		SetHeader(echo.HeaderContentType, echo.MIMEApplicationJSON).
		Get("/seovs/v1/ovs/shutdown")

	if err != nil {
		return err
	}

	fmt.Printf("%s\n", resp.String())

	return nil
}
