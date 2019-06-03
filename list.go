package embeddedShell

import (
	"context"

	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/interface-go-ipfs-core/path"

	// for types
	sh "github.com/ipfs/go-ipfs-api"
)

func (s *Shell) List(ipath string) ([]*sh.LsLink, error) {
	api, err := coreapi.NewCoreAPI(s.node)
	if err != nil {
		return nil, err
	}

	ls, err := api.Unixfs().Ls(context.Background(), path.New(ipath))
	if err != nil {
		return nil, err
	}

	var out []*sh.LsLink
	for l := range ls {
		out = append(out, &sh.LsLink{
			Hash: l.Cid.String(),
			Name: l.Name,
			Size: l.Size,
		})
	}

	return out, nil
}

func (s *Shell) ResolvePath(ipath string) (string, error) {
	api, err := coreapi.NewCoreAPI(s.node)
	if err != nil {
		return "", err
	}

	rp, err := api.ResolvePath(context.Background(), path.New(ipath))
	if err != nil {
		return "", err
	}

	return rp.Cid().String(), nil
}
