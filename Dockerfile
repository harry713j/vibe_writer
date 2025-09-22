# Build stage
FROM golang:alpine AS Builder

# Install git (needed for fetching private modules sometimes)
RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build static binary
RUN GOOS=linux GOARCH=amd64 go build -o vibewriter ./cmd/server

# Final stage
FROM alpine:latest

COPY --from=Builder /app/vibewriter .

CMD [ "./vibewriter" ]