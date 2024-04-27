# -----------------------------------------------------------------------------
#                                    Builder
# -----------------------------------------------------------------------------
FROM golang:1.21-alpine as builder

RUN apk update \
    && apk upgrade \
    && apk add --no-cache make git

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

COPY app.env .

# Create symlink to the application for this container.
ARG APP_NAME=link-identity-api
RUN ln -s "${APP_NAME}" /main
CMD ["/main"]
