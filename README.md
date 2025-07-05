# a giant canvas

I stumbled upon this idea not too long ago - a (more or less) global canvas that anyone
could draw on, _whatever*_ they like, and in real time, maybe even with friends if they
want to.

## technical nerdy stuff

The project consists of a server and a client. A client is a web browser, the server is,
well, a server, written in Go. The client connects to the server over WebSockets, and
then they exchange updates in JSON. It's not too complicated, but I wanted to test myself
if I can build something using a protocol that I have never really used.

The web frontend is built with Preact, so it's really, really lightweight and also not
very complicated.
