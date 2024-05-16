# stage 1: build
#FROM golang:1.14 as build
#LABEL stage=intermediate
#WORKDIR /app
#COPY . .
#RUN make build

# stage 2: scratch
FROM alpine:latest as scratch
RUN apk --no-cache add ca-certificates
WORKDIR /opt/ftpdts
COPY bin/ftpdts ./ftpdts
COPY *.ini ./config/
COPY tmpl/*.tmpl /tmpl/
VOLUME /opt/ftpdts/data
RUN ln -s config/ftpdts.ini ftpdts.ini
CMD ["./ftpdts"]
