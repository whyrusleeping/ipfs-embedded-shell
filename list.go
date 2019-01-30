package embeddedShell

import (
	core "github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-path"

	// for types
	sh "github.com/ipfs/go-ipfs-api"
)

func (s *Shell) List(ipath string) ([]*sh.LsLink, error) {
	p, err := path.ParsePath(ipath)
	if err != nil {
		return nil, err
	}

	nd, err := core.Resolve(s.ctx, s.node.Namesys, s.node.Resolver, p)
	if err != nil {
		return nil, err
	}

	var out []*sh.LsLink
	for _, l := range nd.Links() {
		out = append(out, &sh.LsLink{
			Hash: l.Cid.String(),
			Name: l.Name,
			Size: l.Size,
		})
	}

	return out, nil
}

func (s *Shell) ResolvePath(ipath string) (string, error) {
	p, err := path.ParsePath(ipath)
	if err != nil {
		return "", err
	}

	nd, err := core.Resolve(s.ctx, s.node.Namesys, s.node.Resolver, p)
	if err != nil {
		return "", err
	}

	return nd.Cid().String(), nil
}
