package bot

import (
	"bytes"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jung-kurt/gofpdf"
	"image"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var channelID int64 = -1002073737722

type File struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

func (h *BotHandler) finalizeAddClient(chatId int64) {

	res, err := h.repo.ClientUser().Get(chatId)
	if err != nil {
		log.Println("Error creating client user:", err)
		return
	}

	lat, long, err := getLatAndLang(res.Location)
	if err != nil {
		log.Println("Error getting lat and lang:", err)
		return
	}

	url := fmt.Sprintf("https://www.google.com/maps?q=%f,%f", lat, long)
	fileName := generateFileName(res.UserName, res.ContractId)

	// CreateOne a new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 15) // Decrease font size for better text arrangement

	// Add content to the PDF
	pdf.SetX(10) // Set initial X position
	pdf.CellFormat(40, 10, "Contract ID:", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 10, res.ContractId, "", 0, "L", false, 0, "")
	pdf.Ln(-1) // Move to the next line without spacing

	pdf.SetX(10) // Reset X position
	pdf.CellFormat(40, 10, "Phone Number:", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 10, res.PhoneNumber, "", 0, "L", false, 0, "")
	pdf.Ln(-1) // Move to the next line without spacing

	pdf.SetX(10) // Reset X position
	pdf.CellFormat(40, 10, "Address:", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 10, res.Address, "", 0, "L", false, 0, "")
	pdf.Ln(-1) // Move to the next line without spacing

	pdf.SetX(10) // Reset X position
	pdf.CellFormat(40, 10, "Payment sum:", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 10, res.PaymentSum, "", 0, "L", false, 0, "")
	pdf.Ln(-1)

	pdf.SetX(10)
	pdf.CellFormat(40, 10, "Comment:", "", 0, "L", false, 0, "")
	pdf.MultiCell(0, 10, res.Comment, "", "L", false)
	pdf.Ln(-1)

	pdf.SetX(10)
	pdf.CellFormat(40, 10, "Location:", "", 0, "L", false, 0, "")
	pdf.Ln(-1)

	// Add a clickable link to the location
	pdf.SetX(10)
	pdf.WriteLinkString(0, url, "Map Link")
	pdf.Ln(5)

	// Get image dimensions
	imageWidth, imageHeight := getImageDimensions(res.AddressFotoPath)

	// Calculate the aspect ratio for resizing
	imageAspectRatio := imageWidth / imageHeight
	maxWidth := 140.0
	maxHeight := maxWidth / imageAspectRatio

	pdf.Cell(0, 10, "Address Photo:")
	pdf.ImageOptions(res.AddressFotoPath, 10.0, pdf.GetY()+10, maxWidth, maxHeight, false, gofpdf.ImageOptions{}, 0, "")

	// Get image dimensions
	imageWidth, imageHeight = getImageDimensions(res.AddressFotoPath)

	pdf.AddPage()
	pdf.Cell(0, 10, "Payment foto:")
	pdf.ImageOptions(res.PaymentFotoPath, 10.0, pdf.GetY()+10, maxWidth, maxHeight, false, gofpdf.ImageOptions{}, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		log.Println("Error saving PDF to buffer: ", err)
		return
	}

	doc := tgbotapi.NewDocument(channelID, tgbotapi.FileBytes{Name: fileName, Bytes: buf.Bytes()})
	h.bot.Send(doc)
	h.SendTextMessage(res, "Kiritgan malumotlaringiz jo'natildi")
}

// Function to get the dimensions of an image file
func getImageDimensions(imagePath string) (float64, float64) {
	file, err := os.Open(imagePath)
	if err != nil {
		log.Println("Error opening image file:", err)
		return 0, 0
	}
	defer file.Close()

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		log.Println("Error decoding image configuration:", err)
		return 0, 0
	}

	return float64(img.Width), float64(img.Height)
}

func generateFileName(name, contactID string) string {
	fileName := fmt.Sprintf("%s_%s.pdf", name, contactID)
	return fileName
}

// Method to download an image through the bot and upload it to the images folder
func (h *BotHandler) savePhotoToFolder(bot *tgbotapi.BotAPI, fileID string) (string, error) {
	log.Println("FileID: ", fileID)
	// Download the image from Telegram using the bot
	file, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		return "", err
	}

	// Download the file contents
	resp, err := http.Get("https://api.telegram.org/file/bot" + bot.Token + "/" + file.FilePath)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the file data
	imageData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Define the path to save the image
	imagesFolder := "images/"
	err = os.MkdirAll(imagesFolder, os.ModePerm)
	if err != nil {
		return "", err
	}
	imageFilePath := filepath.Join(imagesFolder, fileID+".jpg")

	// Write the image data to a file
	err = ioutil.WriteFile(imageFilePath, imageData, 0644)
	if err != nil {
		return "", err
	}

	return imageFilePath, nil
}

func getLatAndLang(location string) (float64, float64, error) {
	latLang := strings.Split(location, ",")
	if len(latLang) != 2 {
		return 0, 0, fmt.Errorf("invalid location format")
	}

	lat, err := strconv.ParseFloat(latLang[0], 64)
	if err != nil {
		return 0, 0, err
	}

	lang, err := strconv.ParseFloat(latLang[1], 64)
	if err != nil {
		return 0, 0, err
	}

	return lat, lang, nil
}
