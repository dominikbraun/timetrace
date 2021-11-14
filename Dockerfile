# This Dockerfile builds a lightweight distribution image for Docker Hub.
# It only contains the application without any source code.
FROM golang:latest as builder

# The timetrace release to be downloaded from GitHub.
ARG VERSION
WORKDIR /
COPY *.go go.* ./
RUN CGO_ENABLED=0 go build -mod=mod -a -installsuffix temp -ldflags "-extldflags '-static' -X 'main.version=${VERSION}'" .


# The final stage. This is the image that will be distrubuted.
FROM alpine:3.11.5 AS final

RUN apk add -U --no-cache tzdata

LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.name="timetrace"
LABEL org.label-schema.description="A simple CLI for tracking your working time."
LABEL org.label-schema.url="https://github.com/dominikbraun/timetrace"
LABEL org.label-schema.vcs-url="https://github.com/dominikbraun/timetrace"
LABEL org.label-schema.version=${VERSION}

COPY --from=builder ["/timetrace", "/bin/timetrace"]

# Create a symlink for musl, see https://stackoverflow.com/a/35613430.
#RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

RUN mkdir /etc/timetrace && \
    echo "store: '/data'" >> /etc/timetrace/config.yml

RUN mkdir /data

ENTRYPOINT ["/bin/timetrace"]
