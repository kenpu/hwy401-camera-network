package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func download(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("response error:" + resp.Status)
	}

	bufReader := new(bytes.Buffer)
	bufReader.ReadFrom(resp.Body)

	return bufReader.Bytes(), nil
}

func isDiff(buf1, buf2 []byte) bool {
	if buf1 == nil || buf2 == nil {
		return true
	}

	if len(buf1) != len(buf2) {
		return true
	}

	for i := range buf1 {
		if buf1[i] != buf2[i] {
			return true
		}
	}

	return false
}

func main() {
	var (
		loc      int
		output   string
		interval int
	)
	flag.IntVar(&loc, "location", 96, "Location index")
	flag.StringVar(&output, "out", "", "File to save to")
	flag.IntVar(&interval, "interval", 0, "Interval in seconds")

	flag.Parse()

	if output == "" {
		fmt.Println("No output defined")
		flag.PrintDefaults()
		return
	}

	url := fmt.Sprintf(
		"http://www.cdn.mto.gov.on.ca/"+
			"english/traveller/compass/camera/pictures/loc%d.jpg", loc)

	ticks := time.Tick(time.Duration(interval) * time.Second)
	var oldbuf []byte = nil
	var errCount = 0
	var i = 0
	for _ = range ticks {
		buf, err := download(url)
		if err != nil {
			fmt.Println("Error:", err)
			errCount += 1
			if errCount > 10 {
				fmt.Println("Bail out...")
				break
			}
		}

		if isDiff(oldbuf, buf) {
			out := fmt.Sprintf("%s_%.8d.jpg", output, i)
			fmt.Printf("Save %s\n", out)
			ioutil.WriteFile(out, buf, 0744)
			oldbuf = buf
			i += 1
		} else {
			fmt.Println("Same frame, skipping.")
		}
	}

	fmt.Println("All done.")
}
