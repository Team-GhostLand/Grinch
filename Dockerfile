FROM golang:1.24 AS builder
WORKDIR /src/Team-GhostLand/Grinch/
COPY . /src/Team-GhostLand/Grinch
RUN ["go", "build", "-o", "/bin/grinch"]

FROM alpine/git
WORKDIR /app
COPY --from=builder /bin/grinch /app
COPY ./scripts/ci.sh /app
COPY ./scripts/make_serverpack.sh /app
RUN ["chmod", "-R", "777", "."]

ENTRYPOINT ["/app/ci.sh"]