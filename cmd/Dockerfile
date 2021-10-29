# Собираем в гошке
FROM golang:1.16.8 as build

ENV BIN_FILE /opt/fibonacci-app
ENV CODE_DIR /go/src

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
RUN CGO_ENABLED=0 go build -o ${BIN_FILE} cmd/*.go

# На выходе тонкий образ
FROM alpine:3.9

LABEL SERVICE="fibonacci"
LABEL MAINTAINERS="kozirevss@gmail.com"

ENV BIN_FILE "/opt/fibonacci-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/fibonacci/config.yml
COPY ./configs/config.yml ${CONFIG_FILE}

RUN mkdir -p /var/log

CMD LOG_FILE=/var/log/app.log ${BIN_FILE} -config ${CONFIG_FILE}
