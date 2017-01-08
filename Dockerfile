FROM floreks/kubepi-base

COPY build/led-max7219-client-linux-arm-6 /usr/bin/led-max7219-client

ENTRYPOINT ["/usr/bin/led-max7219-client"]
