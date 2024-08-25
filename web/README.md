# Sensortron Web Interface

Container which provides a minimal web interface to view current
temperature sensor readings and also exposes a REST API for temperature
sensors to submit current readings.

## Setup

Run `podman build -t sensortron .` to build an image.  Example:

    > podman build -t sensortron .
    [1/2] STEP 1/4: FROM docker.io/golang:1.23-alpine AS build
    [1/2] STEP 2/4: COPY . /src
    --> 0b3fceee26a
    [1/2] STEP 3/4: WORKDIR /src
    --> 0efba939713
    [1/2] STEP 4/4: RUN ["go", "build", "-trimpath", "-ldflags=-s -w"]
    go: downloading github.com/go-chi/chi/v5 v5.1.0
    --> a0c5c1a6cb5
    [2/2] STEP 1/3: FROM scratch
    [2/2] STEP 2/3: COPY --from=build /src/sensors /sensors
    --> 16bd359fe1f
    [2/2] STEP 3/3: CMD ["/sensors"]
    [2/2] COMMIT sensors
    --> dbc781e8bd6
    Successfully tagged localhost/sensortron:latest
    dbc781e8bd63b4b3972882359ea588c4cf690dbf6f335040342b21e0673cbb02

Use `podman run` to run the image.  Example:

    > podman run -d --rm -p 1979:1979 -v sensortron:/data --name sensortron sensortron
    a178d699787fcbaf92764b6104cbb4da719c364d406c1fbf69156fa78c13fa41
    > 

Other commands:
- `podman stop sensortron`: stop the container
- `podman logs -f sensortron`: monitor container logs
