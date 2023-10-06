package main

import (
	"AppTransferenciaArquivo/controllers/conf_rede"
	"AppTransferenciaArquivo/controllers/delete"
	"AppTransferenciaArquivo/controllers/download"
	"AppTransferenciaArquivo/controllers/files"
	"AppTransferenciaArquivo/controllers/global"
	"AppTransferenciaArquivo/controllers/upload"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	global.Port = "8080"

	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Erro ao obter o diretório de trabalho:", err)
		return
	}

	router.LoadHTMLGlob(workingDir + "/templates/*.html")

	global.HostGlobal = conf_rede.ConfigLink()

	// Crie um novo QR code com o link
	qr, err := qrcode.New(global.HostGlobal, qrcode.Medium)
	if err != nil {
		fmt.Println("Erro ao criar o QR code:", err)
		return
	}

	// Obtenha a matriz de caracteres representando o QR code
	//qrASCII := qr.ToSmallString(false)

	// Exiba o QR code no terminal
	//fmt.Println(qrASCII)

	// Obtenha a matriz de bytes da imagem PNG do QR code
	qrBytes, err := qr.PNG(256)
	if err != nil {
		fmt.Println("Erro ao gerar imagem PNG do QR code:", err)
		return
	}

	// Converta os bytes da imagem PNG em uma representação base64
	qrBase64 := base64.StdEncoding.EncodeToString(qrBytes)

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload.html", gin.H{
			"Year":   time.Now().Year(),
			"QRCode": qrBase64,
		})
	})

	router.POST("/upload", func(c *gin.Context) {
		upload.Upload(c, workingDir)
	})

	router.GET("/files", files.ListFiles)

	router.POST("/download", func(c *gin.Context) {
		filename := c.PostForm("filename")
		download.DownloadFile(c, filename)
	})

	router.POST("/delete", func(c *gin.Context) {
		filename := c.PostForm("filename")
		delete.DeleteFile(c, workingDir, filename)
	})

	go conf_rede.OpenBrowser(global.HostGlobal)

	router.Run(":" + global.Port)
}
