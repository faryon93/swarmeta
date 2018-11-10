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
	"net/http"

	"github.com/docker/docker/client"

	"github.com/faryon93/swarmeta/model"
)

// ---------------------------------------------------------------------------------------
//  types
// ---------------------------------------------------------------------------------------

const (
    ctxKey = "__ctx"
)

// ---------------------------------------------------------------------------------------
//  types
// ---------------------------------------------------------------------------------------

type Context struct {
	Conf   *model.Conf
	Docker *client.Client
}

// ---------------------------------------------------------------------------------------
//  public functions
// ---------------------------------------------------------------------------------------

func (c *Context) With( fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKey, c)))
	}
}

func Ctx(r *http.Request) *Context {
    return r.Context().Value(ctxKey).(*Context)
}
