FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./
ADD assets ./assets
ADD views ./views

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /bookstore-webserver


EXPOSE 8080

#ENV bookstore_title="Generic Bookstore" #set this at runtime

# Run
CMD ["/bookstore-webserver"]

#az group create -g acr-rg -l canadacentral
#az acr create --resource-group acr-rg --name prephcpregistry --sku Standard
#az acr update --name prephcpregistry  --anonymous-pull-enabled
#az acr build --image prephcr/bookstore-webserver:v1 --resource-group acr-rg --registry prephcpregistry  --file Dockerfile . 