package main

import (
	"pokedex/internal/domain"
	"pokedex/internal/client"
	"pokedex/internal/cache"
	"bufio"
	"encoding/json"
	"fmt"
	"time"
	"strings"
	"strconv"
	"os"
	"regexp"
)

func main() {
	run()
}

func run() {
	scanner := bufio.NewScanner(os.Stdin)
	cli := client.NewPokeApiClient()
	c := cache.New(1 * time.Minute)
	c.LoadFile("cache/cache.txt")

	printHeader()
	for {
		fmt.Print("Search a pokemon by name or ID, or 'exit' to quit the program: ")
		scanner.Scan()
		identifier := strings.ToLower(scanner.Text())
		if identifier == "exit" {
			break
		}
		isValid, _ := regexp.MatchString("^[a-zA-Z0-9]*$", identifier)
		if !isValid || len(identifier) == 0 {
			fmt.Println("Input must be alphanumeric only.")
			continue
		}

		_, err := strconv.ParseInt(identifier, 10, 64)
		isNotId := err != nil
		isCached := false
		data := &domain.PokeDataWithLocation{}
		if isNotId {
			if id, isFound := c.Get(identifier); isFound {
				identifier = string(id.([]byte)[:])
			}
		}
		if cachedData, cacheFound := c.Get(identifier); cacheFound {
			pokeData, ok := cachedData.([]byte)
			if ok {
				json.Unmarshal(pokeData, data)
				isCached = true
			}
		}

		fmt.Println("++ RESULTS ++")
		if !isCached {
			basicInfo, err := cli.GetBasicInfo(identifier)
			if err != nil {
				fmt.Println("Error searching pokemon: ", err)
				continue
			}
			data.PokeData = *basicInfo
		}

		printBasicInfo(data.PokeData)

		if !isCached {
			encounters, err := cli.GetEncounters(identifier)
			if err != nil {
				fmt.Println("Error searching pokemon's location area encounters: ", err)
				continue
			}
			encounters = getKantoEncounters(encounters)
			data.Encounters = encounters
		}
		
		printEncounters(data.Encounters)
		fmt.Println("++   END   ++")

		if !isCached {
			c.Set(data.PokeData.Name, fmt.Sprint(data.PokeData.Id), cache.NoExpiration)
			
			dataInString, _ := json.Marshal(data)
			c.Set(
				fmt.Sprint(data.PokeData.Id),
				[]byte(dataInString),
				cache.DefaultExpiration,
			)
		}
	}
	fmt.Println("See you again! Caching data for next time...")
	err := c.Write("cache/cache.txt")
	if err != nil {
		fmt.Println("Error caching data: ", err)
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