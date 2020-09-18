# Build the Go Binary.
FROM golang:1.14.4 as build-customer-service
ENV CGO_ENABLED 0
ARG PACKAGE_NAME
ARG PACKAGE_PREFIX

# Create a location in the container for the source code. Using the
# default GOPATH location.
RUN mkdir -p /service

# Copy the source code into the container.
WORKDIR /service
COPY . .

# Build the service binary. We are doing this last since this will be different
# every time we run through this process.
WORKDIR /service/app/customer-service/
RUN go build


# Run the Go Binary in Alpine.
FROM alpine:3.7
ARG BUILD_DATE
ARG PACKAGE_NAME
ARG PACKAGE_PREFIX
COPY --from=build-customer-service /service/app/customer-service/customer-service /app/main
COPY --from=build-customer-service /service/business/templates /app/business/templates
COPY --from=build-customer-service /service/start-customer-service.sh /app/start.sh

WORKDIR /app

CMD sh /app/start.sh

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="${PACKAGE_NAME}" \
      org.opencontainers.image.authors="Alex Lucaci <alexlucaci@usefulpython.com>" \
      org.opencontainers.image.vendor="Alex Lucaci"


