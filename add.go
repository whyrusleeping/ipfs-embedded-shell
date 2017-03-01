package embeddedShell

import (
	"fmt"
	"io"

	"github.com/ipfs/go-ipfs/importer"
	"github.com/ipfs/go-ipfs/importer/chunk"
	dag "github.com/ipfs/go-ipfs/merkledag"
	ft "github.com/ipfs/go-ipfs/unixfs"
)

func (s *Shell) Add(r io.Reader) (string, error) {
	dag, err := importer.BuildDagFromReader(
		s.node.DAG,
		chunk.DefaultSplitter(r),
	)
	if err != nil {
		return "", fmt.Errorf("add: importing DAG failed: %s", err)
	}
	return dag.Cid().String(), nil
}

// AddLink creates a unixfs symlink and returns its hash
func (s *Shell) AddLink(target string) (string, error) {
	d, _ := ft.SymlinkData(target)
	nd := dag.NodeWithData(d)
	c, err := s.node.DAG.Add(nd)
	if err != nil {
		return "", err
	}
	return c.String(), nil
}
