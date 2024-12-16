package search

import (
	"github.com/Kuniwak/name/eval"
	"github.com/Kuniwak/name/filter"
	"github.com/Kuniwak/name/gen"
	"github.com/Kuniwak/name/mora"
	"github.com/Kuniwak/name/sex"
	"sync"
)

func Parallel(
	familyName []rune,
	in <-chan gen.Generated,
	out chan<- filter.Target,
	filterFunc filter.Func,
	strokesMap map[rune]byte,
	sexFunc sex.Func,
	parallelism int,
) {
	var wg sync.WaitGroup
	wg.Add(parallelism)
	for i := 0; i < parallelism; i++ {
		go func() {
			defer wg.Done()
			Search(familyName, in, out, filterFunc, strokesMap, sexFunc)
		}()
	}
	wg.Wait()
	close(out)
}

func Search(
	familyName []rune,
	in <-chan gen.Generated,
	out chan<- filter.Target,
	filterFunc filter.Func,
	strokesMap map[rune]byte,
	sexFunc sex.Func,
) {
	for generated := range in {
		res, err := eval.Evaluate(familyName, generated.GivenName, strokesMap)
		if err != nil {
			continue
		}

		target := filter.Target{
			Kanji:      generated.GivenName,
			Yomi:       generated.Yomi,
			YomiString: generated.YomiString,
			Strokes:    eval.SumStrokes(generated.GivenName, strokesMap),
			Mora:       mora.Count(generated.Yomi),
			Sex:        sexFunc(generated.YomiString),
			EvalResult: res,
		}
		if filterFunc(target) {
			out <- target
		}
	}
}
