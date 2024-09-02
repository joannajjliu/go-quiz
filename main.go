package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func shuffleQuestions(records [][]string) {
	for i := range records {
		j := rand.Intn(i + 1)
		records[i], records[j] = records[j], records[i]
	}
}

func sendQuestions(done chan bool, inputReader *bufio.Reader, filename string, isRandom bool) {
	file, err := os.Open(filename + ".csv")

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

	records, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println("csv read err:", err)
		done <- true
		return
	}
	if isRandom {
		shuffleQuestions(records)
	}

	for _, record := range records {
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

	fmt.Printf("Of %d questions, you got %d correct!\n", totalQuestions, totalCorrect)
	done <- true
}

func main() {
	filename := flag.String("filename", "problems", "csv filename of the problems.")
	timeLimit := flag.Int("timeLimit", 30, "time limit to complete quiz, in seconds.")
	isRandom := flag.Bool("isRandom", false, "determines if questions should be given in random order.")

	flag.Parse()
	fmt.Println("flags:", *filename, *timeLimit)

	inputReader := bufio.NewReader(os.Stdin)

	fmt.Println("Press enter when you're ready to start the quiz!")
	inputReader.ReadString('\n')

	done := make(chan bool)
	go sendQuestions(done, inputReader, *filename, *isRandom)

	select {
	case <-done:
		fmt.Println("You've solved all the questions!")
		return
	case <-time.After(time.Duration(*timeLimit) * time.Second):
		fmt.Println("You took too long. Better luck next time!")
		return
	}
}

func cleanStr(in string) string {
	return strings.ToLower(strings.TrimSpace(in))
}
