# syntax=docker/dockerfile:1
FROM golang:1.18

WORKDIR /app

RUN curl -sSL https://github.com/brandonweeks/certificates/archive/274f6ccb4186958e3b1e8ad50d3acd2eca20235b.tar.gz \
    | tar -v -C /app --strip-components 1 -xz
RUN make build

COPY demo-config demo-config

EXPOSE 8080
CMD ["bin/step-ca", "demo-config/config/ca.json"]