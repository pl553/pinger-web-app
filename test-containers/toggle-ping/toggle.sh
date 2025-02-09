#!/bin/bash

# allows and disallows ICMP echo requests every 10 seconds

while true; do

  iptables -I INPUT -p icmp --icmp-type echo-request -j DROP  
  sleep 10

  iptables -D INPUT -p icmp --icmp-type echo-request -j DROP
  sleep 10
done


