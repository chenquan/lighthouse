/*
 *    Copyright 2021 chenquan
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package server

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/yunqi/lighthouse/config"
	"github.com/yunqi/lighthouse/internal/goroutine"
	"github.com/yunqi/lighthouse/internal/persistence"
	"github.com/yunqi/lighthouse/internal/persistence/session"
	"github.com/yunqi/lighthouse/internal/persistence/subscription"
	"github.com/yunqi/lighthouse/internal/xlog"
	"github.com/yunqi/lighthouse/internal/xtrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"net"
	"time"
)

type (
	Server interface {
		Stop(ctx context.Context) error
		Run() error
	}
	Option func(server *Options)

	Options struct {
		tcpListen       string
		websocketListen string
		persistence     *config.Persistence
	}
	server struct {
		tcpListen         string
		websocketListen   string
		tcpListener       net.Listener //tcp listeners
		websocketListener *websocket.Conn
		sessionStore      session.Store
		subscriptionStore subscription.Store
		log               *xlog.Log
		tracer            trace.Tracer
	}
)

func WithTcpListen(tcpListen string) Option {
	return func(opts *Options) {
		opts.tcpListen = tcpListen
	}
}
func WithPersistence(persistence *config.Persistence) Option {
	return func(opts *Options) {
		opts.persistence = persistence
	}
}

func WithWebsocketListen(websocketListen string) Option {
	return func(opts *Options) {
		opts.websocketListen = websocketListen
	}
}

func NewServer(opts ...Option) *server {
	options := loadServerOptions(opts...)
	s := &server{}
	s.init(options)
	s.log = xlog.LoggerModule("server")
	return s
}
func loadServerOptions(opts ...Option) *Options {
	options := new(Options)
	for _, opt := range opts {
		opt(options)
	}
	if options.tcpListen == "" {
		options.tcpListen = ":1883"
	}
	return options
}

func (s *server) ServeTCP() {
	//propagator := otel.GetTextMapPropagator()
	s.tracer = otel.GetTracerProvider().Tracer(xtrace.Name)

	defer func() {
		err := s.tcpListener.Close()
		if err != nil {
			s.log.Error("tcpListener close", zap.Error(err))
		}
	}()
	var tempDelay time.Duration

	for {
		accept, err := s.tcpListener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				time.Sleep(tempDelay)
				continue
			}
			return
		}
		// 创建一个客户端连接

		c := newClient(s, accept)
		// 监听该连接
		goroutine.Go(func() {
			c.listen()
		})

	}
}

func (s *server) init(opts *Options) {
	s.tcpListen = opts.tcpListen
	s.websocketListen = opts.websocketListen
	s.log = xlog.LoggerModule("server")

	// session store
	sessionStore, ok := persistence.GetSessionStore(opts.persistence.Session.Type)
	if !ok {
		s.log.Panic("invalid session store")
	}

	if store, err := sessionStore(&opts.persistence.Session); err != nil {
		s.log.Panic("session store", zap.Error(err))
	} else {
		s.sessionStore = store
		s.log.Info("session store", zap.String("type", opts.persistence.Session.Type))
	}

	// subscriptionStore store
	subscriptionStoreFunc, ok := persistence.GetSubscriptionStore(opts.persistence.Subscription.Type)
	if !ok {
		s.log.Panic("invalid subscriptionStore store")
	}

	if subscriptionStore, err := subscriptionStoreFunc(&opts.persistence.Subscription); err != nil {
		s.log.Panic("subscriptionStore store", zap.Error(err))
	} else {
		s.subscriptionStore = subscriptionStore
		s.log.Info("subscriptionStore store", zap.String("type", opts.persistence.Session.Type))
	}

	ln, err := net.Listen("tcp", s.tcpListen)
	if err != nil {
		s.log.Panic("start tcp error", zap.String("tcp", s.tcpListen), zap.Error(err))
	}
	s.log.Info("start tcp", zap.String("TCP", s.tcpListen))
	s.tcpListener = ln

}
