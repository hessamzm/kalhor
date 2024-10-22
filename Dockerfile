# مرحله 1: ساخت برنامه
FROM ubuntu:22.04 AS builder

# نصب ابزارهای لازم
RUN apt update && \
    apt install -y software-properties-common && \
    add-apt-repository -y ppa:longsleep/golang-backports && \
    apt update && \
    apt install -y golang-go curl

# تنظیمات محیط
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOARCH=amd64 \
    GOOS=linux

# دایرکتوری کاری
WORKDIR /app

# کپی فایل‌های go.mod و go.sum
COPY go.mod go.sum ./

# دانلود وابستگی‌ها
RUN GOPROXY=https://goproxy.cn go mod download

# کپی سایر فایل‌های پروژه
COPY . .

# ساخت برنامه
RUN go build -o main .

# مرحله 2: تصویر نهایی
FROM ubuntu:22.04

# تنظیم دایرکتوری کاری
WORKDIR /www/wwwroot/app.kalhorgold.ir

# کپی فایل اجرایی و سایر منابع
COPY --from=builder /app/main .
COPY --from=builder /app/public ./public/
COPY --from=builder /app/storage ./storage/
COPY --from=builder /app/resources ./resources/
COPY --from=builder /app/.env .

# اجرای برنامه
ENTRYPOINT ["./main"]
