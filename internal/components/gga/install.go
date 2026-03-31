package gga

import (
	"github.com/Thrasno/conpas-ai/internal/installcmd"
	"github.com/Thrasno/conpas-ai/internal/model"
	"github.com/Thrasno/conpas-ai/internal/system"
)

func InstallCommand(profile system.PlatformProfile) ([][]string, error) {
	return installcmd.NewResolver().ResolveComponentInstall(profile, model.ComponentGGA)
}

func ShouldInstall(enabled bool) bool {
	return enabled
}
