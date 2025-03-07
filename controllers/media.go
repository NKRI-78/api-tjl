package controllers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	helper "superapps/helpers"
)

type MediaResponseEntity struct {
	Path     string `json:"path"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	Mime     string `json:"mime"`
}

func Upload(w http.ResponseWriter, r *http.Request) {

	url := os.Getenv("API_URL")

	err := r.ParseMultipartForm(10 << 20) // Max 10 MB
	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}
	defer file.Close()

	parentFolder := "public"
	subFolder := r.FormValue("folder")

	// Create Parent Folder if it doesn't exist
	errParentMkdir := os.MkdirAll(parentFolder, os.ModePerm)
	if errParentMkdir != nil {
		helper.Response(w, 400, true, "Failed to create parent folder: "+errParentMkdir.Error(), map[string]any{})
		return
	}

	// Construct the full path for the subfolder
	subFolderPath := filepath.Join(parentFolder, subFolder)

	// Create Subfolder if it doesn't exist (use MkdirAll to avoid error if it already exists)
	errSubFolder := os.MkdirAll(subFolderPath, os.ModePerm)
	if errSubFolder != nil {
		helper.Response(w, 400, true, "Failed to create subfolder: "+errSubFolder.Error(), map[string]any{})
		return
	}

	// Construct the full file path for saving the file
	filePath := filepath.Join(subFolderPath, handler.Filename)

	// Create and open the file for writing (O_WRONLY) and create it if it doesn't exist (O_CREATE)
	f, err := os.OpenFile(filepath.Clean(filePath), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		helper.Response(w, 400, true, "Failed to create or open the file: "+err.Error(), map[string]any{})
		return
	}
	defer f.Close()

	_, errCopy := io.Copy(f, file)

	if errCopy != nil {
		helper.Response(w, 400, true, errCopy.Error(), map[string]any{})
		return
	}

	var mediaAssign MediaResponseEntity

	mediaAssign.Path = url + "/" + subFolder + "/" + handler.Filename
	mediaAssign.Filename = handler.Filename
	mediaAssign.Size = handler.Size
	mediaAssign.Mime = handler.Header.Get("Content-Type")

	helper.Logger("info", "Upload success")
	helper.Response(w, http.StatusOK, false, "Successfully", mediaAssign)
}
