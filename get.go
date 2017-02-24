package embeddedShell

import (
	"errors"
	"gopkg.in/errgo.v1"

	"github.com/ipfs/go-ipfs/core"
	dag "github.com/ipfs/go-ipfs/merkledag"
	"github.com/ipfs/go-ipfs/path"
	tar "github.com/ipfs/go-ipfs/thirdparty/tar"
	uarchive "github.com/ipfs/go-ipfs/unixfs/archive"
)

// Cat resolves the ipfs path p and returns a reader for that data, if it exists and is availalbe
func (s *Shell) Get(ref, outdir string) error {
	ipfsPath, err := path.ParsePath(ref)
	if err != nil {
		return errgo.Notef(err, "get: could not parse %q", ref)
	}

	nd, err := core.Resolve(s.ctx, s.node.Namesys, s.node.Resolver, ipfsPath)
	if err != nil {
		return errgo.Notef(err, "get: could not resolve %s", ipfsPath)
	}

	pbnd, ok := nd.(*dag.ProtoNode)
	if !ok {
		return errors.New("could not cast Node to ProtoNode")
	}

	r, err := uarchive.DagArchive(s.ctx, pbnd, outdir, s.node.DAG, false, 0)
	if err != nil {
		return err
	}

	ext := tar.Extractor{outdir}

	return ext.Extract(r)
}
