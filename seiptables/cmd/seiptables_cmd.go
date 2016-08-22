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
	"os"
	"path"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"github.com/urfave/cli"
	"github.com/virtsecurity/logger"
	"github.com/virtsecurity/seiptables/endpoint"
)

// Bind the security iptables to addr and port
var Bind string

// Logdir is the directory where is to save the log
var Logdir string

const (
	// LogFileName define the name of log file
	LogFileName = "seiptables.log"
)

func newLoggerMiddleware() (mf echo.MiddlewareFunc, f *os.File) {
	if Logdir == "" {
		logger.Logger.Info("Logdir is empty and use default stdout logger")
		mf = middleware.Logger()
		f = nil
		return
	}

	if _, err := os.Stat(Logdir); err != nil && os.IsNotExist(err) {
		logger.Logger.Info("Logdir is not exist. Make the directory")
		err = os.MkdirAll(Logdir, 0666)
		if err != nil {
			logger.Logger.Info("Make the directory failed and use default stdout logger", err)
			mf = middleware.Logger()
			f = nil
			return
		}
	}

	logpath := path.Join(Logdir, LogFileName)

	file, err := os.OpenFile(logpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		logger.Logger.Info("Open log file failed and use default stdout logger", err)
		mf = middleware.Logger()
		f = nil
		return
	}

	logger.Logger.Info("Use file logger which location at " + logpath)

	loggerConfig := middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			return false
		},
		Format: `{"time":"${time_rfc3339}","remote_ip":"${remote_ip}",` +
			`"method":"${method}","uri":"${uri}","status":${status}, "latency":${latency},` +
			`"latency_human":"${latency_human}","bytes_in":${bytes_in},` +
			`"bytes_out":${bytes_out}}` + "\n",
		Output: file,
	}

	mf = middleware.LoggerWithConfig(loggerConfig)
	f = file

	return
}

func setEchoLoggerOutput(e *echo.Echo, f *os.File) {
	if f != nil {
		e.Logger().SetOutput(f)
	}
}

// StartAction perform start the security iptables daemon
func StartAction(c *cli.Context) error {
	logger.Logger.Info("init echo server")
	e := echo.New()

	logger.Logger.Info("init logger middleware")
	mf, f := newLoggerMiddleware()

	logger.Logger.Info("init echo server logger")
	setEchoLoggerOutput(e, f)

	if f != nil {
		defer f.Close()
	}

	logger.Logger.Info("init echo server logger middleware")
	e.Use(mf)
	logger.Logger.Info("init echo server recover middleware")
	e.Use(middleware.Recover())

	endpoint.ConfigRoute(e)

	logger.Logger.Info("Start echo server")
	e.Run(standard.New(Bind))
	logger.Logger.Info("Shutdown echo server")
	return nil
}
