package v1

// swarmeta - docker swarm service metadata
// Copyright (C) 2018 Maximilian Pachl

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// ---------------------------------------------------------------------------------------
//  imports
// ---------------------------------------------------------------------------------------

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/faryon93/util"
)

// ---------------------------------------------------------------------------------------
//  public functions
// ---------------------------------------------------------------------------------------

// Health renders the health page.
func Health(w http.ResponseWriter, r *http.Request) {
	// only the healthcheck script should have access to the health page
	addr := util.GetRemoteAddr(r)
	if addr != "127.0.0.1" && addr != "::1" {
		logrus.WithField("addr", addr).Warnln("rejected access to health page")
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	// try docker
	_, err := Ctx(r).Docker.Ping(r.Context())
	if err != nil {
		logrus.Errorln("healthcheck failed:", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "{\"status\": \"ok\"}")
}
