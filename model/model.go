// Package model provides the complete representation of the model for a given GEP problem.
package model

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/gmlewis/gep/functions"
	"github.com/gmlewis/gep/gene"
	"github.com/gmlewis/gep/genome"
)

// ScoringFunc is the function that is used to evaluate the fitness of the model.
// Typically, a return value of 0 means that the function is nowhere close to being
// a valid solution and a return value of 1000 (or higher) means a perfect solution.
type ScoringFunc func(g *genome.Genome) float64

// Generation represents one complete generation of the model.
type Generation struct {
	Genomes     []*genome.Genome
	Funcs       []gene.FuncWeight
	ScoringFunc ScoringFunc
}

// New creates a new random generation of the model.
// fs is a slice of function weights.
// fm is the map of available functions to use for creating the generation of the model.
// numGenomes is the number of genomes to use to populate this generation of the model.
// headSize is the number of head symbols to use in a genome.
// numGenesPerGenome is the number of genes to use per genome.
// numTerminals is the number of terminals (inputs) to use within each gene.
// linkFunc is the linking function used to combine the genes within a genome.
// sf is the scoring (or fitness) function.
func New(fs []gene.FuncWeight, fm functions.FuncMap, numGenomes, headSize, numGenesPerGenome, numTerminals int, linkFunc string, sf ScoringFunc) *Generation {
	r := &Generation{
		Genomes:     make([]*genome.Genome, numGenomes, numGenomes),
		Funcs:       fs,
		ScoringFunc: sf,
	}
	n := maxArity(fs, fm)
	tailSize := headSize*(n-1) + 1
	for i := range r.Genomes {
		genes := make([]*gene.Gene, numGenesPerGenome, numGenesPerGenome)
		for j := range genes {
			genes[j] = gene.RandomNew(headSize, tailSize, numTerminals, fs)
		}
		r.Genomes[i] = genome.New(genes, linkFunc)
	}
	return r
}

// Evolve runs the GEP algorithm for the given number of iterations, or until a score of 1000 (or more) is reached.
func (g *Generation) Evolve(iterations int) *genome.Genome {
	// Algorithm flow diagram, figure 3.1, book page 56
	for i := 0; i < iterations; i++ {
		// fmt.Printf("Iteration #%v...\n", i)
		bestGenome := g.getBest() // Preserve the best genome
		if bestGenome.Score >= 1000.0 {
			fmt.Printf("Stopping after generation #%v\n", i)
			return bestGenome
		}
		// fmt.Printf("Best genome (score %v): %v\n", bestGenome.Score, *bestGenome)
		saveCopy := bestGenome.Dup()
		g.replication() // Section 3.3.1, book page 75
		g.mutation()    // Section 3.3.2, book page 77
		// g.isTransposition()
		// g.risTransposition()
		// g.geneTransposition()
		// g.onePointRecombination()
		// g.twoPointRecombination()
		// g.geneRecombination()
		// Now that replication is done, restore the best genome (aka "elitism")
		g.Genomes[0] = saveCopy
	}
	fmt.Printf("Stopping after generation #%v\n", iterations)
	return g.getBest()
}

func (g *Generation) replication() {
	// roulette wheel selection - see www.youtube.com/watch?v=aHLslaWO-AQ
	maxWeight := 0.0
	for _, v := range g.Genomes {
		if v.Score > maxWeight {
			maxWeight = v.Score
		}
	}
	result := make([]*genome.Genome, 0, len(g.Genomes))
	index := rand.Intn(len(g.Genomes))
	beta := 0.0
	for i := 0; i < len(g.Genomes); i++ {
		beta += rand.Float64() * 2.0 * maxWeight
		for beta > g.Genomes[index].Score {
			beta -= g.Genomes[index].Score
			index = (index + 1) % len(g.Genomes)
		}
		result = append(result, g.Genomes[index].Dup())
	}
	g.Genomes = result
}

func (g *Generation) mutation() {
	// Determine the total number of genomes to mutate
	numGenomes := 1 + rand.Intn(len(g.Genomes)-1)
	for i := 0; i < numGenomes; i++ {
		// Pick a random genome
		genomeNum := rand.Intn(len(g.Genomes))
		genome := &g.Genomes[genomeNum]
		// Determine the total number of mutations to perform within the genome
		numMutations := 1 + rand.Intn(2)
		// fmt.Printf("\nMutating genome #%v %v times, before:\n%v\n", genomeNum, numMutations, genome)
		genome.Mutate(numMutations)
		// fmt.Printf("after:\n%v\n", genome)
	}
}

// getBest evaluates all genomes and returns a pointer to the best one.
func (g *Generation) getBest() *genome.Genome {
	bestScore := 0.0
	bestGenome := g.Genomes[0]
	for i := 0; i < len(g.Genomes); i++ {
		g.Genomes[i].Score = g.ScoringFunc(g.Genomes[i])
		if g.Genomes[i].Score > bestScore {
			bestGenome = g.Genomes[i]
			bestScore = g.Genomes[i].Score
		}
	}
	return bestGenome
}

// maxArity determines the maximum number of input terminals for the given set of symbols.
func maxArity(fs []gene.FuncWeight, fm functions.FuncMap) int {
	r := 0
	for _, f := range fs {
		if fn, ok := fm[f.Symbol]; ok {
			if fn.Terminals() > r {
				r = fn.Terminals()
			}
		} else {
			log.Printf("unable to find symbol %v in function map\n", f.Symbol)
		}
	}
	return r
}
