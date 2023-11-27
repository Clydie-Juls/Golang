//********************
//Last name: Marindo
//Language: GO
//Paradigm(s): imperative, procedural, object-oriented, functional
//********************

package main

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/olekukonko/tablewriter"
)

type computation struct {
	equation  string
	timeRange string
	result    float32
}

type config struct {
	salary        float32
	inTime        int
	outTime       int
	noOfWorkHours int
	dayType       string
}

func styleTable(t *table.Writer) {
	(*t).Style().Title.Format = tablewriter.ALIGN_CENTER
	(*t).SetStyle(table.StyleRounded)

	(*t).Style().Options.SeparateColumns = true
	(*t).Style().Options.SeparateRows = true
}

func computeDailySalary(salary float32, inTime int, outTime int, workHours int, dayType string) computation {
	multiplier := map[string]float32{
		"Normal Day":     1.00,
		"Rest Day":       1.30,
		"SNWH":           1.30,
		"SNWH, Rest Day": 1.50,
		"RH":             2.00,
		"RH, Rest Day":   2.60,
	}

	diff := (min(outTime-inTime-1, workHours*100)) / 100
	equations := map[string]string{
		"Normal Day":     "Daily Rate",
		"Rest Day":       fmt.Sprintf("Daily Rate x Rest Day\n\t\t= %.2f x %d / %d x 1.30", salary, diff, workHours),
		"SNWH":           fmt.Sprintf("Daily Rate x SNWH\n\t\t= %.2f x %d / %d  x 1.30", salary, diff, workHours),
		"SNWH, Rest Day": fmt.Sprintf("Daily Rate x SNWH-Rest Day\n\t\t= %.2f x %d / %d  x 1.50", salary, diff, workHours),
		"RH":             fmt.Sprintf("Daily Rate x RH\n\t\t= %.2f x %d / %d  x 2.00", salary, diff, workHours),
		"RH, Rest Day":   fmt.Sprintf("Daily Rate x RH-Rest Day\n\t\t= %.2f x %d / %d  x 2.60", salary, diff, workHours),
	}

	newSalary := float32(diff) * salary / float32(workHours)
	result := salary
	if outTime == inTime {
		if dayType != "Rest Day" && dayType != "SNWH, Rest Day" && dayType != "RH, Rest Day" {
			return computation{"", "", 0}
		}
	} else {
		result = newSalary * multiplier[dayType]
	}

	toTime := min(max(0, outTime-inTime), inTime+workHours*100) % 2400

	return computation{equations[dayType],
		fmt.Sprintf("%04d-%04d", inTime, toTime), result}
}

func computeDayOTSalary(salary float32, currTime int, overtime int, workHours int, dayType string) computation {

	multiplier := map[string]float32{
		"Normal Day":     1.25,
		"Rest Day":       1.69,
		"SNWH":           1.69,
		"SNWH, Rest Day": 1.95,
		"RH":             2.60,
		"RH, Rest Day":   3.38,
	}

	equation := fmt.Sprintf("Hours OT x Hourly Rate\n\t\t= %d x %d ÷ %d x %.2f",
		int(salary), overtime, workHours, multiplier[dayType])
	return computation{equation,
		fmt.Sprintf("%04d-%04d", currTime%2400, (currTime+overtime*100)%2400),
		float32(overtime) * salary / float32(workHours) * multiplier[dayType]}
}

func computeNightShiftSalary(salary float32, nsHours int, workHours int, fromTime int, toTime int) computation {
	equation := fmt.Sprintf("Hours on NS x Hourly Rate x NSD\n\t\t= %d x %.2f ÷ %d x 1.10",
		nsHours, salary, workHours)
	return computation{equation,
		fmt.Sprintf("%04d-%04d", fromTime%2400, toTime%2400),
		float32(nsHours) * salary / float32(workHours) * 1.1}

}

func computeNightOTSalary(salary float32, currTime int, overtime int, workHours int, dayType string) computation {

	multiplier := map[string]float32{
		"Normal Day":     1.375,
		"Rest Day":       1.859,
		"SNWH":           1.859,
		"SNWH, Rest Day": 2.145,
		"RH":             2.86,
		"RH, Rest Day":   371.8,
	}

	equation := fmt.Sprintf("Hours OT x NS-OT Hourly Rate\n\t\t= %d x %d ÷ %d x %.3f",
		int(salary), overtime, workHours, multiplier[dayType])
	return computation{equation,
		fmt.Sprintf("%04d-%04d", currTime, (currTime+overtime*100)%2400),
		float32(overtime) * salary / float32(workHours) * multiplier[dayType]}
}

