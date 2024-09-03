# build stage
FROM docker.io/golang:1.23-alpine AS build
COPY . /src
WORKDIR /src

# cache /go between builds to cache packages and improve build speed
RUN --mount=type=cache,target=/go ["go", "build", "-trimpath", "-ldflags=-s -w"]

# run stage
#
# this stage used to use "FROM scratch", but we need ca-certificates in
# order to connect to the NWS API.
FROM gcr.io/distroless/static
ENV TZ "America/New_York"
COPY --from=build /src/sensortron /sensortron
CMD ["/sensortron"]
