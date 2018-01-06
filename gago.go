package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"

	"github.com/MaxHalford/gago"
)

var (
	corpus = strings.Split("abcdefghijklmnopqrstuvwxyz ", "")
	target = strings.Split("hello andy this isnt going to work", "")
)

// Strings is a slice of strings.
type Strings []string

// Evaluate a Strings slice by counting the number of mismatches between itself
// and the target string.
func (X Strings) Evaluate() (mismatches float64) {
	for i, s := range X {
		if s != target[i] {
			mismatches++
		}
	}
	return
}

// Mutate a Strings slice by replacing it's elements by random characters
// contained in  a corpus.
func (X Strings) Mutate(rng *rand.Rand) {
	gago.MutUniformString(X, corpus, 3, rng)
}

// Crossover a Strings slice with another by applying 2-point crossover.
func (X Strings) Crossover(Y gago.Genome, rng *rand.Rand) {
	gago.CrossGNXString(X, Y.(Strings), 2, rng)
}

// MakeStrings creates random Strings slices by picking random characters from a
// corpus.
func MakeStrings(rng *rand.Rand) gago.Genome {
	return Strings(gago.InitUnifString(len(target), corpus, rng))
}

// Clone a Strings slice..
func (X Strings) Clone() gago.Genome {
	var XX = make(Strings, len(X))
	copy(XX, X)
	return XX
}

func run(phrase string) {

	var ga = gago.Generational(MakeStrings)
	ga.Initialize()

	for i := 1; i < 300; i++ {
		ga.Evolve()
		// Concatenate the elements from the best individual and display the result
		var buffer bytes.Buffer
		for _, letter := range ga.HallOfFame[0].Genome.(Strings) {
			buffer.WriteString(letter)
		}
		fmt.Printf("Result -> %s (%.0f mismatches)\n", buffer.String(), ga.HallOfFame[0].Fitness)
	}
}