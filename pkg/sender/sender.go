package sender

type Sender interface {
	StartServer(addr string, ch *chan []byte)
	Upload(opts ...any) ([]byte, error)
}
