FROM golang as builder

WORKDIR /app
COPY main.go .

RUN apt-get update \
    && apt-get install xz-utils \
    && rm -rf /var/lib/apt/lists/*
RUN CGO_ENABLED=0 go build main.go

RUN mkdir ffmpeg-download \
    && cd ffmpeg-download \
    && curl -o /app/ffmpeg-download/ffmpeg.tar.xz -L https://johnvansickle.com/ffmpeg/builds/ffmpeg-git-amd64-static.tar.xz \
    && tar xvf /app/ffmpeg-download/ffmpeg.tar.xz --strip-components 1 \
    && cp ffmpeg ../ \
    && rm -rf /app/ffmpeg-download/

FROM alpine
RUN addgroup -S nonroot && adduser -S -g nonroot nonroot
WORKDIR /app/
RUN mkdir -p /app && chown -R nonroot:nonroot /app
USER nonroot
COPY --from=builder --chown=nonroot /app/main /app
COPY --from=builder --chown=nonroot /app/ffmpeg /app/ffmpeg
COPY --chown=nonroot small.ogv /app
RUN ls /app -laR
