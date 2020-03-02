### Go SDK

##### Getting Started: 

1. Create .env file in root project.
 
```.env
#Set APP PORT 
APP_PORT=8888 
```

2. Example. 
```go
package main

import (
	"github.com/labstack/echo/v4"
 	"github.com/hitechqb/gosdk"
	"github.com/hitechqb/gosdk/sdkcm"
	"github.com/hitechqb/gosdk/sdkplugin"
	"net/http"
)

func main() {
	sdk.NewService(
		sdk.WithPlugin(sdkplugin.NewConfigPlugin(
			sdkplugin.WithName(sdkcm.KeyMainConfig),
			sdkplugin.WithConfig(sdkcm.NewSDKConfig()),
		)),
		sdk.WithPlugin(sdkplugin.NewHttpServerPlugin(
			sdkplugin.WithName(sdkcm.KeyMainHttpServer),
			sdkplugin.WithHttpServerRouter(router),
		)),
	).Start()
}

func router(s sdk.Service) {
	r := s.GetHttpServer(sdkcm.KeyMainHttpServer)
	r.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello world")
	})
}
```

3. Run APP
```.env
type: 
go run main.go serve
```

4. Result
```.env
📦 Plugins:
1. 👉 CONFIG_MAIN 
2. 👉 HTTP_SERVER_MAIN 
=================================================================================
echo
⇨ http server started on [::]:8888
```

#### Cheer!
