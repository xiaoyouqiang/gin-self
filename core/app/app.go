package app

import (
	"context"
	"fmt"
	"gin-self/app"
	"gin-self/extend/config"
	"gin-self/extend/self_db"
	"gin-self/extend/self_loger"
	"gin-self/extend/self_redis"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	ctx    context.Context
	server *http.Server
	cancel func()
}

func New() *App {
	ctx, cancel := context.WithCancel(context.Background())
	return &App{
		ctx: ctx,
		server: &http.Server{
			Addr:         fmt.Sprintf(":%d", config.Get("server", "port").MustInt(8080)),
			Handler:      app.IncludeRoute(),
			ReadTimeout:  time.Duration(config.Get("server", "read_timeout").MustInt(60)) * time.Second,
			WriteTimeout: time.Duration(config.Get("server", "write_timeout").MustInt(60)) * time.Second,
		},
		cancel: cancel,
	}
}

func (a *App) Run() {
	group, _ := errgroup.WithContext(context.Background())

	group.Go(func() (err error) {
		err = a.start()
		return
	})

	group.Go(func() (err error) {
		a.registerExitSignal()
		return
	})

	group.Go(func() (err error) {
		a.shutdown()
		return
	})

	log.Printf("%s: errgroup exit %v\n", time.Now().Format("2006-01-02 15:04:05"), group.Wait())
}

//开始监听服务
func (a *App) start() error {
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

//注册中断退出信号量
func (a *App) registerExitSignal() {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	log.Printf("%s: exit by signal %v\n", time.Now().Format("2006-01-02 15:04:05"), <-quit)

	//trigger shutdown
	a.cancel()
}

//平滑退出服务
func (a *App) shutdown() {
	<-a.ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

func (a *App) InitMysql() {
	self_db.Open("test")
}

func (a *App) InitRedis() {
	self_redis.Open("master")
}

func (a *App) InitLogger(logType self_loger.LogType) {
	self_loger.GetInstance().SetLogger(logType)
}

