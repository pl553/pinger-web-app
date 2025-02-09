FROM debian:latest

RUN apt update
RUN apt install -y iptables

COPY toggle.sh /toggle.sh
RUN chmod +x /toggle.sh

CMD "/toggle.sh"


