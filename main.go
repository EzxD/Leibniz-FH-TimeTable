package main

import (
	"fmt"

	"github.com/EzxD/Leibniz-FH-TimeTable/cal"
	"github.com/EzxD/Leibniz-FH-TimeTable/xlsx"
	"github.com/manifoldco/promptui"
	"github.com/xuri/excelize/v2"
)

func main() {
	prompt := promptui.Prompt{
		Label:   "Calendar ID",
		Default: "primary",
	}

	calID, err := prompt.Run()
	cal.CalendarId = calID

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	timetableName := promptui.Prompt{
		Label:   "Calendar ID",
		Default: "primary",
	}

	tableName, err := timetableName.Run()
	xlsx.UnmergeAllCells(tableName)

	//fmt.Printf("Your username is %q\n", result)

	promptSelect := promptui.Select{
		Label: "What do you want to do?",
		Items: []string{"Create", "Delete"},
	}

	_, result, err := promptSelect.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %q\n", result)
	if result == "Create" {
		createAll()
	}
	if result == "Delete" {
		cal.DeleteAllUniEvents()
	}

}

func createAll() {
	f, err := excelize.OpenFile("unmerged.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	err, kws := xlsx.GetAllWeeks(f)

	if err != nil {
		fmt.Printf("failed to parse all weeks %v\n", err)
	}

	for k, v := range kws {
		xlsx.AddWeekToCal(k, v, f)
	}
}
