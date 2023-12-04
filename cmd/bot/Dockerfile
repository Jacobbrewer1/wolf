FROM golang:1.20.1

LABEL org.opencontainers.image.source='https://github.com/Jacobbrewer1/wolf'
LABEL org.opencontainers.image.description="A discord bot that can do a lot of things."
LABEL org.opencontainers.image.licenses='GNU General Public License v3.0'

WORKDIR /app/wolf

COPY ./bin/app /app/wolf

CMD ["./app"]