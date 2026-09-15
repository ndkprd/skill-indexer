# syntax=docker/dockerfile:1

# Builds the skill-indexer CLI as a static binary and packages it into a
# minimal, non-root runtime image. The image only contains the CLI — it
# does not bundle any skills/ directory; mount your own at run time.
#
# Build:
#   docker build -t skill-indexer .
#
# Run (skills/ and public/ are relative to your current directory):
#   docker run --rm \
#     -v "$(pwd)/skills:/skills:ro" \
#     -v "$(pwd)/public:/public" \
#     skill-indexer --skill-dir /skills --output-dir /public

FROM golang:1.27-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/skill-indexer .

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /out/skill-indexer /usr/local/bin/skill-indexer

ENTRYPOINT ["/usr/local/bin/skill-indexer"]
CMD ["--help"]
