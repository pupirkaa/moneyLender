package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func readFile() []string {
	var transactionData []string
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file: %v", err)
		os.Exit(1)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		transactionData = append(transactionData, scanner.Text())
	}
	return transactionData
}

func printResult(m map[string]int) {
	mk := make([]string, len(m))
	i := 0
	for k, _ := range m {
		mk[i] = k
		i++
	}
	sort.Strings(mk)

	for i := 0; i < len(mk); i++ {
		fmt.Println(mk[i], " ", m[mk[i]])
	}
}

func main() {
	transactionData := readFile()

	debts := make(map[string]int)
	for i := 0; i < len(transactionData); i++ {
		splitedString := (strings.Split(transactionData[i], " "))
		if len(splitedString) != 5 {
			fmt.Fprintln(os.Stderr, "your data is incorrect")
			os.Exit(1)
		}
		var (
			lender = splitedString[0]
			lendee = splitedString[4]
			money  = splitedString[2]
		)

		money = strings.TrimSuffix(money, "$")
		amountOfMoney, err := strconv.Atoi(money)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse money amount: %v", err)
			os.Exit(1)
		}
		debts[lender] += amountOfMoney
		debts[lendee] -= amountOfMoney
	}

	//fmt.Println(debts)
	printResult(debts)

}
