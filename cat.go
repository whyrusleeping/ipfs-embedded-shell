package embeddedShell

import (
	"fmt"
	"io"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/path"
	unixfsio "github.com/ipfs/go-ipfs/unixfs/io"
)

// Cat resolves the ipfs path p and returns a reader for that data, if it exists and is availalbe
func (s *Shell) Cat(p string) (io.ReadCloser, error) {
	ipfsPath, err := path.ParsePath(p)
	if err != nil {
		return nil, fmt.Errorf("cat: could not parse %q: %s", p, err)
	}
	nd, err := core.Resolve(s.ctx, s.node.Namesys, s.node.Resolver, ipfsPath)
	if err != nil {
		return nil, fmt.Errorf("cat: could not resolve %s: %s", ipfsPath, err)
	}
	dr, err := unixfsio.NewDagReader(s.ctx, nd, s.node.DAG)
	if err != nil {
		return nil, fmt.Errorf("cat: failed to construct DAG reader: %s", err)
	}
	return dr, nil
}
