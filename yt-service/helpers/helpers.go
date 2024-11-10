package helpers

import (
	"context"
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
	// frames := []string{"🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘"}
	frames := []string{"🕛", "🕐", "🕑", "🕒", "🕓", "🕔", "🕕", "🕖", "🕗", "🕘", "🕙", "🕚"}
	// dots := []string{"➔", "➜", "➞", "→", "⇢"}
	// dots := []string{"🡪","🡮", "🡫","🡯", "🡨", "🡬", "🡩","🡭"}
	// dots := []string{"🟠", "🟡", "🟢", "🔵", "🟣", "⚫️"}
	dots := []string{"➣", " ➢", "  ➣", "   ➢", "    ➣", "     ➢", "      ➣","       ➢"}
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
	ticker := time.NewTicker(75 * time.Millisecond)

	defer ticker.Stop()

	// frames := []string{"🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘"}
	frames := []string{"🕛", "🕐", "🕑", "🕒", "🕓", "🕔", "🕕", "🕖", "🕗", "🕘", "🕙", "🕚"}
	dots := []string{"🡪","🡮", "🡫","🡯", "🡨", "🡬", "🡩","🡭"}
	// dots := []string{"➣", " ➢", "  ➣", "   ➢", "    ➣", "     ➢"}
	// dots := []string{"", ".", "..", "...", "...."}
	// dots := []string{"➛", "➔", "➜", "➞", "→", "⇢"}

	i, j := 0, 0

	for {
		select {
		case <-done:
			// Print the final elapsed time when done
			// fmt.Printf("\n\n\r⏲ Elapsed time to process video: %v\n", time.Since(startTime).Round(time.Second))
			return
		
		case <-ticker.C:
			// Display the elapsed time with changing frames and dots
			fmt.Printf("\r%s Merging%s 		Elapsed Time   %s %v", dots[j], dots[j],  frames[i], time.Since(startTime).Round(time.Second))

			// Update frames and dots
			i = (i + 1) % len(frames)
			j = (j + 1) % len(dots)
		}
	}
}



// MergeMedia runs the Python script to merge media, listening for context cancellation.
func MergeMedia(ctx context.Context, videoPath, audioPath, outputPath string) error {
	// Use exec.CommandContext to attach context to the command
	cmd := exec.CommandContext(ctx, "python", "./yt-service/helpers/python/merge_media.py", videoPath, audioPath, outputPath)
	cmd.Env = append(cmd.Env, "PYTHONIOENCODING=utf-8")

	// Capture the output
	output, err := cmd.CombinedOutput()

	if ctx.Err() == context.Canceled {
		// Context was canceled, log and exit gracefully
		fmt.Println("\nTask canceled, cleaning up resources, exits in 3 seconds")
		return ctx.Err()
	}
	if err != nil {
		fmt.Printf("\n\n *************** \n\n")
		return fmt.Errorf("error merging media: %v, output: %s", err, output)
	}
	return nil
}

func DeleteFromRoot(file string) (ok bool) {
	 // Define the common suffix to look for
	 suffix := "TEMP_MPY_wvf_snd.mp4"

	// Get the current working directory (root directory)
	//  rootDir, err := os.Getwd()
	//  if err != nil {
	// 	 fmt.Printf("Error getting root directory: %v\n", err)
	// 	 return
	//  }
 
	//  // Read the directory contents
	//  entries, err := os.ReadDir(rootDir)
	//  if err != nil {
	// 	 fmt.Printf("Error reading root directory: %v\n", err)
	// 	 return
	//  }
 
	//  fmt.Printf("Files ending with '%s':\n", suffix)
	// var fullPath string
	//  for _, entry := range entries {
	// 	 // Check if the entry is a file and if it ends with the specified suffix
	// 	 if !entry.IsDir() && strings.HasSuffix(entry.Name(), suffix) {
	// 		 fullPath = filepath.Join(rootDir, entry.Name())
	// 		 fmt.Println(fullPath)  // Print or process the full path as needed
	// 	 }
	//  }
	fileName := filepath.Base(file)
	fileName = strings.Split(fileName, ".")[0]
	fileName = fileName + suffix
	 	
	// Delete the file
	if err := os.Remove(fileName); err != nil {
		fmt.Printf("Temp file deletion failed: %v\n", err)
		ok = false
	} else {
		ok = true
	}
	return
}

func TermLog() {
	fmt.Printf("Temp files cleared\n\n")
			fmt.Println("\n		*--------------------------------------------------------------------------*")
			fmt.Println("		|                              BYE👋 BYE👋                                 |")
			fmt.Printf("		*--------------------------------------------------------------------------*\n\n\n")
}

func SuccessLog(opPath string, downloadTime time.Duration, now time.Time) {
	fmt.Printf("\n\n💾 video saved to➾ %s\n", opPath)
		fmt.Println("⏲ Elapsed time to download files: ", downloadTime)
		fmt.Printf("⏳ Total download and process duration: %v\n\n", time.Since(now))

		fmt.Println("\n		*--------------------------------------------------------------------------*")
		fmt.Println("		|                              BYE👋 BYE👋                                 |")
		fmt.Printf("		*--------------------------------------------------------------------------*\n\n\n")
		time.Sleep(3 * time.Second)
}