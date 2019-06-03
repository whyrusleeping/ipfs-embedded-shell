package embeddedShell

import (
	"context"
	"errors"
	"io"

	files "github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/interface-go-ipfs-core/path"
)

// Cat resolves the ipfs path p and returns a reader for that data, if it exists and is availalbe
func (s *Shell) Cat(p string) (io.ReadCloser, error) {
	api, err := coreapi.NewCoreAPI(s.node)
	if err != nil {
		return nil, err
	}

	f, err := api.Unixfs().Get(context.Background(), path.New(p))
	if err != nil {
		return nil, err
	}

	rf := files.ToFile(f)
	if rf == nil {
		return nil, errors.New("cannot cat a non-file")
	}

	return rf, nil
}
