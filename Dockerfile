# use official Golang image
FROM golang:1.16.3-alpine3.13

# set working directory
WORKDIR /app

# Copy the source code
COPY . . 

# Download and install the dependencies
RUN go get -d -v ./...

# Build the Go app
RUN go build -o main .

#EXPOSE the port
EXPOSE 8001

# Run the executable
CMD ["./main"]