/**
 *
 * By So http://sooo.site
 * -----
 *    Don't panic.
 * -----
 *
 */

package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Git-So/blog-api/models"
	"github.com/Git-So/blog-api/routers"
	"github.com/Git-So/blog-api/utils/cache"
	"github.com/Git-So/blog-api/utils/conf"
	"github.com/Git-So/blog-api/utils/helper"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/wonderivan/logger"
)

func main() {
	// log
	logPath := ".blog/log.json"
	helper.CreatePath("", logPath)
	helper.WriteFile(logPath, "{}")
	logger.SetLogger(logPath)

	// config
	confing := conf.Get()

	// cache
	ca := cache.New()
	ca.Connect()

	// model
	db := models.Get()
	defer db.Close()

	// router
	router := routers.Get()
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", confing.Server.Port),
		Handler:        router,
		ReadTimeout:    confing.Server.ReadTimeout * time.Second,
		WriteTimeout:   confing.Server.WriteTimeout * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err := gracehttp.Serve(server)
	if err != nil {
		logger.Fatal("gracehttp error: ", err)
	}

}
