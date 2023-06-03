FROM alpine:latest

RUN apk add -v build-base go ca-certificates
RUN apk add --no-cache \
    unzip \
    # this is needed only if you want to use scp to copy later your pb_data locally
    openssh

# Copy your custom PocketBase and build
COPY ./pocketbase-custom /pb
WORKDIR /pb

RUN go build
WORKDIR /

EXPOSE 8080

# start PocketBase
CMD ["/pb/pocketbase-custom", "serve", "--http=0.0.0.0:8080"]
