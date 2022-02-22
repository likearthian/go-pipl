package processors

import "github.com/likearthian/go-pipl/data"

// FuncTransformer executes the given function on each data
// payload, sending the resuling data to the next stage.
//
// While FuncTransformer is useful for simple data transformation, more
// complicated tasks justify building a custom implementation of DataProcessor.
type FuncTransformer struct {
	transform        func(d data.Data) data.Data
	Name             string // can be set for more useful log output
	ConcurrencyLevel int    // See ConcurrentDataProcessor
}

//NewFuncTransformer instantiates a new instance of func transformer
func NewFuncTransformer(transform func(d data.Data) data.Data) *FuncTransformer {
	return &FuncTransformer{transform: transform}
}

// ProcessData runs the supplied func and sends the returned value to outputChan
func (t *FuncTransformer) ProcessData(d data.Data, outputChan chan data.Data, killChan chan error) {
	outputChan <- t.transform(d)
}

// Finish - see interface for documentation.
func (t *FuncTransformer) Finish(outputChan chan data.Data, killChan chan error) {
}

func (t *FuncTransformer) String() string {
	if t.Name != "" {
		return t.Name
	}
	return "FuncTransformer"
}

func (t *FuncTransformer) WithName(name string) *FuncTransformer {
	t.Name = name
	return t
}

// Concurrency defers to ConcurrentDataProcessor
func (t *FuncTransformer) Concurrency() int {
	return t.ConcurrencyLevel
}
