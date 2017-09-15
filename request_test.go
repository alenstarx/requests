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


func TestNewRequest(t *testing.T) {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()


	req := NewRequest()
	assert.Equal(t, req.Err(), nil)
	req.SetUrl("https://baidu.com")
	assert.Equal(t, req.Err(), nil)
	resp := req.Get()
	assert.Equal(t, resp.Err(), nil)
	body := resp.Bytes()
	assert.Equal(t, resp.Err(), nil)
	assert.Equal(t, body != nil ,true )

	req.SetUrl("http://127.0.0.1:8080")
	assert.Equal(t, req.Err(), nil)
	resp=req.Get()
	assert.Equal(t, resp.Err(), nil)
	body = resp.Bytes()
	assert.Equal(t, resp.Err(), nil)
	assert.Equal(t, body != nil ,true )


	time.Sleep(time.Second * 1)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}