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
	"context"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/docker/docker/client"
	"github.com/faryon93/util"
	"github.com/hashicorp/hcl"
	"github.com/sirupsen/logrus"

	"github.com/faryon93/swarmeta/model"
	"github.com/faryon93/swarmeta/v1"
)

// ---------------------------------------------------------------------------------------
//  constants
// ---------------------------------------------------------------------------------------

const (
	HttpCloseTimeout = 5 * time.Second
	ConfigName       = "swarmeta.hcl"
)

// ---------------------------------------------------------------------------------------
//  application entry
// ---------------------------------------------------------------------------------------

func main() {
	var err error
	var colors bool
	flag.BoolVar(&colors, "colors", false, "force color logging")
	flag.Parse()

	// setup logger
	formater := logrus.TextFormatter{ForceColors: colors}
	logrus.SetFormatter(&formater)
	logrus.SetOutput(os.Stdout)

	logrus.Infoln("starting", GetAppVersion())

	// read and decode the configureation file
	buf, err := ioutil.ReadFile(ConfigName)
	if err != nil {
		logrus.Errorln("failed to read config file:", err.Error())
		os.Exit(-1)
	}

	var config model.Conf
	err = hcl.Decode(&config, string(buf))
	if err != nil {
		logrus.Errorln("failed to parse config file:", err.Error())
		os.Exit(-1)
	}

	// setup configuration file
	for _, view := range config.Views {
		err = view.Initialize()
		if err != nil {
			logrus.Errorln("failed to initialize view:", err.Error())
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

	// setup http connection
	ctx := v1.Context{Conf: &config, Docker: docker}
	http.HandleFunc("/api/v1/find", ctx.With(v1.Find))
	srv := &http.Server{Addr: ":8000"}
	go func() {
		logrus.Println("http server is listening on :8000")
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logrus.Errorln("http server failed to start:", err.Error())
			return
		}
	}()
	defer func() {
		ctx, _ := context.WithTimeout(context.Background(), HttpCloseTimeout)
		srv.Shutdown(ctx)
		logrus.Infoln("http server shutdown completed")
	}()

	// wait for stop signals
	util.WaitSignal(os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	logrus.Infoln("received SIGINT / SIGTERM going to shutdown")
}