func calculateSalary(initSalary float32, inTime int, outTime int, workHours int, dateType string) ([]computation, int, int, int) {
	var computations []computation

	if outTime < inTime {
		outTime += 2400
	}

	dailyComp := computeDailySalary(initSalary, inTime, outTime, workHours, dateType)
	computations = append(computations, dailyComp)
	nightShiftHours := 0
	if inTime != outTime {
		tempInTime := inTime
		tempOutTime := outTime
		if inTime >= 0 && inTime <= 600 {
			tempInTime += 2400
			tempOutTime += 2400
		}

		fromTime := max(2200, tempInTime)
		g := min(tempOutTime, tempInTime+workHours*100)
		toTime := min(600+2400, g)

		nightShiftHours = max(0, (toTime-fromTime)/100)
		nsComp := computeNightShiftSalary(initSalary, nightShiftHours, workHours, fromTime, toTime)
		if nsComp.result > 0 {
			computations = append(computations, nsComp)
		}
	}

	overtimeHours := max(0, outTime-inTime-workHours*100-100)
	currentTime := inTime + workHours*100 + 100
	dot := 0
	not := 0
	for overtimeHours > 0 {
		var hoursOfWork int
		if currentTime >= 2200 && currentTime < 600+2400 {
			hoursOfWork = min(600+2400, outTime) - currentTime
			dot += hoursOfWork / 100
			computations = append(computations, computeNightOTSalary(initSalary, currentTime, hoursOfWork/100, workHours, dateType))
			currentTime = min(600+2400, outTime)
		} else {
			tempTime := currentTime
			if currentTime < 2200 {
				hoursOfWork = min(2200, outTime) - currentTime
				currentTime = min(2200, outTime)
			} else {
				hoursOfWork = min(2200+2400, outTime) - currentTime
				currentTime = min(2200+2400, outTime)
			}
			not += hoursOfWork / 100
			computations = append(computations, computeDayOTSalary(initSalary, tempTime, hoursOfWork/100, workHours, dateType))
		}

		overtimeHours -= hoursOfWork
	}

	return computations, dot, not, nightShiftHours
}

func renderData(initSalary float32, inTime int, outTime int, workHours int, dateType string) float32 {
	overtimeHours := max(0, outTime-inTime-(workHours*100)-100)
	// if out time reaches the next day e.g. 0000
	if outTime < inTime {
		overtimeHours += 2400
	}
	overtimeHours /= 100

	computations, dot, not, nsh := calculateSalary(initSalary, inTime, outTime, workHours, dateType)
	otDisplay := ""
	otVal := ""
	if nsh > 0 {
		otDisplay = "Hours of night shift"
		otVal = fmt.Sprintf("%d", not)

	} else {
		otDisplay = "Hourly Overtime (Night Shift Overtime)"
		otVal = fmt.Sprintf("%d (%d)", dot, not)
	}

	var sum float32 = 0
	for _, comp := range computations {
		sum += comp.result
	}

	t := table.NewWriter()
	t.SetTitle("Sample")

	data := []table.Row{
		{"Daily Rate", fmt.Sprintf("%.2f", initSalary), fmt.Sprintf("%.2f", initSalary)},
		{"IN Time", fmt.Sprintf("%04d", inTime), fmt.Sprintf("%04d", inTime)},
		{"OUT Time", fmt.Sprintf("%04d", outTime), fmt.Sprintf("%04d", outTime)},
		{"Day type", dateType, dateType},
		{otDisplay, otVal, otVal},
		{"Salary for the day", fmt.Sprintf("%.2f", sum), fmt.Sprintf("%.2f", sum)},
	}

	t.AppendRows(data, table.RowConfig{AutoMerge: true})

	t.AppendRow(table.Row{"Computations: "})
	for _, comp := range computations {
		t.AppendRow(table.Row{comp.equation, comp.timeRange, fmt.Sprintf("%.2f", comp.result)})
	}

	styleTable(&t)

	fmt.Println(t.Render())
	fmt.Println()

	return sum
}

