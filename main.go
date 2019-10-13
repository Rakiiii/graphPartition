package main

import (
	"fmt"
	"math/big"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"

	bigintlib "github.com/Rakiiii/goBigIntLib"
	graphlib "github.com/Rakiiii/goGraph"
	graphpartitionlib "github.com/Rakiiii/goGraphPartitionLib"
)

var wg sync.WaitGroup

func main() {

	var parser = new(graphlib.Parser)
	var g, err = parser.ParseUnweightedUndirectedGraphFromFile(os.Args[2])
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Graph parsed")

	amountOfGroups, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println(err)
		return
	}
	end := bigintlib.Pow(big.NewInt(2), int64(g.AmountOfVertex())*int64(amountOfGroups))

	start := big.NewInt(0)

	//for i := 0; i < g.AmountOfVertex(); i++ {
	//	start.Add(start, bigintlib.Pow(big.NewInt(2), int64((amountOfGroups*g.AmountOfVertex()-1)-(amountOfGroups*(i+1)-i+1))))
	//}

	fmt.Println("Big int initialized")

	var result = graphpartitionlib.Result{Matrix: nil, Value: -1}

	if os.Args[1] == "-s" {

		result, err = graphpartitionlib.FindBestPartion(g, start, end, amountOfGroups, float64(1))

	} else {
		amount := strings.Trim(os.Args[1], "-")
		am, err := strconv.Atoi(amount)
		if err != nil {
			fmt.Println(err)
			return
		}

		runtime.GOMAXPROCS(am)

		ch := make(chan graphpartitionlib.Result, am)

		wg.Add(am)

		dif := big.NewInt(0)
		dif.Sub(end, start)

		dif.Div(dif, big.NewInt(int64(am)))

		subEnd := big.NewInt(0)
		subEnd.Add(subEnd, start)
		subEnd.Add(subEnd, dif)
		for i := 0; i < am; i++ {

			go graphpartitionlib.AsyncFindBestPartion(g, big.NewInt(0), end, amountOfGroups, float64(10), &wg, ch)
			start.Add(start, dif)
			if i != am-2 {
				subEnd.Add(subEnd, dif)
			} else {
				subEnd = end
			}
		}

		wg.Wait()
		close(ch)

		for i := range ch {
			fmt.Println(i.Value)
			if result.Value < i.Value || result.Value == -1 {
				result = i
			}
		}

	}

	fmt.Println("graph partitioned")

	f, err := os.Create("result_" + os.Args[2])
	if err != nil {
		fmt.Println(err)
		fmt.Println(result.Value)
		for i := 0; i < result.Matrix.Heigh(); i++ {
			for j := 0; j < result.Matrix.Width(); j++ {
				fmt.Print(result.Matrix.GetBool(i, j))
			}
			fmt.Println()
		}
		return
	}
	defer f.Close()

	f.WriteString(strconv.FormatInt(result.Value, 10) + "\n")
	for i := 0; i < result.Matrix.Heigh(); i++ {
		subStr := ""
		for j := 0; j < result.Matrix.Width(); j++ {
			if result.Matrix.GetBool(i, j) {
				subStr = subStr + string("1 ")
			} else {
				subStr = subStr + string("0 ")
			}
		}
		subStr = subStr + "\n"
		f.WriteString(subStr)

	}

	fmt.Println("Finished")

}
