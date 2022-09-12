FROM golang:alpine AS builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

# Fetch dependencies.
# Using go get.
RUN go get -d -v

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" .
# Output from build will be the same name as go.mod module name 
# module github.com/oscaralmgren/hackernews - hackernews


# STEP 2 build a small image
FROM busybox

WORKDIR /app

# Copy our static executable. 
COPY --from=builder /app/hackernews /usr/bin/

# Run the hello binary.
ENTRYPOINT ["hackernews"]