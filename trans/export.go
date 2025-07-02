package trans

import "errors"

type ExportMode int

const (
	Default ExportMode = iota
	Quick
	Dev
	Slim
	Tweakable
)

func SwapServerGitToDev() error {
	return errors.New("not yet implemented")
}

func ResolveServerRemovals() error {
	return errors.New("not yet implemented")
}

func DoExportJsonTransforms(em ExportMode) error {
	return errors.New("not yet implemented")
}
