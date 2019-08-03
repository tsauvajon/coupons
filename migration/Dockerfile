#  Use a tagged version instead of latest to avoid regressions
FROM golang:1.12.7-buster

WORKDIR /root/coupons/
COPY . .
RUN go get
CMD go run migration/migrate.go migration/creation.sql migration/seed.sql
