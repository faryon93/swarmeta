package swarm

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

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// ---------------------------------------------------------------------------------------
//  imports
// ---------------------------------------------------------------------------------------

type EventFunc func(message events.Message)

// ---------------------------------------------------------------------------------------
//  public functions
// ---------------------------------------------------------------------------------------

func ForEachServiceEvent(docker *client.Client, ctx context.Context, fn EventFunc) {
	filter := filters.NewArgs()
	filter.Add("type", "service")
	evts, err := docker.Events(ctx, types.EventsOptions{
		Filters: filter,
	})

	// upon start a dummy message is issued to get things started
	fn(events.Message{})

	running := true
	for running {
		select {
		case <-err:
			running = false

		case event := <-evts:
			fn(event)
		}
	}
}
