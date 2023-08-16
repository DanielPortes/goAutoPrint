package main

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"github.com/kbinani/screenshot"
	"image/jpeg"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

var (
	interrupt = make(chan os.Signal, 1)
)

func main() {
	var timer, repetition int
	paused := false

	signal.Notify(interrupt, os.Interrupt)

	getParametersFromUser(&timer, &repetition)
	fmt.Printf("Timer set to %d seconds\n Repeation set to %d times\n", timer, repetition)

	fmt.Println("Press PgDn to pause or resume capturing...")
	go watchForPause(&paused) // Start watching for PgDn key press in the background

	fmt.Println("Starting capturing...")
	captureLoop(&paused, repetition, timer)

	fmt.Println("Finished capturing.")
	close(interrupt)
}

func captureLoop(paused *bool, repetition int, timer int) {
	for {
		if !*paused {
			captureScreen()
			// in case user choose to not repeat forever
			if repetition > 0 {
				repetition--
				if repetition == 0 {
					break
				}
			}
			time.Sleep(time.Duration(timer) * time.Second)
		}
	}
}

func watchForPause(paused *bool) {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	fmt.Println("Press PgDn to pause and resume capturing...")
	for {
		_, keypress, err := keyboard.GetKey()
		if err != nil {
			fmt.Println("Error reading key:", err)
			continue
		}

		if keypress == keyboard.KeyPgdn {
			*paused = !(*paused)
			if *paused {
				fmt.Println("Capturing paused.")
			} else {
				fmt.Println("Capturing resumed.")
			}
		}
	}
}

func getParametersFromUser(i *int, i2 *int) {
	fmt.Println("Enter timer in seconds and how many repetition to repeat(-1 for infinity):")
	fmt.Println("Example: 5 -1")
	_, err := fmt.Scanf("%d %d", i, i2)
	if err != nil {
		fmt.Println(err)
	}
}

func captureScreen() {
	numDisplays := screenshot.NumActiveDisplays()
	for i := 0; i < numDisplays; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			panic(err)
		}

		dirPath := "./screenshots/"
		filename := fmt.Sprintf("screenshot_%s.jpg", time.Now().Format("2006-01-02 15-04-05"))
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			panic(err)
		}

		file, err := os.Create(filepath.Join(dirPath, filename))
		if err != nil {
			panic(err)
		}

		quality := &jpeg.Options{Quality: 50}
		err = jpeg.Encode(file, img, quality)
		if err != nil {

			panic(err)
		}

		err = file.Close()
		if err != nil {
			panic(err)
		}

		fmt.Printf("Print screen saved to %s\n", filename)
	}
}
