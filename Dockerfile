# Build stage
FROM golang:alpine as builder

COPY httpenv.go /go
RUN go build httpenv.go

# Final stage
FROM alpine

# Create a non-root user and group to run the application
RUN addgroup -g 1000 httpenv \
    && adduser -u 1000 -G httpenv -D httpenv

# Copy the compiled binary from the builder stage and set permissions
COPY --from=builder --chown=httpenv:httpenv /go/httpenv /httpenv

# Install curl for debugging purposes
RUN apk add --no-cache curl

# Expose the port the application will listen on
EXPOSE 8888

CMD ["/httpenv"]
