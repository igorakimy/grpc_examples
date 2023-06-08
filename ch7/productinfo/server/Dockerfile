# Многоэтапная сборка

# Этап I:
# Для сборки этой программы требуется только язык Go и Alpine Linux
FROM golang AS build
ENV location app

ADD ./server ${location}/server
ADD ./protobuf ${location}/protobuf

COPY go.mod go.sum ${location}/

WORKDIR ${location}

# Загружаем все зависимости
RUN go get -d ./...
# Устанавливаем все пакеты
RUN go install ./...

WORKDIR server

# Собираем серверное приложение
RUN CGO_ENABLED=0 go build -o /bin/grpc-productinfo-server

# Этап II:
# Собранные программы на Go - самодостаточные исполняемые файлы
FROM scratch
# Копируем двоичный фай, собранный на предыдущем этапе, в новое место
COPY --from=build /bin/grpc-productinfo-server /bin/grpc-productinfo-server
ENTRYPOINT ["/bin/grpc-productinfo-server"]
EXPOSE 50051