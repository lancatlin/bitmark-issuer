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
	db.AutoMigrate(&Asset{}, &Issue{}, &User{}, &URL{})

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

func withCondition(callback func(*gin.Context), condition Condition) func(*gin.Context) {
	return func(c *gin.Context) {
		user := getUser(c)
		if condition == ConditionLogin && !user.IsLogin {
			// permission denied
			page := struct {
				User
				message
			}{
				User: user,
				message: message{
					Title:      "權限不足",
					Content:    "請登入後再訪問此頁面",
					Target:     "/login",
					TargetName: "登入",
				},
			}
			c.HTML(401, "msg.html", page)
			return
		}
		if condition == ConditionLogout && user.IsLogin {
			// redirect to /
			c.Redirect(303, "/")
			return
		}
		callback(c)
	}
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("./templates/*.html")
	r.Static("/static", "./static")
	r.GET("/", index)
	r.GET("/signup", staticPage("sign-up.html", ConditionLogout))
	r.POST("/signup", withCondition(signUp, ConditionLogout))
	r.GET("/login", staticPage("login.html", ConditionLogout))
	r.POST("/login", withCondition(login, ConditionLogout))
	r.GET("/logout", withCondition(logout, ConditionLogin))
	r.GET("/new", staticPage("new-asset.html", ConditionLogin))
	r.POST("/assets", withCondition(newAsset, ConditionLogin))
	r.GET("/assets/:id", withCondition(assetInfo, ConditionLogin))
	r.GET("/get/:id", getAsset)
	r.POST("/get/:id")
	r.Run()
}
