package main

import (
	"flag"
	"fmt"
	"math"
)

func main() {
	typeFlag := flag.String("type", "", "annuity or diff")
	principal := flag.Float64("principal", 0, "Loan principal")
	payment := flag.Float64("payment", 0, "Monthly payment")
	periods := flag.Int("periods", 0, "Number of months")
	interestRate := flag.Float64("interest", 0, "Interest rate")
	flag.Parse()

	if *typeFlag != "annuity" && *typeFlag != "diff" {
		fmt.Println("Incorrect parameters")
		return
	}

	if *interestRate <= 0 {
		fmt.Println("Incorrect parameters")
		return
	}

	if *principal < 0 || *payment < 0 || *periods < 0 {
		fmt.Println("Incorrect parameters")
		return
	}

	if *typeFlag == "diff" && flag.Lookup("payment").Value.String() != flag.Lookup("payment").DefValue {
		fmt.Println("Incorrect parameters")
		return
	}

	count := 0
	for _, name := range []string{"principal", "payment", "periods"} {
		f := flag.Lookup(name)
		if f.Value.String() != f.DefValue {
			count++
		}
	}
	if count < 2 {
		fmt.Println("Incorrect parameters")
		return
	}

	switch *typeFlag {
	case "diff":
		calculateDifferentiated(*principal, *periods, *interestRate)
	case "annuity":
		switch {
		case !isProvided("payment"):
			calculateAnnuityPayment(*principal, *periods, *interestRate)
		case !isProvided("principal"):
			calculateLoanPrincipal(*payment, *periods, *interestRate)
		case !isProvided("periods"):
			calculateNumberOfPayments(*principal, *payment, *interestRate)
		}
	default:
		fmt.Println("Incorrect parameters")
	}

}

func calculateNumberOfPayments(principal, payment, interestRate float64) {
	i := interestRate / (12 * 100)

	n := math.Log(payment/(payment-i*principal)) / math.Log(1+i)
	totalMonths := int(math.Round(n))

	normalized := humanizeMonths(totalMonths)

	fmt.Println(normalized)

	overPayment := payment*float64(totalMonths) - principal
	fmt.Println()
	fmt.Printf("Overpayment = %d\n", int(overPayment))
}

func calculateAnnuityPayment(principal float64, periods int, interestRate float64) {
	n := float64(periods)
	i := interestRate / (12 * 100)

	a := principal * (i * math.Pow(1+i, n)) / (math.Pow(1+i, n) - 1)
	result := int(math.Ceil(a))
	fmt.Printf("Your annuity payment = %d!\n", result)

	overPayment := result*periods - int(principal)
	fmt.Println()
	fmt.Printf("Overpayment = %d\n", overPayment)
}

func calculateLoanPrincipal(payment float64, periods int, interestRate float64) {
	i := interestRate / (12 * 100)
	n := float64(periods)

	p := payment / ((i * math.Pow(1+i, n)) / (math.Pow(1+i, n) - 1))
	result := int(math.Floor(p))

	fmt.Printf("Your loan principal = %d!\n", result)

	overPayment := int(payment*float64(periods)) - result
	fmt.Println()
	fmt.Printf("Overpayment = %d\n", overPayment)
}

func calculateDifferentiated(principal float64, periods int, interestRate float64) {
	i := interestRate / (12 * 100)
	totalPayment := 0
	for month := 1; month <= periods; month++ {
		payment := differentiatedPayment(principal, periods, month, i)
		fmt.Printf("Month %d: payment is %d\n", month, payment)
		totalPayment += payment
	}

	overPayment := totalPayment - int(principal)
	fmt.Println()
	fmt.Printf("Overpayment = %d\n", overPayment)
}

func differentiatedPayment(principal float64, periods int, month int, rate float64) int {
	base := principal / float64(periods)
	interestComponent := rate * (principal - float64(month-1)*base)
	payment := base + interestComponent
	return int(math.Ceil(payment))
}

func pluralizeWord(word string, n int) string {
	if n == 1 {
		return word
	}
	return word + "s"
}

func isProvided(name string) bool {
	f := flag.Lookup(name)
	return f != nil && f.Value.String() != f.DefValue
}

func humanizeMonths(months int) string {
	years := months / 12
	remainingMonths := months % 12

	if years > 0 {
		return fmt.Sprintf("It will take %d %s and %d %s to repay this loan!", years, pluralizeWord("year", years), remainingMonths, pluralizeWord("month", remainingMonths))
	} else {
		return fmt.Sprintf("It will take %d %s to repay this loan!", remainingMonths, pluralizeWord("month", remainingMonths))
	}
}
