package ytservice

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ShebinSp/yt-downloader/yt-service/helpers"
	"github.com/kkdai/youtube/v2"
)

func DownloadAudio(video *youtube.Video) (string, error) {
	client := youtube.Client{}

	// Filter formats to get audio-only formats
	audioFormats := video.Formats.Type("audio")

	// Sort audio formats by bitrate to get the best quality
	if len(audioFormats) == 0 {
		return "", fmt.Errorf("no audio-only formats available")
	}

	audioFormats.Sort()

	// Select the best format available aft
	audioFormat := audioFormats[0]

	title := helpers.SanitizeFilename(video.Title)

	// Get the audio stream
	stream, _, err := client.GetStream(video, &audioFormat)
	if err != nil {
		return "", fmt.Errorf("error while geting the audio: %v", err)
	}

	outputDir := "./yt-service/temp/"

	// Create output file
	outputPath := filepath.Join(outputDir + title + ".m4a")

	file, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("error while creating file: %v", err)
	}
	defer file.Close()

	// Write the stream to the file
	_, err = io.Copy(file, stream)
	if err != nil {
		return "", fmt.Errorf("error while downloading file: %v", err)
	}

	return outputPath, nil
}
