FROM golang:latest

WORKDIR /app

COPY . .

RUN  go build -o exoplanet ./cmd

EXPOSE 8080

CMD [ "./exoplanet" ]