package udp

import (
	"net"

	"github.com/go-gost/gost/v3/pkg/common/util/udp"
	"github.com/go-gost/gost/v3/pkg/listener"
	"github.com/go-gost/gost/v3/pkg/logger"
	md "github.com/go-gost/gost/v3/pkg/metadata"
	"github.com/go-gost/gost/v3/pkg/registry"
	metrics "github.com/go-gost/metrics/wrapper"
)

func init() {
	registry.ListenerRegistry().Register("udp", NewListener)
}

type udpListener struct {
	ln      net.Listener
	logger  logger.Logger
	md      metadata
	options listener.Options
}

func NewListener(opts ...listener.Option) listener.Listener {
	options := listener.Options{}
	for _, opt := range opts {
		opt(&options)
	}
	return &udpListener{
		logger:  options.Logger,
		options: options,
	}
}

func (l *udpListener) Init(md md.Metadata) (err error) {
	if err = l.parseMetadata(md); err != nil {
		return
	}

	laddr, err := net.ResolveUDPAddr("udp", l.options.Addr)
	if err != nil {
		return
	}

	var conn net.PacketConn
	conn, err = net.ListenUDP("udp", laddr)
	if err != nil {
		return
	}
	conn = metrics.WrapPacketConn(l.options.Service, conn)

	l.ln = udp.NewListener(
		conn,
		laddr,
		l.md.backlog,
		l.md.readQueueSize, l.md.readBufferSize,
		l.md.ttl,
		l.logger)
	return
}

func (l *udpListener) Accept() (conn net.Conn, md md.Metadata, err error) {
	conn, err = l.ln.Accept()
	return
}

func (l *udpListener) Addr() net.Addr {
	return l.ln.Addr()
}

func (l *udpListener) Close() error {
	return l.ln.Close()
}
