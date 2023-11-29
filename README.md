# Simple GO TCP Server to recieve and send MODBUS TCP Packets to CLICK PLCs

Sending works! yay


## How to build:
```
go build
```

## Usage:
### Server
```
    ./tcpserver [port]

    - port: Port number to listen to. Default: 502
```
### Client
```
    ./tcpserver [ip] [port] [ddos]

    - ip: ip to send packets to
    - port: port server is listening on
    - ddos: activated ddos, option is "ddos"
```
