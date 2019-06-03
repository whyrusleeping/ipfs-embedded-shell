package embeddedShell

import (
	"context"
	"fmt"

	"github.com/ipfs/go-ipfs/core/coreapi"
	dag "github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-unixfs"
	"github.com/ipfs/interface-go-ipfs-core/options"
	"github.com/ipfs/interface-go-ipfs-core/path"
)

func (s *Shell) NewObject(template string) (string, error) {
	node := new(dag.ProtoNode)
	switch template {
	case "":
		break
	case "unixfs-dir":
		node.SetData(unixfs.FolderPBData())
	default:
		return "", fmt.Errorf("unknown template %s", template)
	}
	err := s.node.DAG.Add(s.ctx, node)
	if err != nil {
		return "", err
	}

	return node.Cid().String(), nil
}

// TODO: extract all this logic from the core/commands/object.go to avoid dupe code
func (s *Shell) Patch(root, action string, args ...string) (string, error) {
	api, err := coreapi.NewCoreAPI(s.node)
	if err != nil {
		return "", err
	}

	insertpath := args[0]
	childhash := path.New(args[1])

	switch action {
	case "add-link":
		h, err := api.Object().AddLink(context.Background(), path.New(root), insertpath, childhash)
		if err != nil {
			return "", err
		}

		return h.Cid().String(), nil
	default:
		return "", fmt.Errorf("unsupported action (impl not complete)")
	}
}

//TODO: hrm, maybe this interface could be better
func (s *Shell) PatchLink(root, npath, childhash string, create bool) (string, error) {
	api, err := coreapi.NewCoreAPI(s.node)
	if err != nil {
		return "", err
	}

	h, err := api.Object().AddLink(context.Background(), path.New(root), npath, path.New(childhash), options.Object.Create(create))
	if err != nil {
		return "", err
	}

	return h.Cid().String(), nil
}
