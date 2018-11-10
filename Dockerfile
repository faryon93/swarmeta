# ----------------------------------------------------------------------------------------
# Image: Builder
# ----------------------------------------------------------------------------------------
FROM golang:alpine as builder

# setup the environment
ENV TZ=Europe/Berlin

# install dependencies
RUN apk --update --no-cache add git gcc musl-dev tzdata
WORKDIR /work
ADD ./ ./

# build the go binary
RUN go build -ldflags \
        '-X "main.BuildTime='$(date -Iminutes)'" \
         -X "main.GitCommit='$(git rev-parse --short HEAD)'" \
         -X "main.GitBranch='$(git rev-parse --abbrev-ref HEAD)'" \
         -X "main.BuildNumber='$CI_BUILDNR'" \
         -s -w' \
         -v -o /tmp/swarmeta .

# ----------------------------------------------------------------------------------------
# Image: Deployment
# ----------------------------------------------------------------------------------------
FROM alpine:latest
MAINTAINER Maximilian Pachl <m@ximilian.info>

# setup the environment
ENV TZ=Europe/Berlin

RUN apk --update --no-cache add ca-certificates tzdata bash su-exec curl

# add relevant files to container
COPY --from=builder /tmp/swarmeta /usr/sbin/swarmeta

# make binary executable
RUN chown nobody:nobody /usr/sbin/swarmeta && \
    chmod +x /usr/sbin/swarmeta

HEALTHCHECK --interval=1s --timeout=2s --start-period=1s \
    CMD curl -f http://127.0.0.1:8000/api/v1/health || exit 1

EXPOSE 8000
CMD /usr/sbin/swarmeta
