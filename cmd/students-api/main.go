package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/patelaryan0914/students-api/internal/config"
	student "github.com/patelaryan0914/students-api/internal/http/handlers"
	"github.com/patelaryan0914/students-api/internal/storage/sqlite"
)

func main() {
	//load config
	cfg:=config.MustLoad()
	//logger setup
	//database setup
	storage,err := sqlite.New(cfg)
	if err!=nil{
		log.Fatal("Database Connection error ",err)
	}
	slog.Info("Storage Initialized",slog.String("env",cfg.Env),slog.String("version","1.0.0"))
	//setup router
	router:=http.NewServeMux()
	router.HandleFunc("POST /api/students",student.New(storage))
	router.HandleFunc("GET /api/students/{id}",student.GetByID(storage))
	//setup server
	server := http.Server {
		Addr : cfg.Addr,
		Handler:router ,
	}
	slog.Info("Server Started",slog.String("address",cfg.Addr))
	done :=make(chan os.Signal,1)
	signal.Notify(done,os.Interrupt,syscall.SIGINT,syscall.SIGTERM)
	go func(){
		err := server.ListenAndServe()
		if err!= nil{
			log.Fatal("Failed to Start Server")
		}
	}()

	<-done

	slog.Info("Shutting Down the server")

	ctx,cancel:=context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()
	err=server.Shutdown(ctx)
	if err!=nil{
		slog.Error("Failed to Shutdown server",slog.String("error",err.Error()))
	}

	slog.Info("Server Shutdown Successfully")


}