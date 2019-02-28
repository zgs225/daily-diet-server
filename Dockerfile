FROM golang:1.12

# Dependencies
RUN go get -u -v github.com/sirupsen/logrus

# Build
ENV APP_NAME daily-diet-server
ENV APP_SRC $GOPATH/src/github.com/zgs225/$APP_NAME
WORKDIR $APP_SRC
COPY . $APP_SRC
RUN go build -v -o $GOPATH/bin/serverd

CMD ["serverd"]
