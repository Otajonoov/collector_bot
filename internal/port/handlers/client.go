package handlers

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
	tb "gopkg.in/tucnak/telebot.v2"
	"image"
	"iman_tg_bot/internal/model"
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

func (h *Handler) AddClient(m *tb.Message) {
	userInfoCache := NewCache()

	steps := []struct {
		prompt  string
		setter  func(clientID int64, value string)
		hasNext bool
	}{
		{"Contract ID", func(clientID int64, value string) {
			userInfoCache.SetClientInfo(clientID, "contract_id", value, 180)
		}, true},
		{"Phone Number", func(clientID int64, value string) {
			userInfoCache.SetClientInfo(clientID, "phone_number", value, 180)
		}, true},
		{"Address", func(clientID int64, value string) {
			userInfoCache.SetClientInfo(clientID, "address", value, 180)
		}, true},
		{"Payment sum", func(clientID int64, value string) {
			userInfoCache.SetClientInfo(clientID, "payment_sum", value, 180)
		}, true},
		{"Comment", func(clientID int64, value string) {
			userInfoCache.SetClientInfo(clientID, "comment", value, 180)
		}, true},
		{"Location", func(clientID int64, value string) {
			userInfoCache.SetClientInfo(clientID, "location", value, 180)
		}, true},
		{"Address foto", func(clientID int64, value string) {
			userInfoCache.SetClientInfo(clientID, "address_foto_path", value, 180)
		}, true},
		{"Payment foto", func(clientID int64, value string) {
			userInfoCache.SetClientInfo(clientID, "payment_foto_path", value, 180)
		}, true},
	}

	h.handleSteps(m, userInfoCache, steps, 0)
}

func (h *Handler) handleSteps(m *tb.Message, userInfoCache *Cache, steps []struct {
	prompt  string
	setter  func(userId int64, value string)
	hasNext bool
}, currentStep int) {
	if currentStep >= len(steps) {
		h.finalizeAddClient(m, userInfoCache)
		return
	}

	step := steps[currentStep]
	h.b.Send(m.Sender, step.prompt)

	h.b.Handle(tb.OnText, func(msg *tb.Message) {
		step.setter(msg.Chat.ID, msg.Text)
		currentStep++
		h.handleSteps(m, userInfoCache, steps, currentStep)
	})

	h.b.Handle(tb.OnPhoto, func(msg *tb.Message) {
		if !steps[currentStep].hasNext {
			log.Println("Photo received at unexpected step:", currentStep)
			return
		}

		// Save the photo to the images folder
		photoPath, err := h.savePhotoToFolder(msg.Photo)
		if err != nil {
			log.Println("Error saving address photo:", err)
			return
		}

		_, exists := userInfoCache.GetClientInfo(msg.Chat.ID, "address_foto_path")
		if !exists {
			step.setter(msg.Chat.ID, photoPath)
			currentStep++
			h.handleSteps(m, userInfoCache, steps, currentStep)
			return
		}
		_, exists = userInfoCache.GetClientInfo(msg.Chat.ID, "payment_foto_path")
		if !exists {
			step.setter(msg.Chat.ID, photoPath)
			currentStep++
			h.handleSteps(m, userInfoCache, steps, currentStep)
			return
		}

		h.handleSteps(m, userInfoCache, steps, currentStep)
	})

	h.b.Handle(tb.OnLocation, func(msg *tb.Message) {
		// Check if the message contains location data
		if msg.Location == nil {
			log.Println("No location data found in the message")
			return
		}

		lat := float32(msg.Location.Lat)
		long := float32(msg.Location.Lng)
		step.setter(msg.Chat.ID, fmt.Sprintf("%f,%f", lat, long))

		currentStep++
		h.handleSteps(m, userInfoCache, steps, currentStep)
	})
}

func (h *Handler) finalizeAddClient(m *tb.Message, userInfoCache *Cache) {

	chatID := m.Chat.ID

	contact_id, _ := userInfoCache.GetClientInfo(chatID, "contract_id")
	phone_number, _ := userInfoCache.GetClientInfo(chatID, "phone_number")
	address, _ := userInfoCache.GetClientInfo(chatID, "address")
	payment_sum, _ := userInfoCache.GetClientInfo(chatID, "payment_sum")
	comment, _ := userInfoCache.GetClientInfo(chatID, "comment")
	location, _ := userInfoCache.GetClientInfo(chatID, "location")
	address_foto_path, _ := userInfoCache.GetClientInfo(chatID, "address_foto_path")
	payment_foto_path, _ := userInfoCache.GetClientInfo(chatID, "payment_foto_path")

	res, err := h.repo.ClientUser().CreateOne(&model.Client{
		ContractId:      contact_id,
		PhoneNumber:     phone_number,
		Address:         address,
		PaymentSum:      payment_sum,
		Comment:         comment,
		Location:        location,
		AddressFotoPath: address_foto_path,
		PaymentFotoPath: payment_foto_path,
		ChatId:          chatID,
	})
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
	fileName := h.generateFileName(m.Sender, contact_id)

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
	//pdf.CellFormat(0, 10, res.Location, "", 0, "L", false, 0, "")
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

	if payment_foto_path != "" {
		pdf.AddPage()
		pdf.Cell(0, 10, "Payment foto:")
		pdf.ImageOptions(res.PaymentFotoPath, 10.0, pdf.GetY()+10, maxWidth, maxHeight, false, gofpdf.ImageOptions{}, 0, "")
	}

	// Write the PDF content to a buffer
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		log.Println("Error saving PDF to buffer: ", err)
		h.b.Send(m.Sender, "Error saving PDF")
		return
	}

	// Send the PDF document to the group
	_, err = h.b.Send(tb.ChatID(channelID), &tb.Document{
		File:     tb.FromReader(&buf),
		FileName: fileName,
	})
	if err != nil {
		log.Println("Error sending PDF to group:", err)
		return
	}
	h.b.Send(m.Sender, "Malumot muvaffaqiyatli yuklandi.")
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

func (h *Handler) generateFileName(sender *tb.User, contactID string) string {
	fileName := "client_information.pdf"
	if sender != nil && (sender.FirstName != "" || sender.LastName != "") {
		fileName = fmt.Sprintf("%s_%s_%s.pdf", sender.FirstName, sender.LastName, contactID)
	}
	return fileName
}

func (h *Handler) savePhotoToFolder(photo *tb.Photo) (string, error) {
	fileURL, err := h.b.FileURLByID(photo.FileID)
	if err != nil {
		return "", err
	}

	resp, err := http.Get(fileURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	photoData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	imagesFolder := "images/"
	err = os.MkdirAll(imagesFolder, 0755)
	if err != nil {
		return "", err
	}

	photoPath := filepath.Join(imagesFolder, uuid.New().String()+"photo.jpg")
	err = ioutil.WriteFile(photoPath, photoData, 0644)
	if err != nil {
		return "", err
	}

	return photoPath, nil
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
