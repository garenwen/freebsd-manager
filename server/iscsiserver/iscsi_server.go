package iscsiserver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/trace"
	"syscall"
	"time"

	"github.com/garenwen/freebsd-manager/server"

	"github.com/garenwen/freebsd-manager/handle"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"golang.org/x/sync/errgroup"
)

type IscsiServer struct {
	context       context.Context
	shutdownFn    context.CancelFunc
	childRoutines *errgroup.Group

	gin     *gin.Engine
	httpSrv *http.Server
}

func NewIscsiServer() server.Server {
	rootCtx, shutdownFn := context.WithCancel(context.Background())
	childRoutines, childCtx := errgroup.WithContext(rootCtx)

	return &IscsiServer{
		context:       childCtx,
		shutdownFn:    shutdownFn,
		childRoutines: childRoutines,
	}
}

func (zs *IscsiServer) Start() {

	go ListenToSystemSignals(zs)
	glog.Infoln("Gin Server staring")
	zs.gin = zs.newGin()

	zs.httpSrv = &http.Server{
		Addr:           "0.0.0.0:8880",
		Handler:        zs.gin,
		ReadTimeout:    70 * time.Second,
		WriteTimeout:   70 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := zs.httpSrv.ListenAndServe(); err != nil {

		glog.Errorln("server was shutdown gracefully")
		zs.Shutdown(1, "Startup failed")
		return

	}
}

func (zs *IscsiServer) newGin() *gin.Engine {

	route := gin.Default()

	// accounts := gin.Accounts{"measure": "measure"}
	// route.Use(gin.BasicAuth(accounts))

	route.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	iscsiHandler := handle.NewIscsiHandler()
	v1 := route.Group("api/v1")
	{
		v1.POST("/create_iscsi", iscsiHandler.HandleCreateIscsi)
	}

	return route
}

func (zs *IscsiServer) Shutdown(code int, reason string) {
	glog.Infoln("Shutdown started", "code", code, "reason", reason)

	err := zs.httpSrv.Shutdown(zs.context)
	if err != nil {
		glog.Fatal("Failed to shutdown server", "error", err)
	}
	zs.shutdownFn()
	err = zs.childRoutines.Wait()

	glog.Error("Shutdown completed", "reason", err)

	os.Exit(code)
}

var exitChan = make(chan int)

func ListenToSystemSignals(server server.Server) {
	signalChan := make(chan os.Signal, 1)
	ignoreChan := make(chan os.Signal, 1)
	code := 0

	signal.Notify(ignoreChan, syscall.SIGHUP)
	signal.Notify(signalChan, os.Interrupt, os.Kill, syscall.SIGTERM)

	select {
	case sig := <-signalChan:
		// Stops trace if profiling has been enabled
		trace.Stop()
		server.Shutdown(0, fmt.Sprintf("system signal: %s", sig))
	case code = <-exitChan:
		server.Shutdown(code, "startup error")
	}
}
