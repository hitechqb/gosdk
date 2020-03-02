package sdkplugin

import (
	"context"
	"fmt"
	"github.com/hitech/gosdk"
	"github.com/hitech/gosdk/sdkcm"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type simpleHttpServer struct {
	*plugin
	logger *logrus.Logger
	server *http.Server
	cf     sdkcm.SDKConfig
}

func NewSimpleHttpServer(options ...PluginOption) *simpleHttpServer {
	logger := sdkcm.GetLogger()
	cf := sdkcm.GetConfig()

	s := &simpleHttpServer{
		plugin: &plugin{},
		logger: logger,
		cf:     cf,
	}

	for _, o := range options {
		o(s)
	}

	return s
}

func (s *simpleHttpServer) Run(sdk.Service) error {
	srv := &http.Server{
		// Handler:      ,
		Addr:         fmt.Sprintf(`:%s`, s.cf.AppPort()),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start Server
	go func() {
		s.logger.Infof(`🎉 Server is listening on port %s`, s.cf.AppPort())
		if err := srv.ListenAndServe(); err != nil {
			s.logger.Fatalln(err)
		}
	}()

	// Graceful Shutdown
	s.waitForShutdown(srv)

	return nil
}

func (s *simpleHttpServer) Stop() error {
	return nil
}

func (s *simpleHttpServer) GetServer() interface{} {
	return s.server
}

func (s *simpleHttpServer) waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	s.logger.Infof("Shutting down")
	os.Exit(0)
}
