package fetcher

import (
	"github.com/cavaliercoder/grab"
	"fmt"
	"os"
	"time"
	"github.com/fatih/color"
)

type Fetcher struct {}

//// get URL to download from command args
//if len(os.Args) < 2 {
//fmt.Fprintf(os.Stderr, "usage: %s url [url]...\n", os.Args[0])
//os.Exit(1)
//}
//
//urls := os.Args[1:]

func Fetch(files []string) {

	// start file downloads, 3 at a time
	fmt.Printf("Downloading %d files...\n", len(files))
	respch, err := grab.GetBatch(3, "./lib", files...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	// start a ticker to update progress every 200ms
	t := time.NewTicker(200 * time.Millisecond)
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	// monitor downloads
	completed := 0
	inProgress := 0
	responses := make([]*grab.Response, 0)
	for completed < len(files) {
		select {
		case resp := <-respch:
		// a new response has been received and has started downloading
		// (nil is received once, when the channel is closed by grab)
			if resp != nil {
				responses = append(responses, resp)
			}

		case <-t.C:
		// clear lines
			if inProgress > 0 {
				fmt.Printf("\033[%dA\033[K", inProgress)
			}

		// update completed downloads
			for i, resp := range responses {
				if resp != nil && resp.IsComplete() {
					// print final result
					if resp.Error != nil {
						fmt.Fprintf(os.Stderr, "Error downloading %s: %v\n", red(resp.Request.URL()), resp.Error)
					} else {
						fmt.Printf("Finished %s %v / %d bytes (%d%%)\n", yellow(resp.Filename), yellow(resp.BytesTransferred()), resp.Size, int(100*resp.Progress()))
					}

					// mark completed
					responses[i] = nil
					completed++
				}
			}

		// update downloads in progress
			inProgress = 0
			for _, resp := range responses {
				if resp != nil {
					inProgress++
					fmt.Printf("Downloading %s %v / %d bytes (%d%%)\033[K\n", yellow(resp.Filename), red(resp.BytesTransferred()), resp.Size, int(100*resp.Progress()))
				}
			}
		}
	}

	t.Stop()

	fmt.Printf("%d files successfully downloaded.\n", len(files))
}
