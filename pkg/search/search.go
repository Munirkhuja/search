package search

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

type Result struct {
	Phrase  string
	Line    string
	LineNum int64
	ColNum  int64
}

func All(ctx context.Context,phrase string,files [] string)<-chan []Result{
	goroutinesCount:=len(files)
	wg:=sync.WaitGroup{}
	wg.Add(goroutinesCount)
	ch:=make(chan []Result)
	defer close(ch)
	for i := 0; i < goroutinesCount; i++ {
		go func (ch chan<-[]Result,path string,findPhrase string)  {
			defer wg.Done()
			file, err := os.Open(path)
			lineNum:=int64(0)
			if err == nil {
				defer func() {
					if cerr := file.Close(); cerr != nil {
						log.Print(cerr)
					}
				}()
				var result []Result
				reader := bufio.NewReader(file)
				for {
					lineNum++
					line, berr := reader.ReadString('\n')
					colNum:=int64(strings.Index(line,findPhrase))
					if colNum>=0 {
						line = strings.TrimSuffix(line, "\n")
						line = strings.TrimSuffix(line, "\r")
						result = append(result, Result{Phrase: findPhrase,Line: line,LineNum: lineNum,ColNum: colNum})
					}

					if berr == io.EOF {
						break
					}
				}
				ch<-result
			}
		}(ch,files[i],phrase)
	}
	wg.Wait()
	<-ctx.Done()
	return ch
}
func Any(ctx context.Context,phrase string,files [] string)<-chan []Result{
	ch:=make(chan []Result)
	return ch
}