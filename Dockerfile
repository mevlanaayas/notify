# We specify the base image we need for our
# go application
FROM golang:1.11.11-alpine3.9
# install git to get dependencies
RUN apk add --no-cache git
# get dependencies
RUN go get -u github.com/gin-gonic/gin
RUN go get -u github.com/joho/godotenv
RUN go get -u github.com/sendgrid/sendgrid-go
RUN go get github.com/getsentry/sentry-go
# We create an /app directory within our
# image that will hold our application source
# files
RUN mkdir /app
# We copy everything in the root directory
# into our /app directory
ADD . /app
# We specify that we now wish to execute 
# any further commands inside our /app
# directory
WORKDIR /app
# we run go build to compile the binary
# executable of our Go program
RUN go build -o main .
# expose port
EXPOSE 7070
# Our start command which kicks off
# our newly created binary executable
CMD ["/app/main"]