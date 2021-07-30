package main

import (
	"pokedex/internal/domain"
	"pokedex/internal/client"
	"bufio"
	"fmt"
	"strings"
	"os"
	"regexp"
)

func main() {
	run()
}

func run() {
	scanner := bufio.NewScanner(os.Stdin)
	cli := client.NewPokeApiClient()

	printHeader()
	for {
		fmt.Print("Search a pokemon by name or ID, or 'exit' to quit the program: ")
		scanner.Scan()
		str := scanner.Text()
		if str == "exit" {
			break
		}
		isValid, _ := regexp.MatchString("^[a-zA-Z0-9]*$", str)
		if !isValid {
			fmt.Println("Input must be alphanumeric only.")
			continue
		}
		fmt.Println("++ RESULTS ++")
		basicInfo, err := cli.GetBasicInfo(str)
		if err != nil {
			fmt.Println("Error searching pokemon: ", err)
			continue
		}
		printBasicInfo(*basicInfo)

		encounters, err := cli.GetEncounters(str)
		if err != nil {
			fmt.Println("Error searching pokemon's location area encounters: ", err)
			continue
		}
		encounters = getKantoEncounters(encounters)
		printEncounters(encounters)
		fmt.Println("++   END   ++")
	}
}

func printHeader() {
	fmt.Println("*********************************************")
	fmt.Println()
	fmt.Println("                   POKEDEX                   ")
	fmt.Println()
	fmt.Println("*********************************************")
}

func printBasicInfo(results domain.PokeData) {
	fmt.Println("ID:      ", results.Id)
	fmt.Println("Name:    ", results.Name)
	fmt.Println("Type(s): ")
	for i, typeItem := range results.Types {
		fmt.Printf("  %d. %s\n", i+1, typeItem)
	}
	fmt.Println("Stats: ")
	for i, statItem := range results.Stats {
		fmt.Printf("  %d. %s\n", i+1, statItem.StatName)
	}
}

func printEncounters(results []domain.LocationAreaEncounter) {
	fmt.Println("Encounter Location(s) and Method(s) in Kanto: ")
	if len(results) == 0 {
		fmt.Println(" - ")
		return
	}
	for i, ec := range results {
		fmt.Printf("  %d. Name: %s\n", i+1, ec.Name)
		fmt.Println("     Version details: ")
		for _, vd := range ec.VersionDetails {
			fmt.Println("      * Max chance: ", vd.MaxChance)
			fmt.Println("        Version name: ", vd.VersionName)
			fmt.Println("        Encounter details: ")
			for _, ed := range vd.EncounterDetails {
				fmt.Println("          - Chance: ", ed.Chance)
				fmt.Println("            Max level: ", ed.MaxLevel)
				fmt.Println("            Method name: ", ed.MethodName)
			}
		}
	}
}

func getKantoEncounters(results []domain.LocationAreaEncounter) []domain.LocationAreaEncounter {
	filteredEncounters := make([]domain.LocationAreaEncounter, 0)
	for _, ec := range results {
		if strings.Contains(ec.Name, "kanto") {
			filteredEncounters = append(filteredEncounters, ec)
		}
	}
	return filteredEncounters
}