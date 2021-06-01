# grpc_training

初始化过程

```go
	s := grpc.NewServer(grpc.UnaryInterceptor(interceptors.OrderUnaryServerInterceptor1),
		grpc.ChainUnaryInterceptor(interceptors.OrderUnaryServerInterceptor2, interceptors.OrderUnaryServerInterceptor3))
```

gopath\pkg\mod\google.golang.org\grpc@v1.38.0\server.go=>561

```go
// NewServer creates a gRPC server which has no service registered and has not
// started to accept requests yet.
func NewServer(opt ...ServerOption) *Server {
	opts := defaultServerOptions
	for _, o := range opt {
		o.apply(&opts)
	}
	s := &Server{
		lis:      make(map[net.Listener]bool),
		opts:     opts,
		conns:    make(map[string]map[transport.ServerTransport]bool),
		services: make(map[string]*serviceInfo),
		quit:     grpcsync.NewEvent(),
		done:     grpcsync.NewEvent(),
		czData:   new(channelzData),
	}
	chainUnaryServerInterceptors(s)
	chainStreamServerInterceptors(s)
	s.cv = sync.NewCond(&s.mu)
	if EnableTracing {
		_, file, line, _ := runtime.Caller(1)
		s.events = trace.NewEventLog("grpc.Server", fmt.Sprintf("%s:%d", file, line))
	}

	if s.opts.numServerWorkers > 0 {
		s.initServerWorkers()
	}

	if channelz.IsOn() {
		s.channelzID = channelz.RegisterServer(&channelzServer{s}, "")
	}
	return s
}
```



defaultServerOptions 继承\pkg\mod\google.golang.org\grpc@v1.38.0\server.go=>143 serverOptions

```go
看下serverOptions结构
=============
type serverOptions struct {
	creds                 credentials.TransportCredentials
	codec                 baseCodec
	cp                    Compressor
	dc                    Decompressor
	unaryInt              UnaryServerInterceptor		// 一元拦截器
	streamInt             StreamServerInterceptor		// 流拦截器
	chainUnaryInts        []UnaryServerInterceptor		// 一元拦截器链表
	chainStreamInts       []StreamServerInterceptor		// 流拦截器链表
	inTapHandle           tap.ServerInHandle
	statsHandler          stats.Handler
	maxConcurrentStreams  uint32
	maxReceiveMessageSize int
	maxSendMessageSize    int
	unknownStreamDesc     *StreamDesc
	keepaliveParams       keepalive.ServerParameters
	keepalivePolicy       keepalive.EnforcementPolicy
	initialWindowSize     int32
	initialConnWindowSize int32
	writeBufferSize       int
	readBufferSize        int
	connectionTimeout     time.Duration
	maxHeaderListSize     *uint32
	headerTableSize       *uint32
	numServerWorkers      uint32
}
```

chainUnaryServerInterceptors(s)和
chainStreamServerInterceptors(s)