// https://adventofcode.com/2020/day/7
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
	"strings"
)

type Bag struct {
	Count  int
	Colour string
}

func parseBags(bagStr string, targetBag string) ([]Bag, bool) {
	canContain := false
	bags := make([]Bag, 0)

	bagStrParts := strings.Split(bagStr, ", ")
	for _, bag := range bagStrParts {
		if bag == "no other bags" {
			continue
		}

		// Remove the count and get only the bag name
		spaceIdx := strings.Index(bag, " ")
		bagIdx := strings.Index(bag, " bag")
		count, _ := strconv.Atoi(bag[:spaceIdx])
		bag = bag[spaceIdx+1 : bagIdx]

		if bag == targetBag {
			canContain = true
		}

		bags = append(bags, Bag{
			Count:  count,
			Colour: bag,
		})
	}

	return bags, canContain
}

func getOuterBags(bags map[string][]Bag, target string) []string {
	outerBags := make([]string, 0)

	for bag, innerBags := range bags {
		for _, innerBag := range innerBags {
			if innerBag.Colour == target {
				outerBags = append(outerBags, bag)
				break
			}
		}
	}

	return outerBags
}

func distinct(slices []string) []string {
	result := make([]string, 0)

	sort.Strings(slices)

	for idx, slice := range slices {
		if idx > 0 && slice == slices[idx-1] {
			continue
		}

		result = append(result, slice)
	}

	return result
}

func countInnerBags(bags map[string][]Bag, targetBag string, count int) int {
	innerBags := bags[targetBag]
	for _, innerBag := range innerBags {
		innerCount := countInnerBags(bags, innerBag.Colour, 0)
		count += innerBag.Count + innerBag.Count*innerCount
	}

	return count
}

func main() {
	fc, err := ioutil.ReadFile("input.txt")
	if err != nil {
		log.Println(err)
		return
	}

	lines := strings.Split(string(fc), "\n")
	bags := make(map[string][]Bag)

	targetBag := "shiny gold"
	possibleBags := make([]string, 0)

	for _, line := range lines {
		bagIdx := strings.Index(line, "bag")
		currBag := line[:bagIdx-1]

		containIdx := strings.Index(line, "contain")
		insideBags := line[containIdx+len("contain")+1 : len(line)-1] // Remove the . at the same time

		// Save the possible bags that can store the shiny gold bag together here to avoid another loop
		insides, canContain := parseBags(insideBags, targetBag)
		bags[currBag] = insides

		if canContain {
			possibleBags = append(possibleBags, currBag)
		}
	}

	possibleBags = distinct(possibleBags)

	allBags := make([]string, 0)
	allBags = append(allBags, possibleBags...)

	// Check for every possible bag, what other bags can contain it recursively to the top level of bag
	for {
		outerBags := make([]string, 0)
		for _, possibleBag := range possibleBags {
			outerBags = append(outerBags, getOuterBags(bags, possibleBag)...)
		}

		outerBags = distinct(outerBags)
		allBags = append(allBags, outerBags...)
		possibleBags = outerBags

		if len(possibleBags) == 0 {
			break
		}
	}

	allBags = distinct(allBags)
	fmt.Println(len(allBags))

	// Part 2
	count := countInnerBags(bags, targetBag, 0)
	fmt.Println(count)
}
