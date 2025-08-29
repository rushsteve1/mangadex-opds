FROM golang:alpine AS build
    WORKDIR /app
    COPY . .
    RUN go mod download
    RUN CGO_ENABLED=0 go build .

FROM alpine:latest
    USER nobody:nobody
    COPY --from=build /app/mangadex-opds /bin/mangadex-opds
    EXPOSE 4444
    ENTRYPOINT ["/bin/mangadex-opds"]