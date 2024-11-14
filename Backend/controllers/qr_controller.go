package controllers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"organize-this/helpers"
	"organize-this/models"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

// Generate creates a new QR code, uploads it to S3, and then returns the url to the frontend.
func (handler Handler) Generate(w http.ResponseWriter, request *http.Request) {
	//Get parameters
	byteData, err := io.ReadAll(request.Body)
	if err != nil {
		logAndRespond(w, "Error parsing request", err)
		return
	}

	var parsedData map[string]string
	if err = json.Unmarshal(byteData, &parsedData); err != nil {
		logAndRespond(w, "Error parsing json", err)
		return
	}

	//Validate parameters
	category, stringID := parsedData["category"], parsedData["id"]
	if category == "" {
		logAndRespond(w, "Missing category", nil)
		return
	}

	if stringID == "" {
		logAndRespond(w, "Missing id", nil)
		return
	}

	id, err := strconv.ParseUint(stringID, 10, 64)
	if err != nil {
		logAndRespond(w, fmt.Sprintf("ID must be type integer: %v", stringID), nil)
		return
	}

	claims := request.Context().Value("user_claims").(jwt.MapClaims)
	userID := claims["username"].(string)

	// Verify entity exists
	entity := models.Entity{
		ID: id,
	}

	validEntity, model := buildEntity(entity, models.Parent{}, category, "")
	if !validEntity {
		logAndRespond(w, fmt.Sprintf("Invalid category %v.", category), nil)
		return
	}

	dberr := handler.Repository.GetOne(model, userID)
	if dberr != nil {
		logAndRespond(w, fmt.Sprintf("Entity category of %v with id %v not found.", category, stringID), nil)
		return
	}

	// Check if object exists
	bucketName := "organize-this-local"
	fileName := base64.StdEncoding.EncodeToString([]byte(userID+"_QR_"+category+"_"+stringID)) + ".jpeg"
	fileURL := "https://organize-this-local.s3.us-east-2.amazonaws.com/" + fileName

	_, err = handler.S3Client.HeadObject(request.Context(), &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	})

	if err != nil {
		var notFound *types.NotFound
		if errors.As(err, &notFound) {
			helpers.SuccessResponse(w, fileURL)
			return
		}
	}

	// Build QR
	url := "http://localhost:5173/" + category + "/" + stringID
	fileLocation := "assets/" + fileName

	qrc, err := qrcode.New(url)
	if err != nil {
		logAndRespond(w, "could not generate QRCode: %v", err)
		return
	}

	writer, err := standard.New(fileLocation)
	if err != nil {
		logAndRespond(w, "standard.New failed: %v", err)
		return
	}

	// Save TMP file
	if err = qrc.Save(writer); err != nil {
		logAndRespond(w, "could not save image: %v", err)
		return
	}

	// Upload to S3
	file, err := os.Open(fileLocation)
	if err != nil {
		logAndRespond(w, "Couldn't open file: %v", err)
		return
	}

	defer file.Close()
	_, err = handler.S3Client.PutObject(request.Context(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   file,
	})

	if err != nil {
		logAndRespond(w, fmt.Sprintf("Couldn't upload file: %v\n", err), err)
		return
	}

	err = s3.NewObjectExistsWaiter(handler.S3Client).Wait(
		request.Context(), &s3.HeadObjectInput{Bucket: aws.String(bucketName), Key: aws.String(fileName)}, time.Minute)
	if err != nil {
		logAndRespond(w, fmt.Sprintf("Failed attempt to wait for object %s to exist.\n", fileName), err)
		return
	}

	// Clean Up TMP file
	e := os.Remove(fileLocation)
	if e != nil {
		logAndRespond(w, "could not remove tmp image: %v", e)
		return
	}

	helpers.SuccessResponse(w, fileURL)
	return
}
