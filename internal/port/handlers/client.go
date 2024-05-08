package handlers

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"
	"image"
	"iman_tg_bot/internal/model"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type File struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

var (
	clientData model.Client
	lat        float32
	long       float32
)

func (h *Handler) AddClient(m *tb.Message) {
	steps := []struct {
		prompt  string
		setter  func(string) error
		hasNext bool
	}{
		{"Contract ID", func(text string) error { clientData.ContractId = text; return nil }, true},
		{"Phone Number", func(text string) error { clientData.PhoneNumber = text; return nil }, true},
		{"Address", func(text string) error { clientData.Address = text; return nil }, true},
		{"Payment sum", func(text string) error { clientData.PaymentSum = text; return nil }, true},
		{"Comment", func(text string) error { clientData.Comment = text; return nil }, true},
		{"Location", func(text string) error { clientData.Location = text; return nil }, false},
		{"Address foto", func(text string) error { clientData.AddressFoto = text; return nil }, true},
		{"Payment foto", func(text string) error { clientData.PaymentFoto = text; return nil }, true},
	}

	h.handleSteps(m, &clientData, steps, 0)
}

func (h *Handler) handleSteps(m *tb.Message, clientData *model.Client, steps []struct {
	prompt  string
	setter  func(string) error
	hasNext bool
}, currentStep int) {
	if currentStep >= len(steps) {
		h.finalizeAddClient(m)
		return
	}

	step := steps[currentStep]
	h.b.Send(m.Sender, step.prompt)

	h.b.Handle(tb.OnText, func(msg *tb.Message) {
		if err := step.setter(msg.Text); err != nil {
			log.Println("Error setting client data:", err)
			return
		}
		currentStep++
		h.handleSteps(m, clientData, steps, currentStep)
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

		if clientData.AddressFoto == "" {
			clientData.AddressFoto = photoPath
		} else if clientData.PaymentFoto == "" {
			clientData.PaymentFoto = photoPath
		}

		currentStep++
		h.handleSteps(m, clientData, steps, currentStep)
	})

	h.b.Handle(tb.OnLocation, func(msg *tb.Message) {
		// Check if the message contains location data
		if msg.Location == nil {
			log.Println("No location data found in the message")
			return
		}

		// Extract latitude and longitude from the message
		lat = msg.Location.Lat
		long = msg.Location.Lng

		clientData.LocationLatitude = fmt.Sprintf("%f", lat)
		clientData.LocationLongitude = fmt.Sprintf("%f", long)

		currentStep++
		h.handleSteps(m, clientData, steps, currentStep)
	})
}

func (h *Handler) finalizeAddClient(m *tb.Message) {
	// Generate the file name
	fileName := h.generateFileName(m.Sender, clientData.ContractId)

	// Create a new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 15) // Decrease font size for better text arrangement

	// Add content to the PDF
	pdf.SetX(10) // Set initial X position
	pdf.CellFormat(40, 10, "Contract ID:", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 10, clientData.ContractId, "", 0, "L", false, 0, "")
	pdf.Ln(-1) // Move to the next line without spacing

	pdf.SetX(10) // Reset X position
	pdf.CellFormat(40, 10, "Phone Number:", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 10, clientData.PhoneNumber, "", 0, "L", false, 0, "")
	pdf.Ln(-1) // Move to the next line without spacing

	pdf.SetX(10) // Reset X position
	pdf.CellFormat(40, 10, "Address:", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 10, clientData.Address, "", 0, "L", false, 0, "")
	pdf.Ln(-1) // Move to the next line without spacing

	pdf.SetX(10) // Reset X position
	pdf.CellFormat(40, 10, "Payment sum:", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 10, clientData.PaymentSum, "", 0, "L", false, 0, "")
	pdf.Ln(-1)

	pdf.SetX(10)
	pdf.CellFormat(40, 10, "Comment:", "", 0, "L", false, 0, "")
	pdf.MultiCell(0, 10, clientData.Comment, "", "L", false) // Use MultiCell for comment
	pdf.Ln(-1)

	pdf.SetX(10)
	pdf.CellFormat(40, 10, "Location:", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 10, "Open in Maps", "", 0, "L", false, 0, "")
	pdf.Ln(-1)

	// Add a clickable link to the location
	pdf.SetX(10)
	pdf.WriteLinkString(0, fmt.Sprintf("https://www.google.com/maps?q=%f,%f", lat, long), "")
	pdf.Ln(7)

	// Get image dimensions
	imageWidth, imageHeight := getImageDimensions(clientData.AddressFoto)

	// Calculate the aspect ratio for resizing
	imageAspectRatio := imageWidth / imageHeight
	maxWidth := 140.0
	maxHeight := maxWidth / imageAspectRatio

	pdf.Cell(0, 10, "Address Photo:")
	if clientData.AddressFoto != "" {
		pdf.ImageOptions(clientData.AddressFoto, 10.0, pdf.GetY()+10, maxWidth, maxHeight, false, gofpdf.ImageOptions{}, 0, "")
	}

	pdf.AddPage()
	pdf.Cell(0, 10, "Payment foto:")
	if clientData.PaymentFoto != "" {
		pdf.ImageOptions(clientData.PaymentFoto, 10.0, pdf.GetY()+10, maxWidth, maxHeight, false, gofpdf.ImageOptions{}, 0, "")
	}

	// Write the PDF content to a buffer
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		log.Println("Error saving PDF to buffer: ", err)
		h.b.Send(m.Sender, "Error saving PDF")
		return
	}

	if err := h.repo.ClientUser().Create(&clientData); err != nil {
		h.log.Error("failed to create client user", zap.Error(err))
		return
	}

	// Send the PDF document to the user
	h.b.Send(m.Sender, &tb.Document{File: tb.FromReader(&buf), FileName: fileName}, &tb.SendOptions{
		ReplyMarkup: AdminButtons,
	})
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

func (h *Handler) generateFileName(sender *tb.User, rasrochqaId string) string {
	fileName := "client_information.pdf"
	if sender != nil && (sender.FirstName != "" || sender.LastName != "") {
		fileName = fmt.Sprintf("%s_%s_%s.pdf", sender.FirstName, sender.LastName, rasrochqaId)
	}
	return fileName
}

// Method to save the received photo file to the images folder
func (h *Handler) savePhotoToFolder(photo *tb.Photo) (string, error) {
	// Get the file URL
	fileURL, err := h.b.FileURLByID(photo.FileID)
	if err != nil {
		return "", err
	}

	// Download the photo data
	resp, err := http.Get(fileURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the photo data
	photoData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Create the images folder if it doesn't exist
	imagesFolder := "images/"
	err = os.MkdirAll(imagesFolder, 0755)
	if err != nil {
		return "", err
	}

	// Save the photo data to a file in the images folder
	photoPath := filepath.Join(imagesFolder, uuid.New().String()+"photo.jpg")
	err = ioutil.WriteFile(photoPath, photoData, 0644)
	if err != nil {
		return "", err
	}

	return photoPath, nil
}
