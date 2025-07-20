FROM golang:1.24 AS builder
WORKDIR /src/Team-GhostLand/Grinch/
COPY . /src/Team-GhostLand/Grinch
RUN ["go", "build", "-o", "/bin/grinch"]

FROM alpine/git
WORKDIR /app
RUN ["apk", "add", "bash"]
RUN ["apk", "add", "libc6-compat"]
COPY --from=builder /bin/grinch /app
COPY ./scripts/ci.sh /app
COPY ./scripts/make_serverpack.sh /app
RUN ["chmod", "-R", "777", "."]
RUN ["mv", "make_serverpack.sh", "grinch-serverpack"]

ENTRYPOINT ["/app/ci.sh"]