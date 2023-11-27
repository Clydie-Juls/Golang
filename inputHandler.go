package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type inputCommands struct {
	function interface{}
}

func inputDailySalary(currSalary *float32) {
	for {
		fmt.Printf("Input daily salary(currently %.2f): ", *currSalary)
		_, err := fmt.Scan(currSalary)
		if err != nil {
			fmt.Println("Invalid input. Try again.")
			continue
		}

		if *currSalary < 0 {
			fmt.Println("Input out of range. Try again.")
			continue
		}

		break
	}
}

func inputMaxRegularWorkHours(currHours *int) {
	for {
		fmt.Printf("Input maximum regular hours of work (currently %d): ", *currHours)
		_, err := fmt.Scan(currHours)
		if err != nil {
			fmt.Println("Invalid input. Try again.")
			continue
		}

		if *currHours < 8 {
			fmt.Println("Number of regular hours of work should be at least 8 hours. Try again.")
			continue
		}

		if *currHours > 24 {
			fmt.Println("You cannot work for more than a day, you need rest. Try again.")
			continue
		}

		break
	}
}

func inputNoOfWorkDays(currDays *int, dayTypes []string) {
	for {
		fmt.Printf("Input number of work days (currently %d):", *currDays)
		_, err := fmt.Scan(currDays)
		if err != nil {
			fmt.Println("Invalid input. Try again.")
			continue
		}

		if *currDays < 0 {
			fmt.Println("No negative values allowed. Try again.")
			continue
		}

		if *currDays > 7 {
			fmt.Println("There's only 7 days in a week. Try again.")
			continue
		}

		for i := range dayTypes {
			if i+1 <= *currDays {
				dayTypes[i] = "Normal Day"
			} else {
				dayTypes[i] = "Rest Day"
			}
		}

		break
	}
}

func inputInTime(currInTime *int) {
	for {
		fmt.Printf("Input In time(currently %04d). Input a non-number to skip: ", *currInTime)
		_, err := fmt.Scan(currInTime)
		if err != nil {
			break
		}

		if !isValidMilitaryTime(*currInTime) {
			fmt.Println("Not a valid military time input. Try again.")
			continue
		}

		break
	}
}

func inputOutTime(currOutTime *int) {
	for {
		fmt.Printf("Input Out time(currently %04d). Input a non-number to skip: ", *currOutTime)
		_, err := fmt.Scan(currOutTime)
		if err != nil {
			break
		}

		if !isValidMilitaryTime(*currOutTime) {
			fmt.Println("Not a valid military time input. Try again.")
			continue
		}

		break
	}
}

func isValidMilitaryTime(time int) bool {
	if time < 0 || time > 2359 {
		return false
	}

	if time%100 > 59 {
		return false
	}

	return true
}

func inputDayType(dayTypes []string) {
	for {
		for i := range dayTypes {
			fmt.Printf("Day %d: %s\n", i+1, dayTypes[i])
		}

		day := 0
		fmt.Printf("Input Which day[1-7]: ")
		_, err := fmt.Scan(&day)
		if err != nil || day < 1 || day > 7 {
			fmt.Println("Invalid input. Try again.")
			continue
		} else {
			types := []string{
				"Normal Day",
				"SNWH",
				"RH",
			}

			restTypes := []string{
				"Rest Day",
				"SNWH, Rest Day",
				"RH, Rest Day",
			}

			dayStr := ""
			fmt.Println("Types:")
			for i := range types {
				if i == 0 {
					fmt.Printf("\t%s(Can be rest day depending on the number of work days)\n", types[i])
				} else {
					fmt.Printf("\t%s\n", types[i])
				}
			}
			fmt.Printf("What day type: (currently %s):", dayTypes[day-1])
			bio := bufio.NewReader(os.Stdin)
			line, _, err1 := bio.ReadLine()
			line, _, err1 = bio.ReadLine()
			if err1 != nil {
				return
			}
			dayStr = string(line)
			validInput := false
			for i := 0; i < len(types); i++ {
				if types[i] == dayStr {
					validInput = true
					if strings.Contains(dayTypes[day-1], "Rest Day") {
						dayTypes[day-1] = restTypes[i]
					} else {
						dayTypes[day-1] = dayStr
					}
					break
				}
			}

			if validInput {
				break
			}

			fmt.Println("Input unknown. Try again.")
		}
		break
	}
}
