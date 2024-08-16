FROM docker.io/golang:1.23-alpine AS build
COPY . /src
WORKDIR /src
RUN ["go", "build", "-trimpath", "-ldflags=-s -w"]

FROM scratch
COPY --from=build /src/sensortron /sensortron
CMD ["/sensortron"]

