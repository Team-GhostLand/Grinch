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
	if em == Default || em == Quick {
		return nil //They don't have any JSON transforms - early-return
	}

	if em == Dev {
		return errors.New("not yet implemented")
	}

	return errors.New("not yet implemented") //We don't care about the other modes yet - I'm under a THIGHT deadline.
}
