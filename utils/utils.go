package utils

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

/**
 * @description:  讀取csv
 * @param {string} filename
 * @return {*}
 */
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

/**
 * @description: 判斷是否為數字
 * @param {string} s
 * @return {*}
 */
func IsDigital(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
