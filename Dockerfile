FROM golang:1.18.0-alpine3.15 AS build

WORKDIR /app

RUN adduser -D scratchuser

COPY go.* ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 go build -o /noman -ldflags="-s -w"

FROM scratch

WORKDIR /www

USER scratchuser

COPY --from=0 /etc/passwd /etc/passwd
COPY --from=build /noman /noman

ENTRYPOINT ["/noman"]
