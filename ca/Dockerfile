# syntax=docker/dockerfile:1
FROM golang:1.18

WORKDIR /app

RUN curl -sSL https://github.com/brandonweeks/certificates/archive/9cbb770f09d0375cb621feea71db98da443fe7ca.tar.gz \
    | tar -v -C /app --strip-components 1 -xz
RUN make build

COPY demo-config demo-config

EXPOSE 8080
CMD ["bin/step-ca", "demo-config/config/ca.json"]