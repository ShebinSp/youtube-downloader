package helpers

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Path struct {
	Path string
	Err  error
}

func GetDownlaodFolder(pathC chan Path) {
	// Get the current user
	user, err := user.Current()
	if err != nil {
		pathC <- Path{Path: "", Err: err}
	}

	// Build the path to the Downloads folder(Windows and UNIX systems)
	downloadPath := filepath.Join(user.HomeDir, "Downloads")
	pathC <- Path{Path: downloadPath, Err: nil}

}

func SanitizeFilename(filename string) string {
	// Replace invalid characters with underscores
	invalidChars := []string{"<", ">", ":", "\"", "/", "\\", "|", "?", "*"}

	for _, char := range invalidChars {
		filename = strings.ReplaceAll(filename, char, "_")
		filename = strings.Trim(filename, " ")
	}
	return filename
}

func ClearTemp(vFile, aFile string) {
	err := os.Remove(vFile)
	if err != nil {
		fmt.Printf("temp deletion failed - vidoe: %v", err)
	}
	err = os.Remove(aFile)
	if err != nil {
		fmt.Printf("temp deletion failed - audio: %v", err)
	}
}

func ShowSpinner(done chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	frames := []string{"🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘"}
	dots := []string{"➛", "➔", "➜", "➞", "→", "⇢"}
	i, j := 0, 0

	for {
		select {
		case <-done:
			fmt.Print("\r✅ Download Complete!        \n")
			return
		default:
			// Display the spinner with changing dots
			fmt.Printf("\r%s Downloading%s", frames[i], dots[j])

			// Update frames and dots
			i = (i + 1) % len(frames)
			j = (j + 1) % len(dots)

			time.Sleep(150 * time.Millisecond)
		}
	}
}

func ShowElapsedTime(done chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	startTime := time.Now()
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	frames := []string{"🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘"}
	// dots := []string{"", ".", "..", "...", "...."}
	dots := []string{"➛", "➔", "➜", "➞", "→", "⇢"}

	i, j := 0, 0

	for {
		select {
		case <-done:
			// Print the final elapsed time when done
			fmt.Printf("\n\n\r⏲ Elapsed time to process video: %v\n", time.Since(startTime).Round(time.Second))
			return
		case <-ticker.C:
			// Display the elapsed time with changing frames and dots
			fmt.Printf("\r%s Merging%s  Elapsed Time: %v", frames[i], dots[j], time.Since(startTime).Round(time.Second))

			// Update frames and dots
			i = (i + 1) % len(frames)
			j = (j + 1) % len(dots)
		}
	}
}

func MergeMedia(videoPath, audioPath, outputPath string) error {
	// The media merging with the executable build by Nuitka
	// cmd := exec.Command("./merge_media.dist/merge_media.exe", videoPath, audioPath, outputPath)

	// The media merging using python script
	cmd := exec.Command("python", "./yt-service/helpers/python/merge_media.py", videoPath, audioPath, outputPath)

	// Set TEMP environment variable for PyInstaller
	cmd.Env = append(cmd.Env, "PYTHONIOENCODING=utf-8")

	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("error merging media: %v, output: %s", err, output)
	}
	
	return nil
}
