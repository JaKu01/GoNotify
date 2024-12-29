# This multi-stage Dockerfile is used to keep the resulting docker image as small as possible.
# See https://docs.docker.com/build/building/multi-stage/

# The build stage is used to build the application to an executable.
FROM golang:alpine AS build
WORKDIR /app
# The sqlite driver library (gorm.io/driver/sqlite, i.e. indirectly github.com/mattn/go-sqlite3) uses cgo to link sqlite and thus requires the gcc C compiler infrastructure.
RUN apk add gcc libc-dev
COPY . .
RUN go mod download
RUN CGO_ENABLED=1 go build ./cmd/notify

# The execute stage is as small as possible and only copies the executable and files that are actually needed at runtime from the build stage.
FROM alpine
WORKDIR /app
# The musl-libc provided by the alpine image is not sufficient for cgo (used for sqlite), so install a glibc.
RUN apk add --no-cache libc-dev
COPY --from=build /app/template ./template
COPY --from=build /app/notify ./notify
EXPOSE 8080
CMD ["./notify"]