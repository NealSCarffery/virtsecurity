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

package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
	"github.com/virtsecurity/seovs/cmd"
)

func main() {
	cli.AppHelpTemplate = fmt.Sprintf(`%s
    website: http://www.cosmosos.org
    support: support@cosmosos.org
    By. Cosmos os system team`, cli.AppHelpTemplate)

	cli.VersionFlag = cli.BoolFlag{
		Name:  "print-version, V",
		Usage: "print only the version",
	}

	app := cli.NewApp()
	app.Name = "security openvswitch"
	app.Usage = "The daemon of security openvswitch"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "run security openvswitch",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "bind, b",
					Value:       "0.0.0.0:30720",
					Usage:       "security openvswitch bind to `addr:port`",
					Destination: &cmd.Bind,
					EnvVar:      "SECURITY_OPENVSWITCH_BIND",
				},

				cli.StringFlag{
					Name:        "logdir, l",
					Value:       "",
					Usage:       "security openvswitch log to `directory`",
					Destination: &cmd.Logdir,
					EnvVar:      "SECURITY_OPENVSWITCH_LOGDIR",
				},
			},
			Action: cmd.RunAction,
		},

		{
			Name:    "httpstatsd",
			Aliases: []string{"h"},
			Usage:   "Get http status",
			Action:  cmd.HTTPStatsdAction,
		},

		{
			Name:    "shutdown",
			Aliases: []string{"s"},
			Usage:   "shutdown security openvswitch",
			Action:  cmd.ShutdownAction,
		},
	}

	app.Action = func(c *cli.Context) error {
		fmt.Println("There are not commands or options")
		fmt.Println("See [seovs help]")
		return nil
	}

	app.Run(os.Args)
}
