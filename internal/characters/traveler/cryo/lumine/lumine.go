package lumine

import (
	"github.com/genshinsim/gcsim/internal/characters/traveler/cryo/aether"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func NewChar(s *core.Core, w *character.CharWrapper, p info.CharacterProfile) error {
	return aether.NewLumineChar(s, w, p)
}
