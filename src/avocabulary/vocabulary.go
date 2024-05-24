package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)


var vocabulary = "english_russian.txt"

func init() {
	dir, _ := filepath.Split(os.Args[0])
	vocabulary = filepath.Join(dir, vocabulary)
}

func getFileNames() (inFileName, outFileName string, err error) {
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "-help") {
		err := fmt.Errorf("Usage is: %s <input file> <output file>", filepath.Base(os.Args[0]))
		return "", "", err
	}
	if len(os.Args) > 1 {
		inFileName = os.Args[1]
		if len(os.Args) > 2 {
			outFileName = os.Args[2]
		}
	}
	if inFileName != "" && inFileName == outFileName {
		log.Fatal("Output file name must differ from input one!")
	}
	return inFileName, outFileName, nil
}

func getValidator(vocabulary string) (func(string, string) bool, error) {
	arrayOfBytes, err := ioutil.ReadFile(vocabulary)
	if err != nil {
		return nil, err
	}
	wholeVocabulary := string(arrayOfBytes)
	foreignWord := make(map[string]string)
	
	rows := strings.Split(wholeVocabulary, "\n")
	for _, row := range rows {
		words := strings.Split(row, ":=")
		if len(words) == 2 {
			words[1] = strings.TrimRight(words[1], "\r\n")
			foreignWord[words[0]] = words[1]
		}
	}
	return func(firstWord, secondWord string) bool {
		secondWord = strings.TrimRight(secondWord, "\r\n")
		
		if patternWord, found := foreignWord[firstWord]; found {
			if strings.Compare(patternWord, secondWord) == 0 {
				return true 
			}
		} else if patternWord, found := foreignWord[secondWord]; found {
			if strings.Compare(patternWord, firstWord) == 0 {
				return true 
			}
		}
		return false
	}, nil
}

func compare(line string, validator func(string, string) bool) (string){
	words := strings.Split(line, ":=")
	
	if len(words) == 2 {
		if isCorrect := validator(words[0], words[1]); isCorrect {
			return fmt.Sprintf("%s %s\r\n", line, "CORRECT")
		}
	}
	return fmt.Sprintf("%s\r\n", line)
}

func check (ioReader io.Reader, ioWriter io.Writer) (err error) {
	bufioReader := bufio.NewReader(ioReader)
	bufioWriter := bufio.NewWriter(ioWriter)
	defer func() {
		if err == nil {
			err = bufioWriter.Flush() 
		}
	}()
	var validator func(string, string) bool
	if validator, err = getValidator(vocabulary); err != nil {
		return err
	}
	endOfFile := false
	for !endOfFile {
		var row string
		row, err = bufioReader.ReadString('\n')
		row = strings.TrimRight(row, "\r\n")
		
		if err == io.EOF {
			err = nil
			endOfFile = true
		} else if err != nil {
			return err
		}
		row = compare(row, validator)
		if _, err = bufioWriter.WriteString(row); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	inFileName, outFileName, err := getFileNames()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ioReader := os.Stdin
	ioWriter := os.Stdout
	if inFileName != "" {
		if ioReader, err = os.Open(inFileName); err != nil {
			log.Fatal(err)
		}
		defer ioReader.Close()
	}
	if outFileName != "" {
		if ioWriter, err = os.Create(outFileName); err != nil {
			log.Fatal(err)
		}
		defer ioWriter.Close()
	}
	
	if err = check(ioReader, ioWriter); err != nil {
		log.Fatal(err) 
	} 
}
