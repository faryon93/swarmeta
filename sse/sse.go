package sse

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
	"net/http"

	"github.com/pkg/errors"
)

// ---------------------------------------------------------------------------------------
//  type
// ---------------------------------------------------------------------------------------

// EventStream represents an SSE connection.
type EventStream struct {
	writer http.ResponseWriter
}

// ---------------------------------------------------------------------------------------
//  public functions
// ---------------------------------------------------------------------------------------

// Upgrade upgrades the given http session to an EventStream session.
func Upgrade(w http.ResponseWriter) (*EventStream, error) {
	_, ok := w.(http.Flusher)
	if !ok {
		return nil, errors.New("text/eventstream unsupported")
	}

	// setup headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	return &EventStream{w}, nil
}

// ---------------------------------------------------------------------------------------
//  public members
// ---------------------------------------------------------------------------------------

// Write writes a message to the client.
func (e *EventStream) Write(buf string) {
	e.writer.Write([]byte("data: " + buf + "\r\n\r\n"))
	e.writer.(http.Flusher).Flush()
}
