FROM        prom/busybox:latest
MAINTAINER  Joel Hermanns <joel.hermanns@gmail.com>

COPY bin/server.linux                /bin/prom-timestream-adapter

EXPOSE     4000
ENTRYPOINT [ "/bin/prom-timestream-adapter" ]
