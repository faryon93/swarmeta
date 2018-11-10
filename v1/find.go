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
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"

	"github.com/cespare/xxhash"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"

	"github.com/faryon93/util"

	"github.com/faryon93/swarmeta/model"
	"github.com/faryon93/swarmeta/sse"
	"github.com/faryon93/swarmeta/swarm"
)

// ---------------------------------------------------------------------------------------
//  constants
// ---------------------------------------------------------------------------------------

const (
	DefaultView = "@default"

	QueryMagicPrefix = "_"
	QueryFollow      = QueryMagicPrefix + "follow"
	QueryView        = QueryMagicPrefix + "view"
)

// ---------------------------------------------------------------------------------------
//  types
// ---------------------------------------------------------------------------------------

type meta map[string]string

// ---------------------------------------------------------------------------------------
//  public functions
// ---------------------------------------------------------------------------------------

func Find(w http.ResponseWriter, r *http.Request) {
	docker := Ctx(r).Docker
	view := getView(r)
	filter := getFilter(r)
	if filter.Len() == 0 {
		http.Error(w, "filter required", http.StatusNotAcceptable)
		return
	}

	// the user wants a one time answer
	if r.URL.Query().Get(QueryFollow) == "" {
		services, err := renderMeta(docker, view, filter)
		if err != nil {
			logrus.Errorln("failed to fetch service list:", err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		util.Jsonify(w, services)

	} else {
		// stream the service list as Event Stream (Server-Sent-Events)
		stream, err := sse.Upgrade(w)
		if err != nil {
			http.Error(w, "streaming unsupported!", http.StatusBadRequest)
			return
		}

		// everytime a service is updated/created/deleted
		// the list is sent to the client
		// TODO: wouldn't it be better to just have one central event listener?
		lastHash := uint64(0)
		ctx, cancel := context.WithCancel(r.Context())
		swarm.ForEachServiceEvent(docker, ctx, func(events.Message) {
			services, err := renderMeta(docker, view, filter)
			if err != nil {
				logrus.Errorln("failed to fetch service list:", err.Error())
				cancel()
				return
			}

			// transform to json
			buf, err := json.Marshal(services)
			if err != nil {
				logrus.Errorln("failed to transform json:", err.Error())
				cancel()
				return
			}

			// only send the update if something has actually chanced
			hash := xxhash.Sum64(buf)
			if lastHash != hash {
				lastHash = hash
				stream.Write(string(buf))
			}
		})
	}
}

// ---------------------------------------------------------------------------------------
//  helper functions
// ---------------------------------------------------------------------------------------

// getFilter extracts the docker label filter from request parameters.
func getFilter(r *http.Request) filters.Args {
	filter := filters.NewArgs()

	for labelName, labelValues := range r.URL.Query() {
		if strings.HasPrefix(labelName, QueryMagicPrefix) {
			continue
		}

		for _, value := range labelValues {
			filter.Add("label", labelName+"="+value)
		}
	}

	return filter
}

// getView returns the view configured by the user.
// If an invalid view was requests by the user the default
// view is returned.
func getView(r *http.Request) *model.View {
	viewName := r.URL.Query().Get(QueryView)
	conf := Ctx(r).Conf

	// if the requests view could not be
	// processed (e.g. not defined, template errors...)
	// we use the default view for processing.
	view, ok := conf.Views[viewName]
	if !ok || view == nil || !view.IsOkay {
		view = conf.Views[DefaultView]
	}

	return view
}

// renderMeta return the rendered metadata for the matching
// swarm services. The metadata is generated with the given view.
func renderMeta(dock *client.Client, vw *model.View, filt filters.Args) ([]meta, error) {
	// find all services matching the filter
	list, err := dock.ServiceList(context.Background(), types.ServiceListOptions{
		Filters: filt,
	})
	if err != nil {
		return nil, err
	}

	// transform the service list according to the view definition
	services := make([]meta, len(list))
	for i, service := range list {
		services[i], err = vw.Render(&service.Spec)
		if err != nil {

			return nil, err
		}
	}

	return services, nil
}
