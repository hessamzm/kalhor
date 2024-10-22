FROM ubuntu:22.04 AS builder

# نصب ابزارهای لازم
RUN apt-get update && \
    apt-get install -y curl golang git && \
    rm -rf /var/lib/apt/lists/*

# تنظیمات محیط
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOARCH=amd64 \
    GOOS=linux

# دایرکتوری کاری
WORKDIR /

# کپی فایل‌های go.mod و go.sum
COPY /go.mod .
COPY /go.sum .

# دانلود وابستگی‌ها
RUN GOPROXY=https://goproxy.cn go mod download

# کپی سایر فایل‌های پروژه
COPY / .

# ساخت برنامه
RUN go build -o main .

# مرحله دوم: تصویر نهایی
FROM ubuntu:22.04

# تنظیم دایرکتوری کاری
WORKDIR /www/wwwroot/app.kalhorgold.ir

# کپی فایل اجرایی و سایر منابع
COPY --from=builder /main .
COPY --from=builder /public ./public/
COPY --from=builder /storage ./storage/
COPY --from=builder /resources ./resources/
COPY --from=builder /.env .

# اجرای برنامه
ENTRYPOINT ["./main"]

