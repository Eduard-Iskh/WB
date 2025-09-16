FROM golang:1.25

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download

COPY . .
#изенить(хз что тут пока)
RUN go build -v -o /usr/local/bin/app ./...


#хз, что пока нужно 
CMD ["app"]