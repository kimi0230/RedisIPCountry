package utils

import (
	"encoding/csv"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"
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

/**
 * @description: unit test 錯誤時跳error
 * @param {*testing.T} t
 * @param {bool} v
 * @return {*}
 */
func AssertTrue(t *testing.T, v bool) {
	t.Helper()
	if !v {
		t.Error("assert false but get a true value")
	}
}

func RandomString(up int) string {
	rand.Seed(rand.Int63())
	return strconv.Itoa(rand.Intn(up))
}
