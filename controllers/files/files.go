package files

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type File struct {
	Name string
	Size string
}

func ListFiles(c *gin.Context) {
	// Obtenha o diretório desejado para listar os arquivos
	directory := "files"

	// Crie uma estrutura para armazenar os nomes dos arquivos
	fileNames := []File{}

	// Percorra os arquivos no diretório
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Verifique se o caminho é um arquivo
		if !info.IsDir() {
			file := File{
				Name: info.Name(),
				Size: formatSize(info.Size()),
			}

			fileNames = append(fileNames, file)

		}

		return nil
	})

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "files.html", gin.H{
		"FileNames": fileNames,
		"Year":      time.Now().Year(),
	})

}

func formatSize(size int64) string {
	const (
		B  = 1 << (10 * 0) // 1 byte
		KB = 1 << (10 * 1) // 1 kilobyte
		MB = 1 << (10 * 2) // 1 megabyte
		GB = 1 << (10 * 3) // 1 gigabyte
		TB = 1 << (10 * 4) // 1 terabyte
	)

	switch {
	case size >= TB:
		return fmt.Sprintf("%.2f TB", float64(size)/TB)
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/KB)
	default:
		return fmt.Sprintf("%d B", size)
	}
}
