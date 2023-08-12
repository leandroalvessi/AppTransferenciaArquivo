package main

import (
	conf_rede "AppTransferenciaArquivo/controllers/conf_rede"
	delete "AppTransferenciaArquivo/controllers/delete"
	download "AppTransferenciaArquivo/controllers/download"
	file "AppTransferenciaArquivo/controllers/files"
	"AppTransferenciaArquivo/controllers/upload"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Erro ao obter o diretório de trabalho:", err)
		return
	}

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload.html", gin.H{
			"Year": time.Now().Year(),
		})
	})

	router.POST("/upload", func(c *gin.Context) {
		upload.Upload(c, workingDir)
	})

	router.GET("/files", file.ListFiles)

	router.POST("/download", func(c *gin.Context) {
		filename := c.PostForm("filename")
		download.DownloadFile(c, filename)
	})

	router.POST("/delete", func(c *gin.Context) {
		filename := c.PostForm("filename")
		delete.DeleteFile(c, workingDir, filename)
	})

	router.LoadHTMLGlob(workingDir + "/templates/*.html")

	interfaceRede, err := conf_rede.GetLocalIP()
	if err != nil {
		fmt.Println("Erro ao obter endereço IP:", err)
		os.Exit(1)
	}

	port := "8080"

	hostName := ""
	for _, pair := range interfaceRede {
		fmt.Printf("Interface: %s. Acesse em http://%s:%s\n", pair.Name, pair.IP, port)

		if pair.Name == "Wi-Fi" {
			hostName = "http://" + pair.IP + ":" + port + "/"
		} else {
			hostName = "http://" + pair.IP + ":" + port + "/"
		}
	}

	go conf_rede.OpenBrowser(hostName)

	router.Run(":" + port)
}
