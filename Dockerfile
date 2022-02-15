FROM golang:1.17.3-alpine as builder
LABEL maintainer "deflinhec <deflinhec@gmail.com>"

COPY . /workspace
WORKDIR /workspace
ENV GO111MODULE=on
RUN apk add bash make && \
    make all

FROM golang:1.17-alpine AS production
COPY --from=builder /go/bin /go/bin
RUN apk add bash
WORKDIR /workspace

ENTRYPOINT [ "/go/bin/gsx2json" ]
CMD [ "--port 8080" ]