package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
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
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	wg := &sync.WaitGroup{}
	defer close(pathCh)
	defer close(sig)

	fmt.Println("⌨️ Enter the YouTube video url to download⟹")
	fmt.Scanf("%s", &videoID)
	fmt.Println()

	go helpers.GetDownlaodFolder(pathCh)

	wg.Add(1)
	go helpers.ShowSpinner(downloadDone, wg)
	fileInfo, err := ytservice.DownloadYoutubeVideo(videoID)
	if err != nil {
		log.Fatalf("Failed to download the video: %v\n", err)
	}
	go func() {
		<-sig
		fmt.Println("Signal received, Stoping the program...")
		helpers.ClearTemp(fileInfo.VideoPath, fileInfo.AudioPath)
		os.Exit(0)
	}()
	downloadDone <- true
	downloadTime := time.Since(now)

	// Start merging with elapsed time display
	fmt.Println("\nMergeing video and audio...")
	wg.Add(1)
	go helpers.ShowElapsedTime(mergeDone, wg)

	path := <-pathCh
	opPath := path.Path
	if path.Err != nil {
		log.Println("Failed to save file to Download folder")
		opPath = filepath.Join("/C/")
	}

	outputPath := filepath.Join(opPath, fileInfo.VideoName)

	if err != nil {
		log.Printf("Sorry: %v\n", err)
	}

	err = helpers.MergeMedia(fileInfo.VideoPath, fileInfo.AudioPath, outputPath)
	if err != nil {
		helpers.ClearTemp(fileInfo.VideoPath, fileInfo.AudioPath)
		log.Printf("Failed to merge the video: %v\n", err)
		return
	}
	mergeDone <- true
	wg.Wait()
	defer helpers.ClearTemp(fileInfo.VideoPath, fileInfo.AudioPath)

	fmt.Printf("\n\n💾video saved successfully to %s🏁\n", opPath)
	fmt.Println("⏲ Elapsed time to download files: ", downloadTime)
	fmt.Printf("⏳ Total download and process duration: %v\n\n", time.Since(now))

	fmt.Println("\n*--------------------------------------------------------------------------*")
	fmt.Println("|                              BYE👋 BYE👋                                 |")
	fmt.Printf("*--------------------------------------------------------------------------*\n\n\n")
	time.Sleep(3 * time.Second)
}
