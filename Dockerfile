# Use the official image as a parent image.
FROM golang

# Set the working directory.
WORKDIR app

# Copy the file from your host to your current location.
COPY monitor .

# Run the command inside your image filesystem.
RUN go get github.com/lib/pq
RUN go get github.com/sergi/go-diff/diffmatchpatch
RUN go build .

# Add metadata to the image to describe which port the container is listening on at runtime.
EXPOSE 8080

# Run the specified command within the container.
CMD [ "./monitor" ]
