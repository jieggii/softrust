package nvdadapter

import "net/http"

type Adapter struct {
	client *http.Client
}
