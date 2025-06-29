package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const inputFile = "message.txt"

func main() {

	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal("Cannot open the file: ", err)
	}
	defer file.Close()
	fmt.Printf("Reading data from %s\n", inputFile)
	fmt.Println("=====================================")
	// curr_output := ""
	/*
		for {
			read_data := make([]byte, 8)
			eof, err := file.Read(read_data)
			if eof == 0 {
				break
			}
			if err != nil {
				log.Fatal("Cannot read 8 bytes from file ", err)
			}
			str := string(read_data[:eof])
			parts := strings.Split(str, "\n")

			curr_output += parts[0]

			if len(parts) == 1 {
				continue
			}

			fmt.Printf("read: %s\n", curr_output)
			curr_output = parts[1]

		}
		fmt.Printf("read: %s\n", curr_output) // printing the last end */

	currentLineContents := ""
	for {
		buffer := make([]byte, 8, 8)
		n, err := file.Read(buffer)
		if err != nil {
			if currentLineContents != "" {
				fmt.Printf("read: %s\n", currentLineContents)
				currentLineContents = ""
			}
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Printf("error: %s\n", err.Error())
			break
		}
		str := string(buffer[:n])
		parts := strings.Split(str, "\n")
		for i := 0; i < len(parts)-1; i++ {
			fmt.Printf("read: %s%s\n", currentLineContents, parts[i])
			currentLineContents = ""
		}
		currentLineContents += parts[len(parts)-1]
	}
}
