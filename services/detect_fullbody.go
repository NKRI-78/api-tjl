package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type DetectFullbodyResult struct {
	Message string         `json:"message"`
	Data    map[string]any `json:"data"`
}

type mediaUploadResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    struct {
		Path string `json:"path"`
	} `json:"data"`
}

type openAIChatResponse struct {
	Choices []struct {
		Message struct {
			ToolCalls []struct {
				Type     string `json:"type"`
				Function struct {
					Name      string `json:"name"`
					Arguments string `json:"arguments"`
				} `json:"function"`
			} `json:"tool_calls"`
		} `json:"message"`
	} `json:"choices"`
}

func DetectFullbody(
	folder string,
	file multipart.File,
	fileHeader *multipart.FileHeader,
) (*DetectFullbodyResult, error) {
	if folder == "" {
		return nil, errors.New("Field folder is required")
	}

	if fileHeader == nil {
		return nil, errors.New("Field media is required")
	}

	imageURL, err := uploadMediaToStorage(folder, file, fileHeader)
	if err != nil {
		return nil, err
	}

	base64Image, err := fetchImageAsBase64(imageURL)
	if err != nil {
		return nil, err
	}

	toolCalls, err := analyzeImageWithOpenAI(base64Image)
	if err != nil {
		return nil, err
	}

	if len(toolCalls) == 0 {
		return nil, errors.New("Unable to analyze image. Please try again with a clear standing photo")
	}

	fullbodyVisible := false
	isNude := false

	for _, tool := range toolCalls {
		var args map[string]bool

		err := json.Unmarshal([]byte(tool.Function.Arguments), &args)
		if err != nil {
			continue
		}

		if tool.Function.Name == "get_fullbody" {
			fullbodyVisible = args["humanbody"]
		}

		if tool.Function.Name == "check_nudity" {
			isNude = args["nude"]
		}
	}

	if !fullbodyVisible {
		return &DetectFullbodyResult{
			Message: "Full human body is not visible. Upload blocked.",
			Data: map[string]any{
				"fullbody_visible": false,
				"nude":             false,
			},
		}, nil
	}

	if isNude {
		return &DetectFullbodyResult{
			Message: "Image contains nudity or person is not standing naturally. Upload blocked.",
			Data: map[string]any{
				"fullbody_visible": false,
				"nude":             true,
			},
		}, nil
	}

	return &DetectFullbodyResult{
		Message: "",
		Data: map[string]any{
			"fullbody_visible": fullbodyVisible,
			"nude":             isNude,
		},
	}, nil
}

func uploadMediaToStorage(
	folder string,
	file multipart.File,
	fileHeader *multipart.FileHeader,
) (string, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Kalau mau sama persis seperti Node.js:
	// _ = writer.WriteField("folder", "ktp-scan")
	//
	// Kalau mau pakai request body `folder`, pakai ini:
	_ = writer.WriteField("folder", folder)

	originalName := strings.TrimSpace(fileHeader.Filename)
	ext := filepath.Ext(originalName)
	name := strings.TrimSuffix(originalName, ext)

	if name == "" {
		name = "media"
	}

	filename := fmt.Sprintf("%s_%d%s", name, time.Now().UnixMilli(), ext)

	part, err := writer.CreateFormFile("media", filename)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}

	err = writer.Close()
	if err != nil {
		return "", err
	}

	baseURLMedia := os.Getenv("BASE_URL_MEDIA")

	req, err := http.NewRequest(http.MethodPost, baseURLMedia, &body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", fmt.Errorf("failed upload media: %s", string(resBody))
	}

	var mediaRes mediaUploadResponse
	err = json.Unmarshal(resBody, &mediaRes)
	if err != nil {
		return "", err
	}

	if mediaRes.Data.Path == "" {
		return "", errors.New("media path is empty")
	}

	return mediaRes.Data.Path, nil
}

func fetchImageAsBase64(imageURL string) (string, error) {
	if imageURL == "" {
		return "", errors.New("image url is empty")
	}

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	res, err := client.Get(imageURL)
	if err != nil {
		return "", errors.New("Error fetching image")
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", errors.New("Error fetching image")
	}

	imageBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(imageBytes), nil
}

func analyzeImageWithOpenAI(base64Image string) ([]struct {
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}, error) {
	apiKey := os.Getenv("CHATGPT_API_KEY")

	payload := map[string]any{
		"model":       "gpt-4o",
		"temperature": 0.0,
		"tools": []map[string]any{
			{
				"type": "function",
				"function": map[string]any{
					"name":        "get_fullbody",
					"description": "Checks if a full human body is fully visible in the image.",
					"parameters": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"humanbody": map[string]any{
								"type":        "boolean",
								"description": "True if the full human body is visible from head to toe.",
							},
						},
						"required": []string{"humanbody"},
					},
				},
			},
			{
				"type": "function",
				"function": map[string]any{
					"name":        "check_nudity",
					"description": "Checks if the image contains nudity.",
					"parameters": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"nude": map[string]any{
								"type":        "boolean",
								"description": "True if any nudity is detected.",
							},
						},
						"required": []string{"nude"},
					},
				},
			},
		},
		"tool_choice": "auto",
		"messages": []map[string]any{
			{
				"role": "user",
				"content": []map[string]any{
					{
						"type": "text",
						"text": "Analyze this image and respond:\n" +
							"1) Using `get_fullbody`: Is the full human body visible from head to toe?\n" +
							"2) Using `check_nudity`: Is there any nudity content?\n",
					},
					{
						"type": "image_url",
						"image_url": map[string]any{
							"url":    "data:image/jpeg;base64," + base64Image,
							"detail": "high",
						},
					},
				},
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.openai.com/v1/chat/completions",
		bytes.NewBuffer(payloadBytes),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 90 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("openai error: %s", string(resBody))
	}

	var openAIRes openAIChatResponse
	err = json.Unmarshal(resBody, &openAIRes)
	if err != nil {
		return nil, err
	}

	if len(openAIRes.Choices) == 0 {
		return nil, errors.New("openai response choice is empty")
	}

	return openAIRes.Choices[0].Message.ToolCalls, nil
}
