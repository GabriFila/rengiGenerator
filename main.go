package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// assuming 15 minutes intervals

func timeToHM(in string) (hours, minutes int) {
	hm := strings.Split(in, ":")
	hours, errH := strconv.Atoi(hm[0])
	if errH != nil {
		log.Fatal("could not convert hours")
	}

	minutes, errM := strconv.Atoi(hm[1])
	if errM != nil {
		log.Fatal("could not convert minutes")
	}
	return hours, minutes
}

func extractIntervals(in string, dayOffset int) (start, end int) {
	se := strings.Split(in, "-")
	start := timeToHM(dayOffset)
}

func analyzeRow(csvRow []string) [][]int {
	ret := make([][]int, 7)
	for i := 2; i < len(csvRow); i++ {
		inDayIntervals := strings.Split(csvRow[i], "e")
		for i, v := range inDayIntervals {
			// i is useful
			// transform hour to interval index

			hours, minutes := timeToHM(v)

		}
	}
	return ret
}

func main() {
	// Open the file
	powersFile, powersErr := os.Open("Powers.csv")
	if powersErr != nil {
		log.Fatalln("Couldn't open the powers file", powersErr)
	}

	behaviorsFile, behaviorsErr := os.Open("Behaviors.csv")
	if powersErr != nil {
		log.Fatalln("Couldn't open the behaviors file", behaviorsErr)
	}

	powers := make(map[string][]int)
	behaviors := make(map[string][]int)

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
		if v, ok := powers[appl]; ok {
			powers[appl] = append(v, newPower)
		} else {
			powers[appl] = []int{newPower}
		}
		if err != nil {
			log.Fatal(err)
		}
	}
	// for k, v := range powers {
	// 	fmt.Printf("Appliance: %s -> Power: %d\n", k, v)
	// }

	// Start parsing power file
	behaviorLine := csv.NewReader(behaviorsFile)
	// read first line to remove col
	_, firstLineErr = behaviorLine.Read()
	if firstLineErr != nil {
		log.Fatal(firstLineErr)
	}
	// Iterate through the records
	for {
		// Read each record from csv
		line, err := behaviorLine.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		appl := line[0]
		newPower, err := strconv.Atoi(line[1])
		if v, ok := behaviors[appl]; ok {
			behaviors[appl] = append(v, newPower)
		} else {

			if _, ok := powers[appl]; !ok {
				log.Fatal("No matching power for behavior of ", appl)
			}
			behaviors[appl] = []int{newPower}
		}

		if err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("Appliance: %s -> %s %s %s %s %s %s %s\n", line[0], line[1], line[2], line[3], line[4], line[5], line[6], line[7])
	}
}

/*
Prendo gli intervalli,
Trasformo l'inizio e la fine negli indici dell'array finale
Per ogni indice comrpeso tra l'inizio e la fine
Aggiungo la potenza dell'elettrodomestico alla fine
Scrivo tutto

*/
