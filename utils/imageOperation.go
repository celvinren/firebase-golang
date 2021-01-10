package utils

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"os"
	"strings"
)

func readImageFileToBuffer() []byte {
	imgFile, err := os.Open("images.jpg")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer imgFile.Close()

	// create a new buffer base on file size
	fInfo, _ := imgFile.Stat()
	var size = fInfo.Size()
	buf := make([]byte, size)

	// read file content into buffer
	fReader := bufio.NewReader(imgFile)
	fReader.Read(buf)
	return buf
}

func ConvertBufferToBase64(buf []byte) string {
	imgBase64Str := base64.StdEncoding.EncodeToString(buf)
	log.Println(imgBase64Str)
	return imgBase64Str
}

func ConvertBase64StrToIOReader(base64Str string) io.Reader {
	return base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64Str))
}

func SaveImageFromIOReader(isReader io.Reader, fileName string) {
	m, formatString, err := image.Decode(isReader)
	if err != nil {
		log.Println(err)
	}
	bounds := m.Bounds()
	fmt.Println("base64toJpg", bounds, formatString)

	//Encode from image format to writer
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Println(err)
		return
	}

	err = jpeg.Encode(f, m, &jpeg.Options{Quality: 75})
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("Jpg file", fileName, "created")
}
