package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"time"

	"github.com/ShebinSp/yt-downloader/yt-service/helpers"
	ytservice "github.com/ShebinSp/yt-downloader/yt-service/video"
)

func main() {

	fmt.Println("\n     🎞🎞🎞🎞🎞🎞🎞 YouTube Downloader 🎞🎞🎞🎞🎞🎞🎞")
	fmt.Printf("\n🔴 use this application for educational purpose only! 🔴\n")
	fmt.Printf("🛑 Press CTRL + C to stop the program\n\n")

	var videoID string
	now := time.Now()
	pathCh := make(chan helpers.Path, 1)
	downloadDone := make(chan bool)
	mergeDone := make(chan bool)

	wg := &sync.WaitGroup{}
	defer close(pathCh)
	// ctx, cancel := context.WithCancel(context.Background())
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	fmt.Println("⌨️ Enter the YouTube video URL to download⪼")
	fmt.Print("⪼ ")
	fmt.Scanf("%s", &videoID)
	fmt.Println()

	// Get the Download folder of the user
	go helpers.GetDownlaodFolder(pathCh)

	// Start the spinner and Download the video and audio files
	wg.Add(1)
	go helpers.ShowSpinner(downloadDone, wg)
	fileInfo, err := ytservice.DownloadYoutubeVideo(videoID)
	if err != nil {
		log.Fatalf("Failed to download the video: %v\n", err)
	}

	// Stop the spinner and save the download time
	downloadDone <- true
	downloadTime := time.Since(now)

	// Start merging with elapsed time display
	fmt.Println("\nMergeing video and audio...")
	wg.Add(1)
	go helpers.ShowElapsedTime(mergeDone, wg)

	// Getting the file path ie, Download folder from buffered channel
	path := <-pathCh
	opPath := path.Path
	if path.Err != nil {
		log.Println("Failed to find the Download folder")
		opPath = filepath.Join("C:/")
	}

	// Creating the file name with complete path to Download folder
	outputPath := filepath.Join(opPath, fileInfo.VideoName)

	// Call to merge media with ctx, video path, audio path and output path
	err = helpers.MergeMedia(ctx, fileInfo.VideoPath, fileInfo.AudioPath, outputPath)
	if err != nil && err != context.Canceled {
		// If any errors from the python or cmd.CombinedOutput()
		mergeDone <- true
		helpers.ClearTemp(fileInfo.VideoPath, fileInfo.AudioPath)
		log.Printf("Failed to merge the video: %v\n", err)
		return
	} else if err == context.Canceled {
		// In case of program termination from the user
		mergeDone <- true
		time.Sleep(3 * time.Second)
		helpers.ClearTemp(fileInfo.VideoPath, fileInfo.AudioPath)
		ok := helpers.DeleteFromRoot(fileInfo.VideoPath)
		if ok {
			helpers.TermLog()
		} else {
			os.Exit(1)
		}
	} else {
		mergeDone <- true
		helpers.SuccessLog(opPath, downloadTime, now)
		helpers.ClearTemp(fileInfo.VideoPath, fileInfo.AudioPath)
	}

	wg.Wait()
}
