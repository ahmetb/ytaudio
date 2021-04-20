FROM golang:1.15-alpine
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -tags netgo -o /app .

# vimagick/youtube-dl contains ffmpeg, youtube-dl.
FROM debian
RUN apt-get -qqy update && apt-get -qqy install curl ffmpeg python
RUN curl -L https://yt-dl.org/downloads/latest/youtube-dl -o /usr/local/bin/youtube-dl && \
        chmod a+rx /usr/local/bin/youtube-dl

COPY --from=0 /app /app
ENTRYPOINT ["/app"]
