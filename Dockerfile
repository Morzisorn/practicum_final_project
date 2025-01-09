FROM golang:1.23.1

RUN apt-get update && apt-get install -y gcc

WORKDIR /app

# Копируем и устанавливаем зависимости
COPY go.mod go.sum ./ 
RUN go mod download

COPY . .

WORKDIR /app/config

ENV TODO_PORT=7540
ENV TODO_DBFILE=scheduler.db
ENV TODO_PASSWORD=1234
ENV CGO_ENABLED=0 
ENV GOOS=linux 
ENV GOARCH=amd64

WORKDIR /app/cmd/practicum_final

RUN CGO_ENABLED=${CGO_ENABLED} GOOS=${GOOS} GOARCH=${GOARCH} go build -o /practicum-final

EXPOSE ${TODO_PORT}

CMD ["/practicum-final"]

WORKDIR /app