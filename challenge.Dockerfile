FROM golang:latest

# Install Task
RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d

# Install bun and add it to the PATH env variable.
RUN apt-get update && apt-get install -y zip && apt-get clean && rm -rf /var/lib/apt/lists/*
RUN curl -fsSL https://bun.sh/install | bash
ENV PATH="/root/.bun/bin:$PATH"

# Copy the src into the container
COPY ./challenge ./challenge
COPY ./openapi.json ./openapi.json

# Set the container's working directory
WORKDIR /go/challenge/

# Build and run the backend server
CMD ["sh", "-c", "task build --force && task run"]
