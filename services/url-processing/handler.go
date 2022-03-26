package urlprocessing

import (
	"fmt"
	"net/http"
)

func HandleURLProcessing(rw http.ResponseWriter, r *http.Response) {
	fmt.Fprint(rw, "Hello, world!")
	// comment to trigger workflow
}
