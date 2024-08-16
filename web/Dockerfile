FROM docker.io/golang:1.23-alpine AS build
COPY . /src
WORKDIR /src
RUN ["go", "build", "-trimpath", "-ldflags=-s -w"]

# used to use "scratch", but we need ca-certificates in order to connect
# to the NWS API.
#
# FROM scratch
FROM gcr.io/distroless/static
COPY --from=build /src/sensortron /sensortron
CMD ["/sensortron"]

