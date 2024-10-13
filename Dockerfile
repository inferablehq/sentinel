FROM golang:1.23.1-alpine

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .
EXPOSE 8080

# Set up environment variables
ENV STRATEGY=minimal

# !IMPORTANT: DATA_DIR must be a volume for data persistence accross restarts
ARG DATA_DIR=/data
ENV DATA_DIR=${DATA_DIR}

CMD ["./main"]
