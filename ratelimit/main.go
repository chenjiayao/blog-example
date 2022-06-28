package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"ratelimit/flowlimit"
	"ratelimit/middlewares"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.Use(middlewares.IPLimitRaterMiddleware).GET("/ip_rate_limit", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})

	})

	r.GET("/download_1", func(c *gin.Context) {
		b := RandStringBytesMaskImprSrc(10 * 1024 * 1024)
		buf := bytes.NewBuffer(b)
		c.Header("Content-Length", fmt.Sprintf("%d", len(b)))
		c.Header("Content-disposition", "attachment;filename=download")

		//每秒写 10*1024 byte 到 c.Writer
		for range time.Tick(1 * time.Second) {
			_, err := io.CopyN(c.Writer, buf, 10*1024) //10KB/s
			if err == io.EOF {
				break
			}
		}
	})

	r.GET("/download_2", func(c *gin.Context) {
		b := RandStringBytesMaskImprSrc(10 * 1024 * 1024)
		buf := bytes.NewBuffer(b)
		c.Header("Content-Length", fmt.Sprintf("%d", len(b)))
		c.Header("Content-disposition", "attachment;filename=download")

		lr := flowlimit.NewLimitReader(buf)
		lr.SetRateLimit(10 * 1024)
		io.Copy(c.Writer, lr)
	})

	r.Run()
}

//生成一个 n byte 的文件
func RandStringBytesMaskImprSrc(n int) []byte {
	const (
		letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)

	var src = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return b
}
