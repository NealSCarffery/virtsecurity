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

package middleware

import (
	"os"
	"path"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/virtsecurity/logger"
)

// NewLoggerMiddleware create a instance of logger middleware
func NewLoggerMiddleware(logdir string, logfilename string) (mf echo.MiddlewareFunc, f *os.File) {
	if logdir == "" {
		logger.Logger.Info("Logdir is empty and use default stdout logger")
		mf = middleware.Logger()
		f = nil
		return
	}

	if _, err := os.Stat(logdir); err != nil && os.IsNotExist(err) {
		logger.Logger.Info("Logdir is not exist. Make the directory")
		err = os.MkdirAll(logdir, 0666)
		if err != nil {
			logger.Logger.Info("Make the directory failed and use default stdout logger", err)
			mf = middleware.Logger()
			f = nil
			return
		}
	}

	logpath := path.Join(logdir, logfilename)

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
