package embeddedShell

import (
	"errors"
	"fmt"

	core "github.com/ipfs/go-ipfs/core"
	dagutils "github.com/ipfs/go-ipfs/dagutils"
	dag "github.com/ipfs/go-merkledag"
	"github.com/ipfs/go-path"
	unixfs "github.com/ipfs/go-unixfs"
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
	p, err := path.ParsePath(root)
	if err != nil {
		return "", err
	}

	nd, err := core.Resolve(s.ctx, s.node.Namesys, s.node.Resolver, p)
	if err != nil {
		return "", err
	}

	rootnd, ok := nd.(*dag.ProtoNode)
	if !ok {
		return "", errors.New("could not cast Node to ProtoNode")
	}

	insertpath := args[0]
	childhash := args[1]

	childpath, err := path.ParsePath(childhash)
	if err != nil {
		return "", err
	}

	nnode, err := core.Resolve(s.ctx, s.node.Namesys, s.node.Resolver, childpath)
	if err != nil {
		return "", err
	}

	e := dagutils.NewDagEditor(rootnd, s.node.DAG)

	switch action {
	case "add-link":
		err := e.InsertNodeAtPath(s.ctx, insertpath, nnode, nil)
		if err != nil {
			return "", err
		}

		_, err = e.Finalize(s.ctx, s.node.DAG)
		if err != nil {
			return "", err
		}

		return e.GetNode().Cid().String(), nil
	default:
		return "", fmt.Errorf("unsupported action (impl not complete)")
	}
}

//TODO: hrm, maybe this interface could be better
func (s *Shell) PatchLink(root, npath, childhash string, create bool) (string, error) {
	p, err := path.ParsePath(root)
	if err != nil {
		return "", err
	}

	nd, err := core.Resolve(s.ctx, s.node.Namesys, s.node.Resolver, p)
	if err != nil {
		return "", err
	}

	rootnd, ok := nd.(*dag.ProtoNode)
	if !ok {
		return "", errors.New("could not cast Node to ProtoNode")
	}

	childpath, err := path.ParsePath(childhash)
	if err != nil {
		return "", err
	}

	nnode, err := core.Resolve(s.ctx, s.node.Namesys, s.node.Resolver, childpath)
	if err != nil {
		return "", err
	}

	e := dagutils.NewDagEditor(rootnd, s.node.DAG)
	err = e.InsertNodeAtPath(s.ctx, npath, nnode, func() *dag.ProtoNode {
		return dag.NodeWithData(unixfs.FolderPBData())
	})
	if err != nil {
		return "", err
	}

	_, err = e.Finalize(s.ctx, s.node.DAG)
	if err != nil {
		return "", err
	}

	return e.GetNode().Cid().String(), nil
}
