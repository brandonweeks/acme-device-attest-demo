# syntax=docker/dockerfile:1
FROM golang:1.18

WORKDIR /app

RUN curl -sSL https://github.com/brandonweeks/certificates/archive/ea80e562fa058fb47931fa84bfa75de5a37f91b3.tar.gz \
    | tar -v -C /app --strip-components 1 -xz
RUN make build

COPY demo-config demo-config

EXPOSE 8080
CMD ["bin/step-ca", "demo-config/config/ca.json"]