package web

import (
	// "crypto/md5"
	// "encoding/hex"
	// "encoding/json"
	// "fmt"

	"log"
	"net/http"
	"path/filepath"
	"time"

	config "github.com/aibfarm/config"
	factory "github.com/aibfarm/web/factory"
	"github.com/gin-contrib/multitemplate"
	gin "github.com/gin-gonic/gin"
	"github.com/k0kubun/pp"

	jwt "github.com/appleboy/gin-jwt/v2"
)

const identityKey = "id"

type User struct { //!--  用户
	UserID   int64  `gorm:"primary_key;auto_increment;not_null"` //!-
	Name     string `gorm:"index;not_null"`
	UserName string `gorm:"unique;not_null"`
	Reward   int64  `gorm:"index;not_null"`
	// FirstName string
	// LastName  string
}

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// !-- mydash 的web 服务器
func Start_WebServer() {
	// var str string
	var err error

	r := gin.Default()
	// r.Use(cors.Default())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong-AIBFARM",
		})
	})

	//!--JWT Login
	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "temple zone",
		Key:         []byte("secret key should stay in the config file"),
		Timeout:     24 * time.Hour,
		MaxRefresh:  24 * time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				UserName: claims[identityKey].(string),
				// FirstName: "Mika", //这里是验证返回的名字
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password

			log.Printf("%s:%s",
				loginVals.Username,
				loginVals.Password,
			)

			if (userID == "temple@9cat.net" && password == "temple123") || (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
				return &User{
					UserName: userID,
					// LastName:  "Bo-Yi",
					// FirstName: "Wu",
				}, nil
			}
			if userID == "mika@9cat.net" && password == "temple1119" {
				return &User{
					UserName: userID,
					// LastName:  "Temple",
					// FirstName: "Mika",
				}, nil
			}

			if userID == "tjs@templexp.com" && password == "aAa33.88" {
				return &User{
					UserName: userID,
					// LastName:  "Temple",
					// FirstName: "Mika",
				}, nil
			}

			if userID == "legend.cord@gmail.com" && password == "vBD$41" {
				return &User{
					UserName: userID,
					// LastName:  "Temple",
					// FirstName: "Mika",
				}, nil
			}

			if userID == "trade@dongchenxie.com" && password == "aTm.230623" {
				return &User{
					UserName: userID,
					// LastName:  "Temple",
					// FirstName: "Mika",
				}, nil
			}

			if userID == "ireneilu@live.hk" && password == "aTm.230629" {
				return &User{
					UserName: userID,
				}, nil
			}

			if userID == "442887516@qq.com" && password == "aTm.230630" {
				return &User{
					UserName: userID,
				}, nil
			}

			if userID == "tonsetisolution@gmail.com" && password == "aTm.240818" {
				return &User{
					UserName: userID,
				}, nil
			}

			if userID == "feiya200@gmail.com" && password == "feiya.240819" {
				return &User{
					UserName: userID,
				}, nil
			}

			if userID == "87619687@qq.com" && password == "jisen.240820" {
				return &User{
					UserName: userID,
				}, nil
			}

			if userID == "colewu9712@gmail.com" && password == "col.240820" {
				return &User{
					UserName: userID,
				}, nil
			}

			if userID == "Robert_anu@hotmail.com" && password == "ra.240820" {
				return &User{
					UserName: userID,
				}, nil
			}

			if userID == "624355527@qq.com" && password == "jc.240820" {
				return &User{
					UserName: userID,
				}, nil
			}

			if userID == "lining41081@gmail.com" && password == "dudu.240821" {
				return &User{
					UserName: userID,
				}, nil
			}

			if userID == "by.bao@hotmail.com" && password == "mika.240825" {
				return &User{
					UserName: userID,
				}, nil
			}

			if userID == "xxxottochanxxx@gmail.com" && password == "zx.240829" {
				return &User{
					UserName: userID,
				}, nil
			}

			if userID == "richard.zhang.9cat@gmail.com" && password == "richard.240903" {
				return &User{
					UserName: userID,
				}, nil
			}

			if userID == "lemoine5332@gmail.com" && password == "lou.2409a20" {
				return &User{
					UserName: userID,
				}, nil
			}

			if userID == "cnimaiv@gmail.com" && password == "js.20240923" {
				return &User{
					UserName: userID,
				}, nil
			}

			if userID == "martin@9cat.net" && password == "temple1204." {
				return &User{
					UserName: userID,
					// LastName:  "Temple",
					// FirstName: "Martin",
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},

		Unauthorized: func(c *gin.Context, code int, message string) {
			// c.JSON(code, gin.H{
			// 	"code":    code,
			// 	"message": message,
			// })

			c.Redirect(http.StatusFound, "/login")
		},

		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	errInit := authMiddleware.MiddlewareInit()
	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	r.POST("/login", authMiddleware.LoginHandler)
	r.GET("/logout", authMiddleware.LogoutHandler)

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	authorized := r.Group("/")
	authorized.GET("/refresh_token", authMiddleware.RefreshHandler)

	authorized.Use(authMiddleware.MiddlewareFunc())
	{
		//!--这里作为验证之后的操作：

		authorized.GET("/hello", func(c *gin.Context) {
			claims := jwt.ExtractClaims(c)
			user, _ := c.Get(identityKey)
			u := user.(*User)
			pp.Print(u)

			c.JSON(200, gin.H{
				"userID":   claims[identityKey],
				"userName": u.UserName,
				// "firstName": u.FirstName,
				// "text":      "Hello World.",
			})
		})

		//!-- homepage
		//主页进入
		authorized.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.tmpl.html", gin.H{
				"menulist": factory.MakeMenu(),
			})
		})

		//!-------OKRice Debuging / Monitoring
		authorized.GET("/aibfarm.debug", func(c *gin.Context) {
			c.HTML(http.StatusOK, "aibfarm.debug.tmpl.html", gin.H{
				"menulist":     factory.MakeMenu(),
				"api_endpoint": config.AIBFarmConfig.API_ENDPOINT, //"ws://192.168.2.131:20889",
			})
		})

		//!--END OF 验证之后的操作：

	}

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	//!-------

	// authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
	// 	"aib": "farm",
	// }))

	r.Static("/assets", "./web/templates/assets")
	r.HTMLRender = loadTemplates("./web/templates")

	r.Run(config.AIBFARM_DASH_PORT) // listen and serve on 0.0.0.0:

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

	adminLayouts, err := filepath.Glob(templatesDir + "/layouts/html-base.html")
	if err != nil {
		panic(err.Error())
	}

	admins, err := filepath.Glob(templatesDir + "/html/*.html")
	if err != nil {
		panic(err.Error())
	}
	// Generate our templates map from our adminLayouts/ and admins/ directories
	for _, admin := range admins {
		layoutCopy := make([]string, len(adminLayouts))
		copy(layoutCopy, adminLayouts)
		files := append(layoutCopy, admin)
		r.AddFromFiles(filepath.Base(admin), files...)
	}
	return r
}
