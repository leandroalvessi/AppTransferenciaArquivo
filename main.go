package main

import (
	conf_rede "AppTransferenciaArquivo/controllers/conf_rede"
	delete "AppTransferenciaArquivo/controllers/delete"
	download "AppTransferenciaArquivo/controllers/download"
	file "AppTransferenciaArquivo/controllers/files"
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

	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Erro ao obter o diretório de trabalho:", err)
		return
	}

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

	// Crie um novo QR code com o link
	qr, err := qrcode.New(hostName, qrcode.Medium)
	if err != nil {
		fmt.Println("Erro ao criar o QR code:", err)
		return
	}

	// Obtenha a matriz de caracteres representando o QR code
	qrASCII := qr.ToSmallString(false)

	// Exiba o QR code no terminal
	fmt.Println(qrASCII)

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

	router.GET("/files", file.ListFiles)

	router.POST("/download", func(c *gin.Context) {
		filename := c.PostForm("filename")
		download.DownloadFile(c, filename)
	})

	router.POST("/delete", func(c *gin.Context) {
		filename := c.PostForm("filename")
		delete.DeleteFile(c, workingDir, filename)
	})

	go conf_rede.OpenBrowser(hostName)

	router.Run(":" + port)
}
