package main

import (
    "github.com/gin-gonic/gin"
    "github.com/notedit/rtmp"
    "log"
    "net/http"
)

func main() {
    // RTMP server setup
    rtmpServer := rtmp.NewRtmpServer()

    rtmpServer.HandlePublish = func(conn *rtmp.Conn) {
        log.Println("New RTMP stream received")
    }

    go func() {
        log.Fatal(rtmpServer.ListenAndServe(":1935"))
    }()

    // HTTP server setup
    router := gin.Default()

    router.GET("/live", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", gin.H{})
    })

    router.Static("/video", "./video")

    router.LoadHTMLGlob("templates/*")

    log.Println("Starting HTTP server on :8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}
