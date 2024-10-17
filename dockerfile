from scratch

label org.opencontainers.image.authors="Vladislav Gardner <vladislavgardner@gmail.com>"

workdir service

copy ./bin/service_start.sh .env ./

cmd ["./service_start.sh"]
