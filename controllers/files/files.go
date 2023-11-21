package files

import (
	"AppTransferenciaArquivo/controllers/global"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
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

	// Crie um novo QR code com o link
	qr, err := qrcode.New(global.HostGlobal+"files", qrcode.Medium)
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

	c.HTML(http.StatusOK, "files.html", gin.H{
		"FileNames": fileNames,
		"Year":      time.Now().Year(),
		"QRCode":    qrBase64,
	})

}
func GeneratePDF(pdfPath string, c *gin.Context) error {
	// Create a new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Set font for the entire document
	pdf.SetFont("Arial", "", 12)

	// Set header function
	pdf.SetHeaderFunc(func() {
		pdf.SetY(10) // Adjust the Y position as needed
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(0, 10, "My Document Header")
		pdf.Ln(10) // Move the cursor down by 10 units
	})

	// Add a page to the document
	pdf.AddPage()

	// Set footer function (optional)
	pdf.SetFooterFunc(func() {
		pdf.SetXY(-20, -15)
		pdf.SetFont("Arial", "I", 8)
		pdf.Cell(0, 10, fmt.Sprintf("Page %d", pdf.PageNo()))
	})

	// Set font back to regular for the main content
	pdf.SetFont("Arial", "", 12)

	// Add content to the document
	for i := 0; i < 50; i++ {
		pdf.Cell(0, 10, fmt.Sprintf("Line %d", i+1))
		pdf.Ln(10)
	}

	// Output the PDF to a file
	return pdf.OutputFileAndClose(pdfPath)
}

func ViewPDF(c *gin.Context) {
	// Define the directory where PDF files are stored
	pdfDirectory := "pdf"

	// Create a filename for the "Hello World" PDF
	filename := c.PostForm("filename") + ".pdf"

	// Create the full path to the PDF file
	pdfPath := path.Join(pdfDirectory, filename)

	// Generate the "Hello World" PDF
	err := GeneratePDF(pdfPath, c)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": "Error generating the PDF file",
		})
		return
	}

	// Read the content of the PDF file
	pdfContent, err := ioutil.ReadFile(pdfPath)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": "Error reading the PDF file",
		})
		return
	}

	// You can use the PDF content as needed, for example, send it to the client
	c.Data(http.StatusOK, "application/pdf", pdfContent)
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
