package embeddedShell

import (
	"context"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"strings"
	"testing"

	cid "github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs/assets"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/mock"
)

type testShell struct {
	cid cid.Cid
	mn  *core.IpfsNode
	s   *Shell
}

func newTestShell(t *testing.T) *testShell {
	mn, err := coremock.NewMockNode()
	if err != nil {
		t.Fatalf("coremock.NewMockNode() failed: %s", err)
	}
	tk, err := assets.SeedInitDocs(mn)
	if err != nil {
		t.Fatalf("assets.SeedInitDocs() failed: %s", err)
	}
	return &testShell{
		cid: tk,
		mn:  mn,
		s:   NewShell(mn),
	}
}

func TestCat(t *testing.T) {
	ts := newTestShell(t)
	rc, err := ts.s.Cat(path.Join("/ipfs/", ts.cid.String(), "about"))
	if err != nil {
		t.Fatal(err)
	}
	ioutil.ReadAll(rc)
	h := sha1.New()
	_, err = io.Copy(h, rc)
	if err != nil {
		t.Fatal(err)
	}
	got := h.Sum(nil)
	want := "da39a3ee5e6b4b0d3255bfef95601890afd80709"
	if want != fmt.Sprintf("%x", got) {
		t.Errorf("hash comparison failed\nWant: %s\nGot:  %x", want, got)
	}
}

func TestAdd(t *testing.T) {
	ts := newTestShell(t)
	h, err := ts.s.Add(strings.NewReader("Hello, World"))
	if err != nil {
		t.Fatal(err)
	}
	if h != "QmTev1ZgJkHgFYiCX7MgELEDJuMygPNGcinqBa2RmfnGFu" {
		t.Fatal("wrong hash from add")
	}
}

func TestTempNode(t *testing.T) {
	ctx := context.Background()

	s, err := NewTmpDirNode(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
}
