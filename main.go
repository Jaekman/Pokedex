package main

import (
	"github.com/Jaekman/pokedexcli/internal/pokecache"
)

type pageState struct {
	previousURL *string
	nextURL     *string
	cache       pokecache.Cache
	userPokeDex map[string]pokemonStruct
}

func main() {
	ps := pageState{
		cache:       pokecache.NewCache(60),
		userPokeDex: make(map[string]pokemonStruct),
	}
	startRepl(&ps)
}
