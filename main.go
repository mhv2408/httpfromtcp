package main

import (
	"fmt"
	"log"
	"os"
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
	for {
		read_data := make([]byte, 8)
		eof, err := file.Read(read_data)
		if eof == 0 {
			break
		}
		if err != nil {
			log.Fatal("Cannot read 8 bytes from file ", err)
		}

		fmt.Printf("read: %s\n", string(read_data))

	}
}
