package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func downloadPkgs(deps ManifestDeps) error {
	srcPath := getSrcPath(setting.WorkDir())
	for importPath, rev := range deps {
		showRev := rev
		if showRev == "" {
			showRev = "latest"
		}

		fmt.Printf("Downloading %s@%s\n", importPath, showRev)
		err := downloadPkg(srcPath, importPath, rev)
		if err != nil {
			return err
		}
	}

	return nil
}

// downloadPkg downloads package and checkouts revision to dir
func downloadPkg(dir, importPath, rev string) error {
	root, vcs, err := VCSForImportPath(importPath)
	if err != nil {
		return err
	}

	pkgRoot := filepath.Join(dir, filepath.FromSlash(root.Root))
	err = os.MkdirAll(filepath.Dir(pkgRoot), os.ModePerm)
	if err != nil {
		return err
	}

	if rev == "" {
		err = vcs.Create(pkgRoot, root.Repo)
	} else {
		err = vcs.CreateAtRev(pkgRoot, root.Repo, rev)
	}

	return nil
}
