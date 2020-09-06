# Use the official image as a parent image.
FROM golang AS build

# Set the working directory.
WORKDIR /monitor/

# Copy the file from your host to your current location.
COPY ./monitor/*.go /monitor/

# Run the command inside your image filesystem.
RUN go get github.com/lib/pq
RUN go get github.com/sergi/go-diff/diffmatchpatch
RUN CGO_ENABLED=0 go build -o /bin/monitor

# Add metadata to the image to describe which port the container is listening on at runtime.
EXPOSE 8080:8080

ENTRYPOINT /bin/monitor

