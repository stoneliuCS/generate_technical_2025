FROM golang:latest

# Install Task
RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d

# Copy the src into the container
COPY ./challenge ./challenge
COPY ./openapi.json ./openapi.json

# Set the container's working directory
WORKDIR /go/challenge/

# Build and run the backend server
RUN task build --force
CMD ["./main_server"]
