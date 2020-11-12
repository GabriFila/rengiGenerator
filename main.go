package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/GabriFila/rengiGenerator/utils"
)

const totWeeks = 1
const stepMinutesResolution = 15
const numIntervals = totWeeks * 7 * 24 * 60 / stepMinutesResolution

func main() {

	data := make([][][]int, 9)
	dataSplitAppl := make(map[string][]int)

	powers := parsePowers("Powers.csv")
	behaviors := parseBehaviors("Behaviors.csv", powers)

	processDataSplitAppl(dataSplitAppl, powers, behaviors, numIntervals)
	writeOutputSingleSplittedByAppliance(dataSplitAppl, "resultSplitSingle.csv", numIntervals)

	// for powerIdx := 0; powerIdx < len(powers); powerIdx++ {
	// 	for behIdx := 0; behIdx < len(behaviors); behIdx++ {
	// 		data[powerIdx][behIdx] = make([]int, numIntervals)
	// 		// processData(data, powers, behaviors, powerIdx, behIdx)
	// 	}
	// }

	titleCols := []string{"Time", "Power"}
	writeOutputGenericPeople(data, "result.csv", titleCols)

}

func getHumanTime(intervalIdx int, intervalLength int) string {
	weeks := (intervalIdx * intervalLength) / (60 * 24 * 7)
	days := ((intervalIdx * intervalLength) % (60 * 24 * 7)) / (60 * 24)
	hours := ((intervalIdx * intervalLength) % (60 * 24 * 7)) % (60 * 24) / (60)
	minutes := ((intervalIdx * intervalLength) % (60 * 24 * 7)) % (60 * 24) % (60)
	return fmt.Sprintf("Week %d DoW %d %02d:%02d", weeks, days, hours, minutes)
}

func getApplBehavior(csvRow []string, applPowers map[string][]int) (appl string, ints []utils.Interval) {
	appl = csvRow[0]
	for i := 2; i < len(csvRow); i++ {
		if len(csvRow[i]) > 0 {
			cur := utils.ExtractDailyIntervals(csvRow[i], 0, i-2)
			ints = append(ints, cur...)
		}
	}
	return appl, ints
}

func processData(data [][][]int, applPowers map[string][]int, behaviors map[string][][]utils.Interval, powerIdx int, behIdx int) {
	for appl := range behaviors {
		for _, v := range behaviors[appl][behIdx] {
			for j := v.Start.GetIndex(15); j <= v.End.GetIndex(15); j++ {
				data[powerIdx][behIdx][j] += applPowers[appl][powerIdx]
			}
		}
	}
}

func processDataSplitAppl(data map[string][]int, applPowers map[string][]int, behaviors map[string][][]utils.Interval, numIntervals int) {
	for appl := range behaviors {
		data[appl] = make([]int, numIntervals)
		for _, v := range behaviors[appl][1] {
			for j := v.Start.GetIndex(15); j <= v.End.GetIndex(15); j++ {
				data[appl][j] += applPowers[appl][0]
			}
		}
	}
}

// print the output, in each column there will be the power consuption of each different person
func writeOutputGenericPeople(data [][][]int, filePath string, columns []string) {
	file, err := os.Create(filePath)
	utils.CheckError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	errFirst := writer.Write(columns)
	utils.CheckError("Cannot write to file", errFirst)
	// personDataLen := len(data[0])
	// for i := 0; i < personDataLen; i++ {
	// 	toWrite := []string{strconv.Itoa(i)}
	// 	for _, personData := range data {
	// 		toWrite = append(toWrite, strconv.Itoa(personData[i]))
	// 	}
	// 	err := writer.Write(toWrite)
	// 	utils.CheckError("Cannot write to file", err)
	// }
}

// print the output, in each column there will be the power consuption of each different person
func writeOutputSingleSplittedByAppliance(data map[string][]int, filePath string, numIntervals int) {
	file, err := os.Create(filePath)
	utils.CheckError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	columns := make([]string, len(data)+2)
	columns[0] = "Interval idx"
	columns[1] = "Human time"
	j := 2
	for appl := range data {
		columns[j] = appl
		j++
	}
	errFirst := writer.Write(columns)
	utils.CheckError("Cannot write to file", errFirst)

	for i := 0; i < numIntervals; i++ {
		toWrite := []string{strconv.Itoa(i), getHumanTime(i, 15)}
		for j := 0; j < len(data); j++ {
			toWrite = append(toWrite, strconv.Itoa(data[columns[j+2]][i]))
		}
		err := writer.Write(toWrite)
		utils.CheckError("Cannot write to file", err)
	}
}

func parseBehaviors(behPath string, powers map[string][]int) (behaviors map[string][][]utils.Interval) {
	behaviors = make(map[string][][]utils.Interval)
	behaviorsFile, behaviorsErr := os.Open(behPath)
	if behaviorsErr != nil {
		log.Fatalln("Couldn't open the behaviors file", behaviorsErr)
	}
	// Start parsing behavior file
	behaviorLine := csv.NewReader(behaviorsFile)
	// read first line to remove col
	_, firstLineErr := behaviorLine.Read()
	if firstLineErr != nil {
		log.Fatal(firstLineErr)
	}
	// Iterate through the records of behaviors file
	for {
		// Read each record from csv
		line, err := behaviorLine.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if err != nil {
			log.Fatal(err)
		}
		appl, newInts := getApplBehavior(line, powers)
		if v, ok := behaviors[appl]; ok {
			behaviors[appl] = append(v, newInts)
		} else {
			behaviors[appl] = [][]utils.Interval{newInts}
		}
	}
	return behaviors
}

func parsePowers(powersFilePath string) (powers map[string][]int) {
	// Open the file
	powersFile, powersErr := os.Open(powersFilePath)
	if powersErr != nil {
		log.Fatalln("Couldn't open the powers file", powersErr)
	}

	powers = make(map[string][]int)
	// Start parsing power file
	powersLine := csv.NewReader(powersFile)
	// read first line to remove col
	var _, firstLineErr = powersLine.Read()
	if firstLineErr != nil {
		log.Fatal(firstLineErr)
	}
	// Iterate through the records
	for {
		// Read each record from csv
		line, err := powersLine.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		appl := line[0]
		newPower, err := strconv.Atoi(line[1])
		if err != nil {
			log.Fatal(err)
		}
		if v, ok := powers[appl]; ok {
			powers[appl] = append(v, newPower)
		} else {
			powers[appl] = []int{newPower}
		}
	}
	return powers
}
