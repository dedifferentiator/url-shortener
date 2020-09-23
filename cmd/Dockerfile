FROM golang:1.14.6-alpine3.12 as stage

WORKDIR /cmd

COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

COPY . .
RUN go build -o main


FROM alpine:3.12.0 as prod

WORKDIR /cmd

COPY --from=stage /cmd/main .
COPY --from=stage /cmd/static static/
ENTRYPOINT ["./main"]
