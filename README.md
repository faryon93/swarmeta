# Swarmeta
Swarmeta is a metadata service discovery for docker swarm. Services can be queried by arbitrary labels. The responses are rendered as predefined views.

## Configuration

    docker_socket = "unix:///var/run/docker.sock"

    view "@default" {
        metadata "name" {
            template = "{{.Spec.Name}}"
            OmitEmpty = true
        }
        metadata "id" {
            template = "{{.ID}}"
        }
        metadata "last_update" {
            template = "{{.UpdatedAt}}"
        }
    }

To get a complete list of properties you can use in the metadata templates see: https://godoc.org/github.com/docker/docker/api/types/swarm#Service

## Query
Just send a request to `/api/v1/find/` all query parameters are interpreted as swarm service label matches.
To get an EventStream (SSE) append `?_follow=true`. To select another view as `@default` you can set the query parameter `?_view=<viewname>`.
