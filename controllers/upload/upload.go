package upload

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context, savePath string) {
	// Faz o parse do formulário com multipart
	err := c.Request.ParseMultipartForm(32 << 20) // Define o tamanho máximo do formulário (32MB neste exemplo)
	if err != nil {
		c.HTML(http.StatusBadRequest, "upload.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	// Obtém todos os arquivos enviados
	form := c.Request.MultipartForm
	files := form.File["file"]

	for _, file := range files {
		// Concatena o caminho absoluto com o nome do arquivo
		filePath := filepath.Join(savePath, "files", file.Filename)

		// Salva o arquivo no sistema de arquivos
		err = c.SaveUploadedFile(file, filePath)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "upload.html", gin.H{
				"Error": err.Error(),
			})
			return
		}
	}

	c.Redirect(http.StatusFound, "/files")
}
