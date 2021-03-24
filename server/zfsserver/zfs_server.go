package zfsserver

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

type ZfsServer struct {
	context       context.Context
	shutdownFn    context.CancelFunc
	childRoutines *errgroup.Group

	gin     *gin.Engine
	httpSrv *http.Server
}

func NewZfsServer() server.Server {
	rootCtx, shutdownFn := context.WithCancel(context.Background())
	childRoutines, childCtx := errgroup.WithContext(rootCtx)

	return &ZfsServer{
		context:       childCtx,
		shutdownFn:    shutdownFn,
		childRoutines: childRoutines,
	}
}

func (zs *ZfsServer) Start() {

	go ListenToSystemSignals(zs)
	glog.Infoln("Gin Server staring")
	zs.gin = zs.newGin()

	zs.httpSrv = &http.Server{
		Addr:           "0.0.0.0:8870",
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

func (zs *ZfsServer) newGin() *gin.Engine {

	route := gin.Default()

	// accounts := gin.Accounts{"measure": "measure"}
	// route.Use(gin.BasicAuth(accounts))

	route.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	zfsHandler := handle.NewZfsHandler()
	v1 := route.Group("apis/storage/v1")
	{
		// 汇聚查询接口
		v1.POST("/create_volume", zfsHandler.HandleCreateVolume)
	}

	return route
}

func (zs *ZfsServer) Shutdown(code int, reason string) {
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
