package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func main() {
	// Command-line flag for setting the time limit
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")

	flag.Parse()
	// Open the CSV file
	file, err := os.Open("problems.csv")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()
	//defer file.Close()

	// Create a CSV reader
	csvReader := csv.NewReader(file)

	// Create a reader for user input
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Press 'Enter' when you are ready to start the quiz.")
	_, _ = reader.ReadString('\n')

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	// Variable to keep track of correct answers
	var correctCount int
	var totalCount int

	answerCh := make(chan string)

	// Read and process the CSV data
	for {
		// Read one record (a slice of strings)
		record, err := csvReader.Read()

		// Check for end of file
		if err == io.EOF {
			break
		}

		// Handle other errors
		if err != nil {
			panic(err)
		}

		// Ask the question
		totalCount++
		fmt.Printf("Question: %s\n", record[0])

		go func() {
			userAnswer, _ := reader.ReadString('\n')
			answerCh <- strings.TrimSpace(userAnswer)
		}()

		// Wait for either the timer to expire or an answer to be submitted
		select {
		case <-timer.C:
			fmt.Println("\nTime's up!")
			fmt.Printf("You scored %d out of %d.\n", correctCount, totalCount)
			return
		case userAnswer := <-answerCh:
			if userAnswer == record[1] {
				correctCount++
			}
		}

		//// Get user's answer
		//fmt.Print("Enter your answer: ")
		//userAnswer, _ := reader.ReadString('\n')
		//userAnswer = strings.TrimSpace(userAnswer)
		//
		//// Check if the answer is correct
		//if userAnswer == record[1] {
		//	correctCount++
		//}
	}
	fmt.Printf("You got %d out of %d correct.\n", correctCount, totalCount)

}
