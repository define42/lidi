package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse the multipart form, with a maximum of 10MB file size.
		err := r.ParseMultipartForm(10 << 28)
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		// Retrieve the file from form data
		file, _, err := r.FormFile("uploadFile")
		if err != nil {
			http.Error(w, "Error retrieving file", http.StatusBadRequest)
			return
		}
		defer file.Close()
		/*
			// Save the file to a temporary location
			tempFile, err := os.Create(filepath.Join(os.TempDir(), handler.Filename))
			if err != nil {
				http.Error(w, "Error saving file", http.StatusInternalServerError)
				return
			}
			defer tempFile.Close()

			// Copy the uploaded file to the temporary file
			_, err = io.Copy(tempFile, file)
			if err != nil {
				http.Error(w, "Error saving file", http.StatusInternalServerError)
				return
			}
		*/
		// Forward the file to the TCP server
		err = forwardToTCPServer(file)
		if err != nil {
			http.Error(w, "Error forwarding file to TCP server", http.StatusInternalServerError)
			return
		}

		// Send a response to the client
		fmt.Fprintf(w, "File uploaded successfully!")
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func forwardToTCPServer(data io.Reader) error {
	// Connect to the TCP server
	conn, err := net.Dial("tcp", "172.16.0.2:5000")
	if err != nil {
		return fmt.Errorf("error connecting to TCP server: %v", err)
	}
	defer conn.Close()
	/*
		// Open the file
		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("error opening file: %v", err)
		}
		defer file.Close()
	*/
	// Send the file contents to the TCP server
	_, err = io.Copy(conn, data)
	if err != nil {
		return fmt.Errorf("error sending file to TCP server: %v", err)
	}

	return nil
}

func main() {
	// Serve the file upload form
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<html>
<body>
<form enctype="multipart/form-data" action="/upload" method="post">
    <input type="file" name="uploadFile" />
    <input type="submit" value="upload" />
</form>
</body>
</html>`)
	})

	// Handle the file upload
	http.HandleFunc("/upload", uploadHandler)

	// Start the HTTP server
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
