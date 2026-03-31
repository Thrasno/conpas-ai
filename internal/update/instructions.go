package update

import (
	"github.com/Thrasno/conpas-ai/internal/system"
)

// updateHint returns a platform-specific instruction string for updating the given tool.
func updateHint(tool ToolInfo, profile system.PlatformProfile) string {
	switch tool.Name {
	case "conpas-ai":
		return gentleAIHint(profile)
	case "engram":
		return engramHint(profile)
	case "gga":
		return ggaHint(profile)
	default:
		return ""
	}
}

func gentleAIHint(profile system.PlatformProfile) string {
	switch profile.OS {
	case "darwin":
		return "brew upgrade conpas-ai"
	case "linux":
		return "curl -fsSL https://raw.githubusercontent.com/Thrasno/conpas-ai/main/scripts/install.sh | bash"
	case "windows":
		return "irm https://raw.githubusercontent.com/Thrasno/conpas-ai/main/scripts/install.ps1 | iex"
	default:
		return ""
	}
}

func engramHint(profile system.PlatformProfile) string {
	switch profile.PackageManager {
	case "brew":
		return "brew upgrade engram"
	default:
		return "conpas-ai upgrade (downloads pre-built binary)"
	}
}

func ggaHint(profile system.PlatformProfile) string {
	switch profile.PackageManager {
	case "brew":
		return "brew upgrade gga"
	default:
		return "See https://github.com/Gentleman-Programming/gentleman-guardian-angel"
	}
}
