package main

import (
	"flag"
	"fmt"
	"kalhor/bootstrap"
	"kalhor/utils"
	"os"
	"os/signal"
	"syscall"

	"github.com/goravel/framework/facades"
)

func main() {

	// اجرای برنامه
	utils.InitDebug()
	// تابعی برای نمایش لاگ در صورت فعال بودن فلگ debug
	flag.Parse()
	// بررسی فلگ
	if utils.KlDebug {
		utils.LogDebug("Debugging is enabled!")
	} else {
		fmt.Println("Running without debugging.")
	}

	// ادامه برنامه شما
	fmt.Println("Program is running")

	// This bootstraps the framework and gets it ready for use.
	bootstrap.Boot()

	// Create a channel to listen for OS signals
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start http server by facades.Route().
	go func() {
		if err := facades.Route().Run(); err != nil {
			facades.Log().Errorf("Route Run error: %v", err)
		}
	}()

	// Listen for the OS signal
	go func() {
		<-quit
		if err := facades.Route().Shutdown(); err != nil {
			facades.Log().Errorf("Route Shutdown error: %v", err)
		}

		os.Exit(0)
	}()

	select {}
}
