package aether

import (
	"github.com/genshinsim/gcsim/internal/characters/traveler/common"
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

type char struct {
	*tmpl.Character
	starSrc   int
	frostglow int
	icepoint  int
}

func NewChar(s *core.Core, w *character.CharWrapper, p info.CharacterProfile) error {
	c := char{starSrc: -1}
	c.Character = tmpl.NewWithWrapper(s, w)
	c.Base.Element = attributes.Cryo
	c.EnergyMax = 60
	c.NormalHitNum = 5
	c.SkillCon = 5
	c.BurstCon = 3
	common.TravelerStoryBuffs(w, p)
	c.addResonanceEnhancements()
	w.Character = &c
	return nil
}

func (c *char) Init() error {
	c.initStellar()
	c.initA4()
	c.initIcepoint()
	c.initConstellations()
	return nil
}

// Foreign Permafrost grants one permanent enhancement for every element the
// Traveler has resonated with. Cryo Traveler necessarily has access to these
// seven live-game bonuses; they are separate from artifact/weapon stats.
func (c *char) addResonanceEnhancements() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.CR] = 0.10
	m[attributes.DEFP] = 0.20
	m[attributes.ER] = 0.20
	m[attributes.EM] = 60
	m[attributes.HPP] = 0.20
	m[attributes.ATKP] = 0.20
	m[attributes.CD] = 0.20
	c.AddStatMod(character.StatMod{
		Base:   modifier.NewBase("cryo-traveler-resonance-enhancements", -1),
		Amount: func() []float64 { return m },
	})
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "frostglow":
		return c.frostglow, nil
	case "icepoint":
		return c.icepoint, nil
	default:
		return c.Character.Condition(fields)
	}
}
