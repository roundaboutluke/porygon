# Build image
FROM golang:1.21-alpine as build

WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN if [ ! -f vendor/modules.txt ]; then go mod download; fi

COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/porygon

# Now copy it into our base image.
FROM gcr.io/distroless/static-debian12 as runner
COPY --from=build /go/bin/porygon /porygon/
COPY default.toml /porygon/
COPY /templates /porygon/templates

WORKDIR /porygon
CMD ["./porygon"]