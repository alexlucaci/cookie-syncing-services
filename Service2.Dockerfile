# Build the Go Binary.
FROM golang:1.14.4 as build-service-2
ENV CGO_ENABLED 0
ARG PACKAGE_NAME
ARG PACKAGE_PREFIX

# Create a location in the container for the source code. Using the
# default GOPATH location.
RUN mkdir -p /service

# Copy the source code into the container.
WORKDIR /service
COPY . .

# Build the admin tool so we can have it in the container. This should not change
# often so do this first.
WORKDIR /service/app/admin
RUN go build

# Build the service binary. We are doing this last since this will be different
# every time we run through this process.
WORKDIR /service/app/partner-service-2/
RUN go build


# Run the Go Binary in Alpine.
FROM alpine:3.7
ARG BUILD_DATE
ARG PACKAGE_NAME
ARG PACKAGE_PREFIX
COPY --from=build-service-2 /service/app/admin/admin /app/admin
COPY --from=build-service-2 /service/app/partner-service-2/partner-service-2 /app/main
COPY --from=build-service-2 /service/start-partner-service-2.sh /app/start.sh

WORKDIR /app

CMD sh /app/start.sh

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="${PACKAGE_NAME}" \
      org.opencontainers.image.authors="Alex Lucaci <alexlucaci@usefulpython.com>" \
      org.opencontainers.image.vendor="Alex Lucaci"


