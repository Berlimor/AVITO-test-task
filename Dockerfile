FROM golang:1.21

WORKDIR /AVITO-test-task

EXPOSE 8000

COPY . /AVITO-test-task

RUN go mod download

RUN go build -v

CMD [ "./main" ]