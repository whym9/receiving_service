package receiver

type Receiver interface {
	StartServer(addr string)
	Receive(opts ...any)
}
