package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func sendQuestions(done chan bool, inputReader *bufio.Reader) {
	file, err := os.Open("problems.csv")

	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("file read!")

	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.FieldsPerRecord = -1
	totalQuestions := 0
	totalCorrect := 0

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			fmt.Printf("Of %d questions, you got %d correct!\n", totalQuestions, totalCorrect)
			fmt.Println("end of file reached:", err)
			done <- true
			return
		}
		if err != nil {
			fmt.Println("csv read err:", err)
			return
		}

		totalQuestions += 1

		lastIdx := len(record) - 1
		q := record[0:lastIdx]
		ans := record[lastIdx]

		fmt.Println(strings.Join(q, ","))

		in, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Println("user input err:", err)
			return
		}

		if cleanStr((in)) == cleanStr((ans)) {
			totalCorrect += 1
		}
	}
}

func main() {
	inputReader := bufio.NewReader(os.Stdin)

	fmt.Println("Press enter when you're ready to start the quiz!")
	inputReader.ReadString('\n')

	done := make(chan bool)
	go sendQuestions(done, inputReader)

	select {
	case <-done:
		fmt.Println("You've solved all the questions!")
		return
	case <-time.After(2 * time.Second):
		fmt.Println("You took too long. Better luck next time!")
		return
	}
}

func cleanStr(in string) string {
	return strings.ToLower(strings.TrimSpace(in))
}