func promptConfigCommands() {
	rows := []table.Row{
		{1, "Input Daily Salary"},
		{2, "Input Max Regular Work Hours"},
		{3, "Input Number of Work Days"},
		{4, "Input Default In Time"},
		{5, "Input Default Out Time"},
		{6, "Input Day type"},
		{7, "Exit"},
	}

	t := table.NewWriter()
	t.SetTitle("Commands")
	t.AppendRows(rows)
	styleTable(&t)

	fmt.Println(t.Render())
	fmt.Println("Enter command:")
}

func promptMenuCommands() {
	rows := []table.Row{
		{1, "Edit configuration"},
		{2, "Generate Payroll"},
		{3, "Exit"},
	}

	t := table.NewWriter()
	t.SetTitle("Commands")
	t.AppendRows(rows)
	styleTable(&t)

	fmt.Println(t.Render())
	fmt.Println("Enter command:")
}

func menuCommands() {
	var currSalary float32 = 500.00
	currHours := 8
	noOfWorkDays := 5
	defaultInTime := 900
	defaultOutTime := 900
	dayTypes := []string{"Normal Day", "Normal Day", "Normal Day", "Normal Day", "Normal Day", "Rest Day", "Rest Day"}

	commands := map[int]interface{}{
		1: func() {
			configCommands(&currSalary, &currHours, &noOfWorkDays, &defaultInTime, &defaultOutTime, dayTypes)
		},
		2: func() { generatePayroll(currSalary, currHours, noOfWorkDays, defaultInTime, defaultOutTime, dayTypes) },
	}
	for {
		command := 0
		promptMenuCommands()
		_, err := fmt.Scan(&command)
		if err != nil {
			fmt.Println("Invalid input. Try again.")
			continue
		}

		val, exists := commands[command]

		if command == len(commands)+1 {
			break
		} else if exists {
			if fn, ok := val.(func()); ok {
				fn()
			}
		} else {
			// The key does not exist in the map.
			fmt.Println("Invalid choice. Please select a valid option.")
		}
	}
}

func generateWeeklyTotal(total float32) {
	t := table.NewWriter()
	t.SetTitle("Weekly Total Payroll")
	t.AppendRow(table.Row{"         Total         ", fmt.Sprintf("%.2f", total)})
	styleTable(&t)

	println(t.Render())
}

func generatePayroll(currSalary float32, currHours int, noOfWorkDays int, defaultInTime int, defaultOutTime int, dayTypes []string) {
	week := []config{
		{currSalary, defaultInTime, defaultOutTime, currHours, dayTypes[0]},
		{currSalary, defaultInTime, defaultOutTime, currHours, dayTypes[1]},
		{currSalary, defaultInTime, defaultOutTime, currHours, dayTypes[2]},
		{currSalary, defaultInTime, defaultOutTime, currHours, dayTypes[3]},
		{currSalary, defaultInTime, defaultOutTime, currHours, dayTypes[4]},
		{currSalary, defaultInTime, defaultOutTime, currHours, dayTypes[5]},
		{currSalary, defaultInTime, defaultOutTime, currHours, dayTypes[6]},
	}

	for i := 0; i < 7; i++ {
		fmt.Printf("On Day %d:\n", i+1)
		inputOutTime(&(week[i].outTime))
	}

	total := float32(0)
	for i := 0; i < 7; i++ {
		total += renderData(week[i].salary, week[i].inTime, week[i].outTime, week[i].noOfWorkHours, week[i].dayType)
	}

	generateWeeklyTotal(total)
}

func configCommands(currSalary *float32, currHours *int, noOfWorkDays *int, defaultInTime *int,
	defaultOutTime *int, dayTypes []string) {
	commands := map[int]interface{}{
		1: func() { inputDailySalary(currSalary) },
		2: func() { inputMaxRegularWorkHours(currHours) },
		3: func() { inputNoOfWorkDays(noOfWorkDays, dayTypes) },
		4: func() { inputInTime(defaultInTime) },
		5: func() { inputOutTime(defaultOutTime) },
		6: func() { inputDayType(dayTypes) },
	}

	for {
		command := 0
		promptConfigCommands()
		_, err := fmt.Scan(&command)
		if err != nil {
			fmt.Println("Invalid input. Try again.")
			continue
		}

		val, exists := commands[command]

		if command == len(commands)+1 {
			break
		} else if exists {
			if fn, ok := val.(func()); ok {
				fn()
			}
		} else {
			// The key does not exist in the map.
			fmt.Println("Invalid choice. Please select a valid option.")
		}
	}

}

func main() {
	println("Welcome to the weekly salary computer")
	menuCommands()
}
