package processors

import (
	"io/ioutil"

	"github.com/likearthian/go-pipl/data"
	"github.com/likearthian/go-pipl/util"
)

// FileReader opens and reads the contents of the given filename.
type FileReader struct {
	filename  string
	fnameFunc func(d data.Data) string
}

// NewFileReader returns a new FileReader that will read the entire contents
// of the given file path and send it at once. For buffered or line-by-line
// reading try using IoReader.
func NewFileReader(filename string) *FileReader {
	return &FileReader{filename: filename}
}

func NewDynamicFileReader(getFileName func(d data.Data) string) *FileReader {
	return &FileReader{fnameFunc: getFileName}
}

// ProcessData reads a file and sends its contents to outputChan
func (r *FileReader) ProcessData(d data.Data, outputChan chan data.Data, killChan chan error) {
	fname := r.filename
	if r.fnameFunc != nil {
		fname = r.fnameFunc(d)
	}

	buf, err := ioutil.ReadFile(fname)
	util.KillPipelineIfErr(err, killChan)
	outputChan <- data.FromRawBytes(buf)
}

// Finish - see interface for documentation.
func (r *FileReader) Finish(outputChan chan data.Data, killChan chan error) {
}

func (r *FileReader) String() string {
	return "FileReader"
}
