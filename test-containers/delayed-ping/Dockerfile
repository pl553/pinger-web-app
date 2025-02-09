FROM debian:latest

RUN apt update
RUN apt install -y iproute2

CMD "sh" "-c" "tc qdisc add dev eth0 root netem delay 250ms && tail -f /dev/null"

