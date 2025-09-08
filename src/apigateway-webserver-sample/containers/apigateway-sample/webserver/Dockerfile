# syntax=docker/dockerfile:1

FROM golang:1.23 AS builder
WORKDIR /apigateway

# Add source code of apigateway sample
COPY . .

# Build artifacts
RUN make -C apigateway-webserver -s clean && \
	make -C apigateway-webserver -s build

# Start from scratch (minimal base image)
FROM scratch
# Or use ubuntu instead
#FROM ubuntu:24.04

# Copy in files that we have built
COPY --from=builder /apigateway/apigateway-webserver/bin /bin

# Define the command to run when the container starts
ENTRYPOINT  ["/bin/apigateway-webserver"]
