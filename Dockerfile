# Use a tagged version instead of latest to avoid regressions
FROM golang:1.12.7-buster as builder

# Copy the local package files to the container's workspace.
WORKDIR /root/coupons/
COPY . .

# Set a location to install dependencies
ENV GOBIN=/go/bin/

# Build the app for alpine, executable name: coupons
RUN go get
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o coupons

# Use alpine for faster pulls and a smaller image
# Use a multi-stage Dockerfile to avoid keeping anything unecessary in the final image
FROM golang:1.12.7-alpine

WORKDIR /root/

# Copy the executable from the previous stage
COPY --from=builder /root/coupons/coupons /root

# Run on port 4000
ENV PORT=4000

# Run the app when the container starts
ENTRYPOINT /root/coupons

EXPOSE 4000