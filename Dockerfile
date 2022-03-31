FROM  golang:1.17

COPY go.mod ./
COPY go.sum ./
COPY *.go ./

EXPOSE  2222

CMD ["./honeypot"]
