package apis

import (
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/pyrin"
)

func InstallHandlers(app core.App, g pyrin.Group) {
	InstallArtistHandlers(app, g)
	InstallAlbumHandlers(app, g)
	InstallTrackHandlers(app, g)
	InstallQueueHandlers(app, g)
	InstallTagHandlers(app, g)
	InstallAuthHandlers(app, g)
	InstallPlaylistHandlers(app, g)
	InstallSystemHandlers(app, g)
}
