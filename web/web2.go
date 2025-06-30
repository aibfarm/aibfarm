package web

import (
	// "crypto/md5"
	// "encoding/hex"
	// "encoding/json"
	// "fmt"

	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	config "github.com/aibfarm/config"
	factory "github.com/aibfarm/web/factory"
	gin "github.com/gin-gonic/gin"
	"github.com/k0kubun/pp"

	jwt "github.com/appleboy/gin-jwt/v2"
)

var StatusMsg *config.ClientData

// !-- mydash 的web 服务器 v2
func WebServerV2() {
	// var str string
	var err error

	//!-700 汇报系统
	StatusMsg = &config.ClientData{
		Name:  "AIBFARM 客服UI接口",
		Ports: fmt.Sprintf("%s", config.AIBFARM_DASH_PORT2),
		DBs:   "N/A",

		Message: "AIBFARM 客服UI接口 - web.WebServerV2() ",
	}
	_, err = config.ReportStatusJSON(StatusMsg) // _statStop

	r := gin.Default()
	// r.Use(cors.Default())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong-AIBFARM-Dashboard-Web",
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

			//TODO - 访问服务器API来确认用户是否是正常用户：
			// if userID == "mika@9cat.net" && password == "temple1119" {
			// 	return &User{
			// 		UserName: userID,
			// 	}, nil
			// }

			u := VerifyUser{
				Username: userID,
				Password: password,
				//TODO F2A Here
			}

			AIBFarm_Web_VerifyUser(&u)
			if u.IsMatch {
				return &User{
					UserName: userID,
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

	r.Static("/assets", "./web/templates/assets")
	r.HTMLRender = loadTemplates("./web/templates")

	r.Run(config.AIBFARM_DASH_PORT2) // listen and serve on 0.0.0.0:

}

type VerifyUser struct {
	Username string
	Password string
	IsMatch  bool
}

type VerifyUser_Response struct {
	Msg  string
	Info string
	Code int
}

func AIBFarm_Web_VerifyUser(user *VerifyUser) error {
	var str string
	var err error

	str = fmt.Sprintf("AIBFarm Web Verify User \n err:: %s", err)
	log.Print(str)

	// app := VerifyUser{
	// 	Username: "tanjiashun@templexp.com",
	// 	Password: "asdfAsdf..12",
	// }

	jsonStr, err := json.Marshal(user)
	if err != nil {
		log.Print(err)
		return err
	}

	req, err := http.NewRequest(
		"POST",
		"https://new.aibfarm.com/aibfarm.web.VerifyUser",
		bytes.NewBuffer(jsonStr),
	)

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("temple", "xXxxxXxx.xXx")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return err
	}
	defer resp.Body.Close()

	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)

	// body, _ := io.ReadAll(resp.Body)
	// fmt.Println("response Body:", string(body))

	r := VerifyUser_Response{}
	user.IsMatch = false
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		log.Print(err)
		return err
	}

	if r.Code == 0 {
		user.IsMatch = true
	} else {
		pp.Print(r)
	}

	return err
}
