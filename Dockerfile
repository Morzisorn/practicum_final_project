FROM golang:1.23.1

RUN apt-get update && apt-get install -y gcc

WORKDIR /app

# Копируем и устанавливаем зависимости
COPY go.mod go.sum ./ 
RUN go mod download

COPY . .

ENV TODO_PORT=7540
ENV TODO_DBFILE=scheduler.db
ENV TODO_PASSWORD=1234

RUN go mod tidy 

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /practicum-final

EXPOSE 7540

CMD ["/practicum-final"]