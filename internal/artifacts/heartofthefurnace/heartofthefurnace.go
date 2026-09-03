package heartofthefurnace

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const buffKey = "heart-of-the-furnace-4pc"

type Set struct {
	Index int
	Count int
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error      { return nil }

// NewSet implements Heart of the Furnace's wearer ATK buff and its party-wide
// Stellar Glimmer reaction bonus. Both effects share the same 12 second status
// and can be refreshed by the wearer from off-field.
func NewSet(c *core.Core, char *character.CharWrapper, count int, _ map[string]int) (info.Set, error) {
	s := &Set{Count: count}

	if count >= 2 {
		buff := make([]float64, attributes.EndStatType)
		buff[attributes.ATKP] = set2ATK
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("heart-of-the-furnace-2pc", -1),
			AffectedStat: attributes.ATKP,
			Amount:       func() []float64 { return buff },
		})
	}

	if count < 4 {
		return s, nil
	}

	teamBuffActive := func() bool {
		for _, holder := range c.Player.Chars() {
			if holder.StatusIsActive(buffKey) {
				return true
			}
		}
		return false
	}
	for _, teammate := range c.Player.Chars() {
		teammate.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBase(buffKey+"-team", -1),
			Amount: func(ai info.AttackInfo) float64 {
				if !teamBuffActive() || !attacks.AttackTagIsStellar(ai.AttackTag) {
					return 0
				}
				return set4StellarDMG
			},
		})
	}

	atkBuff := make([]float64, attributes.EndStatType)
	atkBuff[attributes.ATKP] = set4ATK
	activate := func(atk *info.AttackEvent) {
		if atk.Info.ActorIndex != char.Index() {
			return
		}
		char.AddStatus(buffKey, int(set4Duration*60), true)
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(buffKey+"-atk", int(set4Duration*60)),
			AffectedStat: attributes.ATKP,
			Amount:       func() []float64 { return atkBuff },
		})
	}

	// Triggering a Stellar reaction activates the set even before the reaction
	// damage lands, matching the wording used by other on-reaction sets.
	stellarReaction := func(args ...any) { activate(args[1].(*info.AttackEvent)) }
	c.Events.Subscribe(event.OnStellarConduct, stellarReaction, fmt.Sprintf("%s-conduct-%d", buffKey, char.Index()))
	c.Events.Subscribe(event.OnStellarSwirl, stellarReaction, fmt.Sprintf("%s-swirl-%d", buffKey, char.Index()))

	// Direct Stellar damage can activate the set without triggering a reaction.
	c.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if attacks.AttackTagIsStellar(atk.Info.AttackTag) {
			activate(atk)
		}
	}, fmt.Sprintf("%s-damage-%d", buffKey, char.Index()))

	return s, nil
}
