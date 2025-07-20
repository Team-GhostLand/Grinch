FROM golang:1.24 AS builder
WORKDIR /src/Team-GhostLand/Grinch/Team-GhostLand/Grinch
COPY ./main.go /src/Team-GhostLand/Grinch
COPY ./go.sum /src/Team-GhostLand/Grinch
COPY ./go.mod /src/Team-GhostLand/Grinch
COPY ./cmd/ /src/Team-GhostLand/Grinch/cmd
COPY ./trans/ /src/Team-GhostLand/Grinch/trans
COPY ./util/ /src/Team-GhostLand/Grinch/util
# COPY ./.git/ /src/Team-GhostLand/Grinch/.git
RUN ["pwd"]
RUN ["ls", "-al"]
RUN ["go", "build", "-o", "/bin/grinch"]

FROM alpine/git
WORKDIR /app
COPY --from=builder /bin/grinch /app
COPY ./scripts/ci.sh /app
COPY ./scripts/make_serverpack.sh /app

ENTRYPOINT ["/app/ci.sh"]