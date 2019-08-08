package main

import (
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/joho/godotenv/autoload"
	"os"
	sdk "github.com/bitmark-inc/bitmark-sdk-go"
	"net/http"
	"time"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("sqlite3", "data.sqlite")
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Asset{}, &Issue{}, &User{})

	// init bitmark testnet
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	config := &sdk.Config{
		APIToken: os.Getenv("API_TOKEN"),
		Network: "testnet",
		HTTPClient: httpClient,
	}
	sdk.Init(config)	
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("./templates/*.html")
	r.Static("/static", "./static")
	r.GET("/sign-up", func (c *gin.Context) {
		if getUser(c).IsLogin {
			c.Redirect(303, "/")
			return
		}
		c.HTML(200, "sign-up.html", nil)
	})
	r.POST("/sign-up", signUp)
	r.GET("/login", func (c *gin.Context) {
		if getUser(c).IsLogin {
			c.Redirect(303, "/")
			return
		}
		c.HTML(200, "login.html", nil)
	})
	r.POST("/login", login)
	r.Run()
}
