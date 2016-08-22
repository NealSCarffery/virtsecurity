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

package endpoint

import (
	"encoding/json"
	"net/http"
	"os/exec"

	"github.com/labstack/echo"
)

// ConfigRoute use to config uri and handler of router
func ConfigRoute(e *echo.Echo) {
	e.POST("/seiptables/v1/iptables/rule/add", IptablesRuleAdd)
	e.GET("/seiptables/v1/iptables/help", IptablesHelp)
}

// IptablesRuleAdd handle the get for /seiptables/v1/iptables/rule/add
func IptablesRuleAdd(c echo.Context) error {
	return nil
}

// IptablesHelp handle the get for /seiptables/v1/iptables/help
func IptablesHelp(c echo.Context) error {
	cmd := exec.Command("echo", "-n", `{"Name":"bob", "Age":20}`)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	var person struct {
		Name string
		Age  int
	}
	if err := json.NewDecoder(stdout).Decode(&person); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, &person)
}