FROM golang:1.14-alpine
WORKDIR /src
COPY . .
RUN go build -o /app .

# vimagick/youtube-dl contains ffmpeg, youtube-dl.
FROM vimagick/youtube-dl
COPY --from=0 /app /app
ENTRYPOINT ["/app"]
