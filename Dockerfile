# syntax=docker/dockerfile:1

# Builds the skillstore CLI as a static binary and packages it into a
# minimal, non-root runtime image. The image only contains the CLI — it
# does not bundle any skills/ directory; mount your own at run time.
#
# Build:
#   docker build -t skillstore .
#
# Run (skills/ and public/ are relative to your current directory):
#   docker run --rm \
#     -v "$(pwd)/skills:/skills:ro" \
#     -v "$(pwd)/public:/public" \
#     skillstore --skill-dir /skills --output-dir /public

FROM golang:1.27-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/skillstore .

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /out/skillstore /usr/local/bin/skillstore

ENTRYPOINT ["/usr/local/bin/skillstore"]
CMD ["--help"]
