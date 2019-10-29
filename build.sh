#!/bin/bash

docker build -t jeremyje/ffmpeg-dns-localhost-repro:latest .
docker push jeremyje/ffmpeg-dns-localhost-repro:latest
