package embeddedShell

import (
	"context"
	"io"

	files "github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/interface-go-ipfs-core/path"

	"github.com/ipfs/go-ipfs/core/coreapi"
	tar "github.com/whyrusleeping/tar-utils"
)

// Cat resolves the ipfs path p and returns a reader for that data, if it exists and is availalbe
func (s *Shell) Get(ref, outdir string) error {
	api, err := coreapi.NewCoreAPI(s.node)
	if err != nil {
		return err
	}

	f, err := api.Unixfs().Get(context.Background(), path.New(ref))
	if err != nil {
		return err
	}

	pr, pw := io.Pipe()
	go func() {
		tw, err := files.NewTarWriter(pw)
		if err != nil {
			_ = pw.CloseWithError(err)
			return
		}
		_ = pw.CloseWithError(tw.WriteFile(f, outdir))
	}()

	ext := tar.Extractor{Path: outdir}
	return ext.Extract(pr)
}
