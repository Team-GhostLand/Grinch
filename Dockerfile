FROM golang:1.24 AS builder
WORKDIR /src/Team-GhostLand/Grinch/
COPY . /src/Team-GhostLand/Grinch
RUN ["go", "build", "-o", "/bin/grinch"]
RUN ["/bin/grinch"]

FROM alpine/git
WORKDIR /app
RUN ["apk", "add", "bash"]
COPY --from=builder /bin/grinch /app
RUN ["/app/grinch"]
COPY ./scripts/ci.sh /app
COPY ./scripts/make_serverpack.sh /app
RUN ["chmod", "-R", "777", "."]

ENTRYPOINT ["/app/ci.sh"]