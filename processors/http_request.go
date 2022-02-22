package processors

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/likearthian/go-pipl/data"
	"github.com/likearthian/go-pipl/util"
)

// HTTPRequest executes an HTTP request and passes along the response body.
// It is simply wrapping an http.Request and http.Client object. See the
// net/http docs for more info: https://golang.org/pkg/net/http
type HTTPRequest struct {
	Request []*http.Request
	Client  *http.Client
	ReqFunc func(data.Data) ([]*http.Request, error)
}

// NewHTTPRequest creates a new HTTPRequest and is essentially wrapping net/http's NewRequest
// function. See https://golang.org/pkg/net/http/#NewRequest
func NewHTTPRequest(method, url string, body io.Reader) (*HTTPRequest, error) {
	req, err := http.NewRequest(method, url, body)
	return &HTTPRequest{Request: []*http.Request{req}, Client: &http.Client{}}, err
}

func NewDynamicHttpRequest(reqFn func(data.Data) ([]*http.Request, error)) *HTTPRequest {
	return &HTTPRequest{Client: &http.Client{}, ReqFunc: reqFn}
}

// ProcessData sends data to outputChan if the response body is not null
func (r *HTTPRequest) ProcessData(d data.Data, outputChan chan data.Data, killChan chan error) {
	var requests []*http.Request
	var err error
	if r.ReqFunc != nil {
		requests, err = r.ReqFunc(d)
		util.KillPipelineIfErr(err, killChan)
	} else {
		requests = r.Request
	}

	for _, req := range requests {
		resp, err := r.Client.Do(req)
		util.KillPipelineIfErr(err, killChan)
		if resp != nil && resp.Body != nil {
			dd, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			util.KillPipelineIfErr(err, killChan)
			outputChan <- data.FromRawBytes(dd)
		}
	}
}

// Finish - see interface for documentation.
func (r *HTTPRequest) Finish(outputChan chan data.Data, killChan chan error) {
}

func (r *HTTPRequest) String() string {
	return "HTTPRequest"
}
