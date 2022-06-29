FROM golang as builder

WORKDIR /go/src/

RUN go get golang.org/x/net/publicsuffix
RUN go get github.com/spf13/pflag
RUN go get github.com/go-resty/resty

COPY src .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -a -installsuffix cgo -o skype-devops-bot-wrapper main/main.go



FROM alpine

RUN apk --no-cache --update add ca-certificates
RUN update-ca-certificates

COPY --from=builder /go/src/skype-devops-bot-wrapper /bin/skype-devops-bot-wrapper

CMD ["skype-devops-bot-wrapper"]