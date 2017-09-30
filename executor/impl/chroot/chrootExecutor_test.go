package chroot

import (
	"context"
	"testing"

	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/go-timeless-api/repeatr"
	"go.polydawn.net/go-timeless-api/rio"
	. "go.polydawn.net/repeatr/testutil"
	"go.polydawn.net/rio/client"
	"go.polydawn.net/rio/fs"
	"go.polydawn.net/rio/fs/osfs"
	"go.polydawn.net/rio/stitch"
)

// Base formula full of sensible defaults and ready to run:
var baseFormula = api.Formula{
	Inputs: map[api.AbsPath]api.WareID{
		"/": api.WareID{"tar", "6q7G4hWr283FpTa5Lf8heVqw9t97b5VoMU6AGszuBYAz9EzQdeHVFAou7c4W9vFcQ6"},
	},
	Action: api.FormulaAction{
		Exec: []string{"/bin/echo", "hello world"},
	},
	Outputs: map[api.AbsPath]api.OutputSpec{
		"/": {PackType: "tar", Filters: api.Filter_NoMutation},
	},
	FetchUrls: map[api.AbsPath][]api.WarehouseAddr{
		"/": []api.WarehouseAddr{
			"file://../../../fixtures/busybash.tgz",
		},
	},
}

func TestChrootExecutor(t *testing.T) {
	var (
		unpackTool rio.UnpackFunc = rioexecclient.UnpackFunc
		packTool   rio.PackFunc   = rioexecclient.PackFunc
	)

	WithTmpdir(func(tmpDir fs.AbsolutePath) {
		// Setup assembler and executor.  Both are reusable.
		asm, err := stitch.NewAssembler(unpackTool)
		AssertNoError(t, err)
		exe := Executor{
			osfs.New(tmpDir.Join(fs.MustRelPath("ws"))),
			asm,
			packTool,
		}

		t.Run("hello-world formula should work", func(t *testing.T) {
			frm := baseFormula.Clone()

			rr, err := exe.Run(context.Background(), frm, repeatr.InputControl{}, repeatr.Monitor{})
			WantNoError(t, err)
			t.Logf("%v\n", rr)
		})
	})
}
