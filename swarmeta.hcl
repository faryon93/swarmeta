//docker_socket = "unix:///var/run/docker.sock"
docker_socket = "tcp://127.0.0.1:2375"

view "@default" {
    metadata "name" {
        template = "{{.Spec.Name}}"
    }
    metadata "id" {
        template = "{{.ID}}"
    }
    metadata "last_update" {
        template = "{{.UpdatedAt}}"
    }
}
