FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./
ADD assets ./assets
ADD views ./views

RUN CGO_ENABLED=0 GOOS=linux go build -o /bookstore-webserver


EXPOSE 8080

# Run
CMD ["/bookstore-webserver"]


#az acr build --image prephcr/bookstore-webserver:v1 --resource-group myResourceGroup --registry prephcpregistry  --file Dockerfile . 