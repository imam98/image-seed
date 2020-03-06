package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"sync"

	"gocv.io/x/gocv"
)

func main() {
	webcam, err := gocv.OpenVideoCapture(0)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer webcam.Close()

	wg := sync.WaitGroup{}
	for {
		img := gocv.NewMat()
		if ok := webcam.Read(&img); !ok || img.Empty() {
			log.Println("Error writing image")
			img.Close()
			continue
		}

		imgChannels := gocv.Split(img)

		var pWrite []*io.PipeWriter
		var pRead []*io.PipeReader
		for i := 0; i < 3; i++ {
			pr, pw := io.Pipe()
			pWrite = append(pWrite, pw)
			pRead = append(pRead, pr)
		}

		for i := 0; i < 3; i++ {
			br := bytes.NewReader(imgChannels[i].ToBytes())
			go func(id int) {
				io.Copy(pWrite[id], br)
				pWrite[id].Close()
			}(i)

			go func(id int) {
				wg.Add(1)
				stream(pRead[id], id)
				wg.Done()
			}(i)
		}

		wg.Wait()
		img.Close()
	}
}

func stream(pipeReader *io.PipeReader, id int) {
	defer pipeReader.Close()

	file, err := os.OpenFile("/dev/urandom", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	log.Println("Goroutine", id, "to file...")
	if _, err := io.Copy(file, pipeReader); err != nil {
		log.Println(err)
		return
	}
	log.Println("Goroutine", id, "finished")
}
