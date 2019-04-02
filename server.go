package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"log"
	"net/http"
	"strings"
)

// Server global variables
var port string
var bucket string
var region string

// Upload data to AWS S3
func uploadToS3(fileName string, content *bytes.Buffer) (string, error) {
	newSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return "failure", err
	}

	uploader := s3manager.NewUploader(newSession)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
		Body:   content,
	})
	if err != nil {
		return "failure", errors.New(fmt.Sprintf("Unable to upload %q to %q, %v", fileName, bucket, err))
	}

	return "success", nil

}

// If a file path looks like it has directories, create them
// Returns the path that was created, and the file that was in the path
func getDirectoryPathFromFilePath(filepath string) (string, string) {
	directoryPath := ""

	// Split the file path into paths
	fnSlice := strings.Split(filepath, "/")
	fnSliceLen := len(fnSlice)

	if fnSliceLen > 1 {
		// Remove the filename at the end
		dirSlice := fnSlice[:fnSliceLen-1]

		// Build a path containing only the directories
		directoryPath = strings.Join(dirSlice, "/")
	}

	// Get the name of the file
	fileName := strings.Split(fnSlice[fnSliceLen-1:][0], ".")[0]

	return directoryPath, fileName
}

// Send a response back
func sendResponse(responseWriter http.ResponseWriter, statusCode int, message string) {
	responseWriter.WriteHeader(statusCode)
	fmt.Fprintf(responseWriter, message)
}

type Metadata struct {
	S3Root string `json:"S3Root"`
}

// Handle HTTP Requests
func handleRequest(responseWriter http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		// Information about what the server does.
		sendResponse(
			responseWriter,
			http.StatusOK,
			"This server sends files data to AWS S3. Note: Only multipart/form-data is accepted.",
		)
		return

	case "POST":
		mr, err := request.MultipartReader()

		// Didn't get multipart/form-data, return a 415 status code
		if err != nil {
			sendResponse(
				responseWriter,
				http.StatusUnsupportedMediaType,
				fmt.Sprintln(err),
			)
			return
		}

		// Record of files created during this POST request
		var filenames []string

		var fileName string

		// The filename param in a request part
		var fileLabel string

		// If the filename looks like it contains a directory tree, store it here
		var directoryPath string = ""

		var metadata = Metadata{}
		var metadataContent *bytes.Buffer
		var fileContent *bytes.Buffer

		// Parse through each part of the request
		for {
			part, err := mr.NextPart()
			if err == io.EOF {
				break
			}

			if part != nil {
				buf := new(bytes.Buffer)
				buf.ReadFrom(part)

				// Each request could have a part with metadata in JSON format
				if part.FormName() == "metadata" {
					metadataContent = buf

					// Make metadata available as an object
					json.Unmarshal([]byte(buf.String()), &metadata)

					// If the part isn't metadata, assume it's a file
				} else {
					// Get the filename from the request
					fileName = part.FileName()

					// If a filename wasn't found and the part isn't metadata, then the data is completely unexpected
					if fileName == "" {
						sendResponse(
							responseWriter,
							http.StatusUnsupportedMediaType,
							fmt.Sprintf("Unexpected data: %q", buf.String()),
						)
						return
					}

					fileContent = buf

					// Create the directory structure
					directoryPath, fileLabel = getDirectoryPathFromFilePath(fileName)
				}
			}
		}

		// Handle the _metadata.log file location
		// If the directory path isn't blank, add it to the metadata path
		metadataPath := directoryPath
		if metadataPath != "" {
			metadataPath = directoryPath + "/"
		}

		// Create a corrected path and name for the files
		metadataFullPath := metadataPath + fileLabel + "_metadata.log"
		fileFullPath := fileName

		// If S3Root was provided in the metadata, add as a prefix to path
		if metadata.S3Root != "" {
			metadataFullPath = metadata.S3Root + "/" + metadataFullPath
			fileFullPath = metadata.S3Root + "/" + fileName
		}

		// Add file paths to a list of filenames. (This will be returned in the response body.)
		filenames = append(filenames, metadataFullPath)
		filenames = append(filenames, fileFullPath)

		// Send the metadata to an S3 bucket
		_, err = uploadToS3(metadataFullPath, metadataContent)
		if err != nil {
			sendResponse(
				responseWriter,
				http.StatusUnprocessableEntity,
				fmt.Sprintln(err),
			)
			return
		}

		// Send the file content to an S3 bucket
		_, err = uploadToS3(fileFullPath, fileContent)
		if err != nil {
			sendResponse(
				responseWriter,
				http.StatusUnprocessableEntity,
				fmt.Sprintln(err),
			)
			return
		}

		// Send back a response with files created
		responsePayload, _ := json.Marshal(filenames)
		fmt.Fprintf(responseWriter, string(responsePayload))

	default:
		sendResponse(
			responseWriter,
			http.StatusMethodNotAllowed,
			"Only GET and POST methods are supported.",
		)
	}
}

func main() {
	flag.StringVar(&port, "port", "8081", "Port to launch server on")
	flag.StringVar(&bucket, "bucket", "", "S3 bucket name to send files to")
	flag.StringVar(&region, "region", "us-east-1", "AWS region")
	flag.Parse()

	http.HandleFunc("/", handleRequest)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
