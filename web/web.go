package web

import (
	// "crypto/md5"
	// "encoding/hex"
	// "encoding/json"
	// "fmt"
	"log"
	"net/http"
	"path/filepath"

	config "github.com/aibfarm/aibfarm/config"
	factory "github.com/aibfarm/aibfarm/web/factory"
	"github.com/gin-contrib/multitemplate"
	gin "github.com/gin-gonic/gin"
)

// !-- mydash 的web 服务器
func Start_WebServer() {

	r := gin.Default()
	// r.Use(cors.Default())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong-AIBFARM",
		})
	})

	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"aib": "farm",
	}))

	r.Static("/assets", "./web/templates/assets")
	r.HTMLRender = loadTemplates("./web/templates")

	//!-- homepage
	//主页进入
	authorized.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", gin.H{
			"menulist": factory.MakeMenu(),
		})
	})

	//!-------OKRice Debuging / Monitoring
	authorized.GET("/okriceV2.debug", func(c *gin.Context) {
		c.HTML(http.StatusOK, "okrice.debug.tmpl.html", gin.H{
			"menulist":     factory.MakeMenu(),
			"api_endpoint": config.AIBFARM_CONFIG.API_ENDPOINT, //"ws://192.168.2.131:20889",
		})
	})

	r.Run(config.OKRICE_DASH_PORT) // listen and serve on 0.0.0.0:

}

func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	layouts, err := filepath.Glob(templatesDir + "/layouts/*.tmpl.html")
	if err != nil {
		panic(err.Error())
	}
	log.Printf("%#v", layouts)

	includes, err := filepath.Glob(templatesDir + "/*.tmpl.html")
	if err != nil {
		panic(err.Error())
	}

	// Generate our templates map from our layouts/ and includes/ directories
	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		log.Printf("%#v", include)
		r.AddFromFiles(filepath.Base(include), files...)
	}
	return r
}
