# syntax=docker/dockerfile:1
FROM golang:1.18

WORKDIR /app

RUN curl -sSL https://github.com/brandonweeks/certificates/archive/3a437d034bf03a268af4695ebb85f36aa6a2d6d5.tar.gz \
    | tar -v -C /app --strip-components 1 -xz
RUN make build

COPY demo-config demo-config

EXPOSE 8080
CMD ["bin/step-ca", "demo-config/config/ca.json"]