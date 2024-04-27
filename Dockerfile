# -----------------------------------------------------------------------------
#                                    Builder
# -----------------------------------------------------------------------------
FROM golang:1.20-alpine as builder

RUN apk update \
    && apk upgrade

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /app

# Import the code from the context.
COPY . ./

RUN make build-static

# -----------------------------------------------------------------------------
#                                Production Image
# -----------------------------------------------------------------------------
FROM alpine:3.7

WORKDIR /
RUN mkdir app
# Get all the executables
COPY --from=builder /app/link-identity-api /

# Create symlink to the application for this container.
ARG APP_NAME
RUN ln -s "${APP_NAME}" /main
CMD ["/main"]
