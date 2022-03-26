package urlprocessing

import (
	"fmt"
	"net/http"
)

func HandleURLProcessing(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprint(rw, "Hello, world!")
}
