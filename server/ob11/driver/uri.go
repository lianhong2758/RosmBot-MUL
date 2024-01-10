package driver

import (
	"encoding/base64"
	"net/url"
	"strings"

	"github.com/lianhong2758/RosmBot-MUL/tool"
)

func resolveURI(addr string) (network, address string) {
	network, address = "tcp", addr
	uri, err := url.Parse(addr)
	if err == nil && uri.Scheme != "" {
		scheme, ext, _ := strings.Cut(uri.Scheme, "+")
		if ext != "" {
			network = ext
			uri.Scheme = scheme // remove `+unix`/`+tcp4`
			if ext == "unix" {
				uri.Host, uri.Path, _ = strings.Cut(uri.Path, ":")
				uri.Host = base64.StdEncoding.EncodeToString(tool.StringToBytes(uri.Host)) // special handle for unix
			}
			address = uri.String()
		}
	}
	return
}
