FROM golang:1.20
WORKDIR /app
RUN go install github.com/cosmtrek/air@v1.29.0
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .
EXPOSE 8080
CMD ["./main"] # 'air' is not used in the production environment
