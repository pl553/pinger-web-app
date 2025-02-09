FROM debian:latest

RUN apt update
RUN apt install -y iptables

CMD "sh" "-c" "iptables -A INPUT -p icmp --icmp-type echo-request -j DROP && tail -f /dev/null"

