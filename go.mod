module tcpserver

go 1.21.3

replace tcpserver/types => ./types

replace tcpserver/helper => ./helper

require (
	tcpserver/helper v0.0.0-00010101000000-000000000000
	tcpserver/types v0.0.0-00010101000000-000000000000
)
