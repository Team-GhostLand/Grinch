FROM golang:1.24 as builder
WORKDIR /src
COPY ./main.go /src
COPY ./go.sum /src
COPY ./go.mod /src
COPY ./cmd/ /src
COPY ./trans/ /src
COPY ./util/ /src
RUN ["go", "build", "-o", "/bin/grinch"]

FROM alpine/git
WORKDIR /app
COPY --from=builder /bin/grinch /app
COPY ./scripts/ci.sh /app
COPY ./scripts/make_serverpack.sh /app

ENTRYPOINT ["/app/ci.sh"]