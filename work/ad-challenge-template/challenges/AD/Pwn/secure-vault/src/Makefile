CC ?= gcc
CFLAGS ?= -O2 -fstack-protector

all: server

server: server.c
	$(CC) $(CFLAGS) server.c -o server

clean:
	rm -f server
