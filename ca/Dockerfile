# syntax=docker/dockerfile:1
FROM golang:1.18

WORKDIR /app

RUN curl -sSL https://github.com/brandonweeks/certificates/archive/a663773a133d12616538fcf86b87ae8d47925c51.tar.gz \
    | tar -v -C /app --strip-components 1 -xz
RUN make build

COPY demo-config demo-config

EXPOSE 8080
CMD ["bin/step-ca", "demo-config/config/ca.json"]