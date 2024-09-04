# Sensortron Web Interface

Container which provides a minimal web interface to view current
temperature sensor readings and also exposes a REST API for temperature
sensors to submit current readings.

## Setup

Run `podman build -t sensortron .` to build an image.  Example:

    > podman build -t sensortron .
    [1/2] STEP 1/4: FROM docker.io/golang:1.23-alpine AS build
    [1/2] STEP 2/4: COPY . /src
    --> a6f8a0cdf7f
    [1/2] STEP 3/4: WORKDIR /src
    --> 55a02ac1a29
    [1/2] STEP 4/4: RUN --mount=type=cache,target=/go ["go", "build", "-trimpath", "-ldflags=-s -w"]
    --> 82a6777aa86
    [2/2] STEP 1/4: FROM gcr.io/distroless/static
    [2/2] STEP 2/4: ENV TZ "America/New_York"
    --> Using cache 2da230d044b220f223a5603e3450e37c8d740e7159337fa1b6dcd6b00111e6ee
    --> 2da230d044b
    [2/2] STEP 3/4: COPY --from=build /src/sensortron /sensortron
    --> 7a49ad26cdc
    [2/2] STEP 4/4: CMD ["/sensortron"]
    [2/2] COMMIT sensortron
    --> 52e599a7116
    Successfully tagged localhost/sensortron:latest
    52e599a711603a16f0c4fc5b7ad6eba50ba3ec237bbb0387bd9cae0b3fdabb08

Use `podman run` to run the image.  Example:

    > podman run -d --rm -p 1979:1979 -v sensortron:/data --name sensortron sensortron
    a178d699787fcbaf92764b6104cbb4da719c364d406c1fbf69156fa78c13fa41
    > 

Other commands:
- `podman stop sensortron`: stop the container
- `podman logs -f sensortron`: monitor container logs
