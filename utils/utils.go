package utils

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

func CSVReader(filename string) [][]string {
	csvfile, err := os.Open(filename)
	if err != nil {
		log.Fatalf("open file fault, filename: %s, err: %v", filename, err)
	}
	file := csv.NewReader(csvfile)
	var res [][]string
	for {
		record, err := file.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("read csv fault, err: ", err)
		}
		res = append(res, record)
	}
	return res
}
