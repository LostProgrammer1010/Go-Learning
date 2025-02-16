package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	start()
}

func start() {
	var promptAnswer string
	var timerFlag, fileFlag = setflags()
	questions := generateQuestions(*fileFlag)
	seconds := time.Duration(*timerFlag)

	for stay := true; stay; {
		clearTerm()
		printStartMenu()
		fmt.Scan(&promptAnswer)
		switch promptAnswer {
		case "1":
			startQuiz(questions, seconds)
			stay = false
		default:
			stay = true
		}
	}
}

type Question struct {
	question string
	answer   string
}

func setflags() (timerFlag *int, fileFlag *string) {
	timerFlag = flag.Int("t", 30, "Set Timer for Quiz (Seconds)")
	fileFlag = flag.String("f", "problem.csv", "Name of CSV file with Questions")
	flag.Parse()
	return
}

func generateQuestions(filename string) (questions []Question) {

	file, err := os.Open(filename)

	if err != nil {
		log.Fatal("Unable to read the csv file provided", err)
	}

	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()

	if err != nil {
		fmt.Println("Error Reading Records")
	}

	for _, record := range records {
		questions = append(questions, Question{question: record[0], answer: record[1]})
	}
	return
}

func printStartMenu() {
	prompt :=
		`
Welcome to the CLI Quiz.
1 - Start Quiz

Enter Number for Option: `
	fmt.Print(prompt)
}

func clearTerm() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func startQuiz(questions []Question, quizLength time.Duration) {

	var correctCount int

	time.Duration(10).Seconds()
	timer := time.NewTimer(quizLength * time.Second)
	done := make(chan bool)

	go func() {
		for i, question := range questions {
			clearTerm()
			var answer string
			fmt.Printf("==== Question %d ====\n", i+1)
			fmt.Println(question.question)
			fmt.Print("Type in your answer: ")
			fmt.Scan(&answer)

			if answer == question.answer {
				correctCount++
			}
		}
		done <- true
	}()

	select {
	case <-timer.C:
		clearTerm()
		fmt.Println("Time is Up!!")
	case <-done:

	}

	fmt.Printf("Your score was: %d/%d", correctCount, len(questions))

}
