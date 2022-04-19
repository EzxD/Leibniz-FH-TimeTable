package main

import (
	"fmt"

	"github.com/ezxd/leibniz-fh-timetable/cal"
	"github.com/ezxd/leibniz-fh-timetable/xlsx"
	"github.com/manifoldco/promptui"
	"github.com/xuri/excelize/v2"
)

func main() {
	prompt := promptui.Prompt{
		Label:   "Calendar ID",
		Default: "primary",
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("Your username is %q\n", result)

	promptSelect := promptui.Select{
		Label: "What do you want to do?",
		Items: []string{"Create", "Delete"},
	}

	_, result, err = promptSelect.Run()

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
	f, err := excelize.OpenFile("test.xlsx")
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

	for k, v := range kws {
		xlsx.AddWeekToCal(k, v, f)
	}
}
