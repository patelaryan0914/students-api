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
)

func main() {
	//load config
	cfg:=config.MustLoad()
	//logger setup
	//database setup
	//setup router
	router:=http.NewServeMux()
	router.HandleFunc("GET /",func(w http.ResponseWriter, r *http.Request) {
		 w.Write([]byte("Welcome to Students-api"))
	})
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
	err:=server.Shutdown(ctx)
	if err!=nil{
		slog.Error("Failed to Shutdown server",slog.String("error",err.Error()))
	}

	slog.Info("Server Shutdown Successfully")


}