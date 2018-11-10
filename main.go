package main

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
	"io/ioutil"
	"log"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/hashicorp/hcl"

	"github.com/faryon93/swarmeta/model"
	"github.com/faryon93/swarmeta/v1"
)

// ---------------------------------------------------------------------------------------
//  application entry
// ---------------------------------------------------------------------------------------

func main() {
	var err error

	buf, err := ioutil.ReadFile("swarmeta.hcl")
	if err != nil {
		panic(err)
	}

	var config model.Conf
	err = hcl.Decode(&config, string(buf))
	if err != nil {
		panic(err)
	}

	// setup configuration file
	for _, view := range config.Views {
		err = view.Initialize()
		if err != nil {
			log.Println("failed to initialize view:", err.Error())
			continue
		}
	}

	// create the docker client
	host := config.DockerSocket
	version := client.DefaultVersion
	docker, err := client.NewClient(host, version, nil, nil)
	if err != nil {
		panic(err)
	}

	ctx := v1.Context{Conf: &config, Docker: docker}
	http.HandleFunc("/api/v1/find", ctx.With(v1.Find))
	log.Fatal(http.ListenAndServe(":8000", nil))
}
