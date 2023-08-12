package download

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func DownloadFile(c *gin.Context, fileName string) {
	// Obtenha o caminho absoluto para o arquivo
	filePath := filepath.Join("files", fileName)

	// Verifique se o arquivo existe
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"Error": "Arquivo não encontrado",
		})
		return
	}

	// Defina os cabeçalhos de resposta para o download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "application/octet-stream")
	c.File(filePath)
}
