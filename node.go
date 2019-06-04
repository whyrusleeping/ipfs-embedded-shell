package embeddedShell

import (
	"fmt"
	"github.com/ipfs/go-ipfs-config"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/bootstrap"
	"github.com/ipfs/go-ipfs/plugin/loader"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	"golang.org/x/net/context"
	"io/ioutil"
	"sync"
)

var plugins sync.Once
var pluginErr error

func initPlugins() error {
	// initializes preloaded datastore plugins

	plugins.Do(func() {
		pl, err := loader.NewPluginLoader("")
		if err != nil {
			pluginErr = err
			return
		}

		if err := pl.Initialize(); err != nil {
			pluginErr = err
			return
		}

		if err := pl.Inject(); err != nil {
			pluginErr = err
			return
		}
	})
	return pluginErr
}

func NewDefaultNodeWithFSRepo(ctx context.Context, repoPath string) (*core.IpfsNode, error) {
	if err := initPlugins(); err != nil {
		return nil, err
	}

	r, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, fmt.Errorf("opening fsrepo failed: %s", err)
	}
	n, err := core.NewNode(ctx, &core.BuildCfg{
		Online: true,
		Repo:   r,
	})
	if err != nil {
		return nil, fmt.Errorf("ipfs NewNode() failed: %s", err)
	}
	// TODO: can we bootsrap localy/mdns first and fall back to default?
	err = n.Bootstrap(bootstrap.DefaultBootstrapConfig)
	if err != nil {
		return nil, fmt.Errorf("ipfs Bootstrap() failed: %s", err)
	}
	return n, nil
}

func NewTmpDirNode(ctx context.Context) (*core.IpfsNode, error) {
	dir, err := ioutil.TempDir("", "ipfs-shell")
	if err != nil {
		return nil, fmt.Errorf("failed to get temp dir: %s", err)
	}

	cfg, err := config.Init(ioutil.Discard, 1024)
	if err != nil {
		return nil, err
	}

	if err := initPlugins(); err != nil {
		return nil, err
	}

	err = fsrepo.Init(dir, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init ephemeral node: %s", err)
	}

	repo, err := fsrepo.Open(dir)
	if err != nil {
		return nil, err
	}

	return core.NewNode(ctx, &core.BuildCfg{
		Online: true,
		Repo:   repo,
	})
}
