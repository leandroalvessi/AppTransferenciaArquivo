package download

import (
	"io"
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

	// Abra o arquivo
	file, err := os.Open(filePath)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": "Erro ao abrir o arquivo",
		})
		return
	}
	defer file.Close()

	// Defina o tipo de conteúdo do cabeçalho de resposta
	c.Header("Content-Type", "application/octet-stream")

	// Defina o cabeçalho de resposta para permitir o download do arquivo com o nome original
	c.Header("Content-Disposition", "attachment; filename="+fileName)

	// Copie o conteúdo do arquivo para o corpo da resposta
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": "Erro ao enviar o arquivo",
		})
		return
	}
}
