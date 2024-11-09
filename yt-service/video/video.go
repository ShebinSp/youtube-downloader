package ytservice

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	ytAudio "github.com/ShebinSp/yt-downloader/yt-service/audio"
	"github.com/ShebinSp/yt-downloader/yt-service/helpers"
	"github.com/kkdai/youtube/v2"
)

type FileInfo struct {
	VideoName string
	VideoPath string
	AudioPath string
}

func DownloadYoutubeVideo(videoID string) (FileInfo, error) {
	var fileInfo FileInfo
	var wg sync.WaitGroup

	client := youtube.Client{}

	// Get the video details
	video, err := client.GetVideo(videoID)
	if err != nil {
		return FileInfo{}, fmt.Errorf("error fetching video: %v", err)
	}

	title := video.Title

	wg.Add(1)
	go func() {
		defer wg.Done()
		fileInfo.AudioPath, err = ytAudio.DownloadAudio(video)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
	}()
	if err != nil {
		fmt.Println("err: ", err)
	}

	filename := helpers.SanitizeFilename(title)
	filename = helpers.SanitizeFilename(filename)

	// Get the vidoe format (select the highest available quality)
	// var selectdFormats *youtube.Format

	// formats := video.Formats.Quality("1080p")
	// for _, format := range formats {
	// 	if format.QualityLabel == "2160p" {
	// 		selectdFormats = &format
	// 		break
	// 	} else if format.QualityLabel == "1440p" {
	// 		selectdFormats = &format
	// 		break
	// 	} else if format.QualityLabel == "1080p" {
	// 		selectdFormats = &format
	// 		break
	// 	} else if format.QualityLabel == "720p" {
	// 		selectdFormats = &format
	// 		break
	// 	} else {
	// 		selectdFormats = &formats[0]
	// 	}
	// }

	videoFormats := video.Formats.Type("video")
	videoFormats.Sort()

	stream, _, err := client.GetStream(video, &videoFormats[0])
	if err != nil {
		return FileInfo{}, fmt.Errorf("error getting video stream: %v", err)
	}

	wg.Wait()
	outputDir := "./yt-service/temp"

	// Ensure the output directory exists
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return FileInfo{}, fmt.Errorf("error creating output directory: %v", err)
	}

	// Ensure file has .mp4 ext
	if !strings.HasSuffix(filename, ".mp4") {
		filename += ".mp4"
	}
	fileInfo.VideoName = filename
	outputFilePath := filepath.Join(outputDir, filename)
	fileInfo.VideoPath = outputFilePath

	// Create the output file
	file, err := os.Create(outputFilePath)
	if err != nil {
		return FileInfo{}, fmt.Errorf("error while creating the file: %v", err)
	}
	defer file.Close()

	// Copy the stream to file
	_, err = io.Copy(file, stream)
	if err != nil {
		return FileInfo{}, fmt.Errorf("error while downloading the video: %v", err)
	}

	return fileInfo, nil
}
