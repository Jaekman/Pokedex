package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(ps *pageState, args string) error
}

func commandExit(ps *pageState, args string) error {
	fmt.Println("Exiting the PokeDex")
	os.Exit(0)
	return nil
}

func commandHelp(ps *pageState, args string) error {
	fmt.Println("Usage:")
	for _, cmd := range cliCommands() {
		fmt.Println(cmd.name, "->", cmd.description)
	}
	return nil
}

func commandMap(ps *pageState, args string) error {

	var query string
	if ps.nextURL != nil {
		query = *ps.nextURL
	}

	locationAreaURL := pokeApiURL{
		resource: "/location-area/",
		query:    query,
	}

	printAndSetURL := func() error {
		API, err := locationAreaAPI(baseURL, locationAreaURL, ps.cache)
		if err != nil {
			return err
		}
		for location := range API.Results {
			fmt.Println(API.Results[location].Name)
		}
		if qr, err := urlQueryParse(API.Previous); err != nil {
			return err
		} else {
			ps.previousURL = &qr
		}
		if qr, err := urlQueryParse(API.Next); err != nil {
			return err
		} else {
			ps.nextURL = &qr
		}
		return nil
	}

	err := printAndSetURL()
	if err != nil {
		return err
	}

	return err
}

func commmandMapb(ps *pageState, args string) error {
	var query string
	if ps.previousURL != nil {
		query = *ps.previousURL
	} else {
		return errors.New("you're on the first page")
	}

	locationAreaURL := pokeApiURL{
		resource: "/location-area/",
		query:    query,
	}

	printAndSetURL := func() error {
		API, err := locationAreaAPI(baseURL, locationAreaURL, ps.cache)
		if err != nil {
			return err
		}
		for location := range API.Results {
			fmt.Println(API.Results[location].Name)
		}
		if qr, err := urlQueryParse(API.Previous); err != nil {
			return err
		} else {
			ps.previousURL = &qr
		}
		if qr, err := urlQueryParse(API.Next); err != nil {
			return err
		} else {
			ps.nextURL = &qr
		}
		return nil
	}

	err := printAndSetURL()
	if err != nil {
		return err
	}

	return err
}

func commandExplore(ps *pageState, args string) error {
	fmt.Printf("Exploring %s...\n", args)
	exploreAreaUrl := pokeApiURL{
		resource: "/location-area/",
		query:    args,
	}

	printAndSetURL := func() error {
		API, err := pokeAreaAPI(baseURL, exploreAreaUrl, ps.cache)
		if err != nil {
			return err
		}
		for pokemon := range API.PokemonEncounters {
			fmt.Println(API.PokemonEncounters[pokemon].Pokemon.Name)
		}
		return nil
	}
	err := printAndSetURL()
	if err != nil {
		return err
	}

	return err
}

func commandCatch(ps *pageState, args string) error {
	fmt.Printf("Throwing a Pokeball at %s...\n", args)
	pokemonUrl := pokeApiURL{
		resource: "/pokemon/",
		query:    args,
	}

	printAndSetURL := func() error {
		API, err := pokemonAPI(baseURL, pokemonUrl, ps.cache)
		if err != nil {
			return err
		}
		pokeXP := API.BaseExperience
		catchChance := rand.Intn(700)

		if catchChance > pokeXP {
			ps.userPokeDex[args] = API
			fmt.Printf("You caught the %s!\n", args)
		} else {
			fmt.Printf("%s got away!\n", args)
		}
		return nil
	}
	err := printAndSetURL()
	if err != nil {
		return err
	}

	return err
}

func commandInspect(ps *pageState, args string) error {
	fmt.Printf("Inspecting %s...\n", args)
	if _, ok := ps.userPokeDex[args]; !ok {
		return errors.New("you haven't caught that Pokemon yet")
	}
	pokemon := ps.userPokeDex[args]

	fmt.Printf("Name: %s\nHeight: %d\nWeight: %d\n", pokemon.Name, pokemon.Height, pokemon.Weight)

	for _, stat := range pokemon.Stats {
		fmt.Printf("%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	for _, types := range pokemon.Types {
		fmt.Printf("Type: %s\n", types.Type.Name)
	}

	return nil
}
func commandPokedex(ps *pageState, args string) error {
	fmt.Printf("Your Pokedex:\n")
	for pokemon := range ps.userPokeDex {
		fmt.Println(pokemon)

	}
	return nil
}

func urlQueryParse(urlStr *string) (string, error) {
	if urlStr == nil {
		return "", nil
	}
	u, err := url.Parse(*urlStr)
	if err != nil {
		return "", err
	}
	fullQuery := "?" + u.RawQuery
	return fullQuery, nil
}

func cliCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exits the PokeDex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Diplays PokeDex Commands",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Print 20 Location Areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Go back 20 Location Areas",
			callback:    commmandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Explore the Location Areas",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch a Pokemon in the Location Area",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "View a caught Pokemon's stats",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "View your caught Pokemon",
			callback:    commandPokedex,
		},
	}
}
