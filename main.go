package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	sdk "github.com/bitmark-inc/bitmark-sdk-go"
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/joho/godotenv/autoload"
)

var db *gorm.DB

var env = struct {
	Host string
}{}

func init() {
	var err error
	db, err = gorm.Open("sqlite3", "data.sqlite")
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Asset{}, &Issue{}, &User{})

	env.Host = os.Getenv("HOST")

	// init bitmark testnet
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	config := &sdk.Config{
		APIToken:   os.Getenv("API_TOKEN"),
		Network:    sdk.Testnet,
		HTTPClient: httpClient,
	}
	sdk.Init(config)
}

// Condition define the condition to a page
type Condition int

const (
	// ConditionNotRequire allow all
	ConditionNotRequire Condition = iota
	// ConditionLogin only allow user who has logged in
	ConditionLogin
	// ConditionLogout only allow user who hasn't logged in
	ConditionLogout
)

func staticPage(filePath string, condition Condition) func(*gin.Context) {
	return func(c *gin.Context) {
		user := getUser(c)
		switch condition {
		case ConditionLogin:
			if !user.IsLogin {
				c.Redirect(303, "/login")
				return
			}
		case ConditionLogout:
			if user.IsLogin {
				c.Redirect(303, "/")
				return
			}
		}
		c.HTML(200, filePath, user)
	}
}

func index(c *gin.Context) {
	user := getUser(c)
	file, err := ioutil.ReadFile("README.md")
	if err != nil {
		panic(err)
	}
	html := markdown.ToHTML(file, nil, nil)
	page := struct {
		User
		Content template.HTML
	}{
		User:    user,
		Content: template.HTML(html),
	}
	c.HTML(200, "index.html", page)
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("./templates/*.html")
	r.Static("/static", "./static")
	r.GET("/", index)
	r.GET("/signup", staticPage("sign-up.html", ConditionLogout))
	r.POST("/signup", signUp)
	r.GET("/login", staticPage("login.html", ConditionLogout))
	r.POST("/login", login)
	r.GET("/logout", logout)
	r.GET("/new", staticPage("new-asset.html", ConditionLogin))
	r.POST("/assets", newAsset)
	r.GET("/assets/:id", func(c *gin.Context) {
		user := getUser(c)
		var a Asset
		if err := db.Where("id = ?", c.Param("id")).First(&a).Error; err != nil {
			// not found
			c.String(404, err.Error())
			return
		}
		page := struct {
			User
			Asset
			Host string
		}{
			User:  user,
			Asset: a,
			Host:  env.Host,
		}
		c.HTML(200, "created-asset.html", page)
	})
	r.POST("/assets/:id", staticPage("", ConditionNotRequire))
	r.GET("/assets/:id/get", getAsset)
	r.Run()
}
