package model

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
    "bytes"
    "html/template"
)

// ---------------------------------------------------------------------------------------
//  types
// ---------------------------------------------------------------------------------------

type View struct {
    // populated by hcl
    Metadata map[string]*Metadata `json:"metadata"`
    Labels   bool                 `json:"labels"`

    // internal variables
    IsOkay bool `hcl:"-" json:"-"`
}

// ---------------------------------------------------------------------------------------
//  public members
// ---------------------------------------------------------------------------------------

// Initialize initializes (creates templates, ...) this View.
func (v *View) Initialize() error {
    // compile all metadata templates
    for metaName, meta := range v.Metadata {
        var err error

        meta.template, err = template.New(metaName).Parse(meta.TemplateSrc)
        if err != nil {
            return err
        }
    }

    v.IsOkay = true

    return nil
}

// Render renders this view with the given object.
func (v *View) Render(obj interface{}) (map[string]string, error) {
    metadata := make(map[string]string)

    for name, meta := range v.Metadata {
        // render the template to a buffer
        var doc bytes.Buffer
        err := meta.template.Execute(&doc, obj)
        if err != nil {
            return nil, err
        }

        value := string(doc.Bytes())
        if !meta.OmitEmpty || value != "" {
            metadata[name] = value
        }
    }

    return metadata, nil
}
