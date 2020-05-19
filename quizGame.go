package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {

	csvFileName := flag.String("csv", "problems.csv", "This parameter can be used to pass csv")
	durationOfQuiz := flag.Int("time", 10, "Maximum duration of the test")
	csvFile, err := os.Open(*csvFileName)
	if err != nil {
		fmt.Printf("Unable to open file %s", *csvFileName)
		os.Exit(1)
	}

	fmt.Printf("Opened File %s successfully \n", *csvFileName)

	csvReader := csv.NewReader(csvFile)
	problemRecords, _ := csvReader.ReadAll()
	problems := parseProblems(problemRecords)

	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(problems), func(i, j int) {
		problems[i], problems[j] = problems[j], problems[i]
	})

	gameTimer := time.NewTimer(time.Duration(*durationOfQuiz) * time.Second)

	var score int
	var answer string
	answerChannel := make(chan string)
	for i, problem := range problems {

		fmt.Printf("\nQuestion %d : %s =", i, problem.question)

		go getAnswer(answer, answerChannel)

		select {

		case <-gameTimer.C:
			fmt.Printf("\nYou Scored %d out of %d", score, len(problems))
			return
		case answer = <-answerChannel:
			if strings.EqualFold(answer, problem.answer) {
				score++
			}

		}

	}
	fmt.Printf("\nYou Scored %d out of %d", score, len(problems))
}

func getAnswer(answer string, answerChannel chan string) {

	fmt.Scanf("%s", &answer)
	answerChannel <- answer

}

func parseProblems(inputRecords [][]string) []problem {

	problems := make([]problem, len(inputRecords))

	for i, record := range inputRecords {

		problems[i] = problem{
			question: record[0],
			answer:   strings.TrimSpace(record[1]),
		}

	}

	return problems
}

type problem struct {
	question string
	answer   string
}
