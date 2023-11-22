package delete

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func DeleteFile(c *gin.Context, savePath, fileName string) {
	// Exemplo de código para deletar o arquivo usando a função os.Remove
	filePath := filepath.Join(savePath, "files", fileName)
	err := os.Remove(filePath)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	filePath = filepath.Join(savePath, "pdf", fileName)
	err = os.Remove(filePath + ".pdf")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	c.Redirect(http.StatusFound, "/files")
}
