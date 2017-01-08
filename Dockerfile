FROM floreks/kubepi-base

COPY build/led-max7219-client-linux-arm-6 /usr/bin/led-max7219-client

CMD ["/usr/bin/led-max7219-client"]

EXPOSE 4000