package xlsx

import (
	"fmt"
	"strings"
	"time"

	"github.com/EzxD/Leibniz-FH-TimeTable/cal"
	"github.com/xuri/excelize/v2"
)

func AddWeekToCal(yAxis uint, xAxis uint, f *excelize.File) {
	startX := xAxis + 1
	startY := yAxis + 2
	dateRow := yAxis + 1
	//xValue, time
	dates := make(map[uint]string)
	//yValue, time
	times := make(map[uint]string)

	rows, _ := f.GetRows("dIT21")
	cols, _ := f.GetCols("dIT21")
	for i, v := range rows[dateRow] {
		//fmt.Println(i, v)
		dates[uint(i)] = v
	}
	for i, v := range cols[xAxis] {
		if i >= int(yAxis+8) || i <= int(yAxis+1) {
			continue
		}
		times[uint(i)] = v
	}

	for yAxisTime := 0; yAxisTime <= 5; yAxisTime++ {
		for i, v := range rows[startY+uint(yAxisTime)] {
			if i < int(startX) {
				continue
			}
			if v != "" {
				fmt.Printf("%s findet am %s um %s Uhr statt\n", v, dates[uint(i)], times[uint(startY+uint(yAxisTime))])
				times := strings.Split(times[uint(startY+uint(yAxisTime))], "-")
				start, end := ParseTimeString(dates[uint(i)], times[0], dates[uint(i)], times[1])
				cal.CreateEvent(v, start, end)
			}
		}
	}
}

func ParseTimeString(startDate string, startTime, endDate string, endTime string) (string, string) {
	layout := "02/01/2006 15:04"
	ger, _ := time.LoadLocation("Europe/Berlin")
	time.Local = ger
	start, _ := time.ParseInLocation(layout, startDate+"2022 "+startTime, ger)
	end, _ := time.ParseInLocation(layout, endDate+"2022 "+endTime, ger)

	startString := start.In(ger).Format("2006-01-02T15:04:00")
	endString := end.In(ger).Format("2006-01-02T15:04:00")

	return startString + "+02:00", endString + "+02:00"
}

func GetAllWeeks(f *excelize.File) (error, map[uint]uint) {
	m := make(map[uint]uint)
	rows, err := f.GetRows("dIT21")
	for i, v := range rows {
		for i2, v2 := range v {
			if v2 == "KW" {
				m[uint(i)] = uint(i2)
				//fmt.Println(i, i2)
			}
		}
	}
	return err, m
}

func UnmergeAllCells(name string) error {
	f, err := excelize.OpenFile(name)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	mcells, err := f.GetMergeCells("dIT21")
	for _, mc := range mcells {
		fmt.Println(mc)
		f.UnmergeCell("dIT21", mc.GetStartAxis(), mc.GetEndAxis())
	}
	f.SaveAs("unmerged.xlsx")
	return err
}
