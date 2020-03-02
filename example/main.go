package main

import (
	"fmt"
	"github.com/hitech/gosdk"
	"github.com/hitech/gosdk/sdkcm"
	"github.com/hitech/gosdk/sdkplugin"
	"github.com/labstack/echo/v4"
 	"log"
)

func main() {
 	service := sdk.NewService(
		sdk.WithPlugin(sdkplugin.NewConfigPlugin(
			sdkplugin.WithName(sdkcm.KeyMainConfig),
			sdkplugin.WithConfig(sdkcm.NewSDKConfig()),
		)),
		sdk.WithPlugin(sdkplugin.NewGormPlugin(
			sdkplugin.WithName(sdkcm.KeyMainDb),
		)),
		sdk.WithPlugin(sdkplugin.NewNatPlugin(
			sdkplugin.WithName(sdkcm.KeyMainNat),
			sdkplugin.WithNatPluginRouter(natRouter),
		)),
		//sdk.WithPlugin(sdkplugin.NewGormPlugin("DB_TEST")),
		//sdk.WithPlugin(sdkplugin.NewNatPlugin(key.KeyMainNat)),
		sdk.WithPlugin(sdkplugin.NewRedisPlugin(
			sdkplugin.WithName(sdkcm.KeyMainRedis),
		)),
		sdk.WithPlugin(sdkplugin.NewSimpleHttpServer(
			sdkplugin.WithName(sdkcm.KeyMainHttpServer),
		)),
		//sdk.WithPlugin(sdkplugin.NewHttpServerPlugin(
		//	sdkplugin.WithName(sdkcm.KeyMainHttpServer),
		//	sdkplugin.WithHttpServerRouter(router),
		//)),
	)

	// Show service information.
	service.Info()

	// Start service
	service.Start()
}

func natRouter(s sdk.Service) {
	fmt.Println("nat")
}

func router(r *echo.Echo, s sdk.Service) {
	fmt.Println("echo")
	r.GET("/", func(c echo.Context) error {
		db := s.GetDb("DB")
		if db == nil {
			log.Fatalln("DB nil")
		}

		return c.String(200, "Hello World")
	})
}
