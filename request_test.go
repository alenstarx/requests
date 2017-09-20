package requests

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"context"
	"log"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
)

func testRequestBase(t *testing.T) {
	req := NewRequest()
	assert.Equal(t, req.Err(), nil)
	req.SetUrl("https://baidu.com")
	assert.Equal(t, req.Err(), nil)
	resp := req.Get()
	assert.Equal(t, resp.Err(), nil)
	body := resp.Bytes()
	assert.Equal(t, resp.Err(), nil)
	assert.Equal(t, body != nil, true)
}

func testRequestGet(t *testing.T, r *Request) {
	resp := r.SetUrl("http://127.0.0.1:8080/").Get()
	assert.Equal(t, resp.Err(), nil)
	assert.Equal(t, "requests", resp.GetCookie()["SessionId"])
	resp.StoreCookie(r)
	r.SetUrl("http://127.0.0.1:8080/user")
	assert.Equal(t, r.Err(), nil)
	resp = r.Get()
	assert.Equal(t, resp.Err(), nil)
	body := resp.Bytes()
	assert.Equal(t, resp.Err(), nil)
	assert.Equal(t, body != nil, true)

	param := make(map[string]string)
	param["name"] = "Tom"
	param["age"] = "99"
	body = r.SetParam(param).Get().Bytes()
	assert.Equal(t, string(body), "name=Tom;age=99")
}

func testRequestPost(t *testing.T, r *Request) {
	form := make(map[string]string)
	form["name"] = "Tom"
	form["age"] = "99"
	resp := r.PostForm(form)
	assert.Equal(t, nil, resp.err)
	body := resp.Bytes()
	assert.Equal(t, string(body), "name=Tom;age=99")

}

func initGin(t *testing.T) *gin.Engine {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.SetCookie("SessionId", "requests", 1000, "/", "", true, true)
		c.String(http.StatusOK, "Welcome Gin Server")
	})
	router.GET("/user", func(c *gin.Context) {
		name := c.Query("name")
		age := c.Query("age")
		cookie, err := c.Cookie("SessionId")
		assert.Equal(t, nil, err)
		assert.Equal(t, "requests", cookie)
		c.String(http.StatusOK, "name="+name+";age="+age)
	})
	router.POST("/user", func(c *gin.Context) {
		name := c.PostForm("name")
		age := c.PostForm("age")
		c.String(http.StatusOK, "name="+name+";age="+age)
	})
	return router
}
func TestNewRequest(t *testing.T) {

	srv := &http.Server{
		Addr:    ":8080",
		Handler: initGin(t),
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	testRequestBase(t)
	req := NewRequest()
	assert.Equal(t, req.Err(), nil)
	testRequestGet(t, req)
	testRequestPost(t, req)

	time.Sleep(time.Second * 3)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
