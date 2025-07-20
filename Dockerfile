FROM golang:1.24 AS builder
WORKDIR /src/Team-GhostLand/Grinch/
COPY . /src/Team-GhostLand/Grinch
RUN ["pwd"]
RUN ["ls", "-al"]
RUN ["go", "build", "-o", "/bin/grinch"]

FROM alpine/git
WORKDIR /app
COPY --from=builder /bin/grinch /app
COPY ./scripts/ci.sh /app
COPY ./scripts/make_serverpack.sh /app

ENTRYPOINT ["/app/ci.sh"]