package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

var cnt int

func main() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.GET("/fire-and-forget", fireAndForget)
	router.GET("/load-concurrently", loadConcurrently)
	router.GET("/waterfall", waterfall)

	router.Run("127.0.0.1:8080")
}

func fireAndForget(c *gin.Context) {
	start := time.Now()

	go hGet("http://www.gog.com")
	go hGet("http://www.google.com")
	go hGet("http://www.bing.com")

	timer := time.Since(start).Seconds()

	c.String(http.StatusOK, "request done in %f \n", timer)

	return
}

func loadConcurrently(c *gin.Context) {
	start := time.Now()

	var results []*http.Response
	ch := make(chan *http.Response)

	go hGetC("http://www.gog.com", ch)
	go hGetC("http://www.google.com", ch)
	go hGetC("http://www.bing.com", ch)

	results = append(results, <-ch)
	results = append(results, <-ch)
	results = append(results, <-ch)

	log.Printf("results: %d, %d, %d \n", results[0].StatusCode, results[1].StatusCode, results[2].StatusCode)

	timer := time.Since(start).Seconds()

	c.String(http.StatusOK, "request done (3 calls) in %f\n", timer)
}

func waterfall(c *gin.Context) {
	start := time.Now()

	hGet("http://www.gog.com")
	hGet("http://www.google.com")
	hGet("http://www.bing.com")

	timer := time.Since(start).Seconds()

	c.String(http.StatusOK, "request done in %f \n", timer)

	return
}

func hGet(url string) *http.Response {
	start := time.Now()
	cnt++

	response, error := http.Get(url)
	timer := time.Since(start).Seconds()

	if error != nil {
		log.Printf("get %s failed after %f: %s (%d) \n", url, timer, error, cnt)
	} else {
		log.Printf("get %s done in %f (%d) \n", url, timer, cnt)
	}

	return response
}

func hGetC(url string, c chan *http.Response) {
	c <- hGet(url)
}
