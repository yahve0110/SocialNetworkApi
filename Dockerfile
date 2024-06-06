# Используем официальный образ Go
FROM golang:1.21.7

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы бэкенда
COPY . .

# Запускаем go mod tidy для обновления зависимостей
RUN go mod tidy

# Устанавливаем утилиту migrate
RUN go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Добавляем папку с установленными Go утилитами в PATH
ENV PATH="/go/bin:${PATH}"

# Собираем бинарник бэкенда и выполняем миграции
RUN make all

# Экспонируем порт 8080
EXPOSE 8080

# Указываем команду для запуска бэкенда
CMD ["./api"]
