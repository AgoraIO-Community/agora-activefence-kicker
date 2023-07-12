# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:alpine
RUN apk add git ca-certificates --update

# fetch dependencies github
# RUN go get -u github.com/gin-gonic/gin

ADD . /go/src/github.com/AgoraIO-Community/agora-activefence-kicker

# # fetch dependencies from github (Gin)
# RUN go install github.com/gin-gonic/gin@latest
# # RUN go install github.com/AgoraIO-Community/agora-activefence-kicker
# ADD . /go/src/github.com/AgoraIO-Community/agora-activefence-kicker

ARG APP_ID
ARG CUSTOMER_KEY
ARG CUSTOMER_SECRET
ARG SERVER_PORT
ENV APP_ID $APP_ID
ENV CUSTOMER_KEY $CUSTOMER_KEY
ENV CUSTOMER_SECRET $CUSTOMER_SECRET
ENV SERVER_PORT $SERVER_PORT

# move to the working directory
WORKDIR $GOPATH/src/github.com/AgoraIO-Community/agora-activefence-kicker
# Build the kicker server command inside the container.
RUN go build -o agora-activefence-kicker -v cmd/main.go
# RUN go run main.go
# Run the kicker server by default when the container starts.
ENTRYPOINT ./agora-activefence-kicker

# Document that the service listens on port $SERVER_PORT.
EXPOSE $SERVER_PORT