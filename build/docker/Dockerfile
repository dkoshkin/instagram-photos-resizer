FROM dkoshkin/instagram-photos-resizer-gomod as builder

WORKDIR /src/github.com/dkoshkin/instagram-photos-resizer

# Copy local code to the container image.
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o cli ./cmd/cli/

# production container
FROM alpine
RUN apk add --no-cache ca-certificates

# Copy the binary to the production image from the builder stage.
COPY --from=builder /src/github.com/dkoshkin/instagram-photos-resizer/cli /cli

# Run the web service on container startup.
ENTRYPOINT ["/cli"]