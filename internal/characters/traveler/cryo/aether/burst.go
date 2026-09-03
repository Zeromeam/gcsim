package aether

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

func init() {
	burstFrames = frames.InitAbilSlice(90)
	burstFrames[action.ActionSwap] = 75
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	stacks := c.frostglow
	c.frostglow = 0
	hits := 3
	if stacks == 8 {
		hits = 5
	}

	tag := attacks.AttackTagElementalBurst
	base := burstJavelinDMG[c.TalentLvlBurst()]
	perStack := burstFrostglowBonus[c.TalentLvlBurst()]
	ignoreDef := 0.0
	if c.isStellarSwirl() {
		tag = attacks.AttackTagDirectStellarSwirl
		base = burstSwirlJavelinDMG[c.TalentLvlBurst()]
		perStack = burstSwirlFrostglowBonus[c.TalentLvlBurst()]
		ignoreDef = 1
	}
	ai := info.AttackInfo{
		ActorIndex: c.Index(), AttackTag: tag, ICDTag: attacks.ICDTagElementalBurst,
		ICDGroup: attacks.ICDGroupDefault, StrikeType: attacks.StrikeTypePierce,
		Element: attributes.Cryo, Mult: base + float64(stacks)*perStack,
		IgnoreDefPercent: ignoreDef,
	}
	if tag == attacks.AttackTagElementalBurst {
		ai.Durability = 25
	}
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 4)
	for i := range hits {
		ai.Abil = fmt.Sprintf("Frostbound Javelin %d", i+1)
		c.Core.QueueAttack(ai, ap, 0, 24+i*10)
		ai.Durability = 0
	}

	if c.Base.Cons >= 6 && stacks > 0 {
		for _, teammate := range c.Core.Player.Chars() {
			if teammate.Index() == c.Index() {
				continue
			}
			teammate.AddReactBonusMod(character.ReactBonusMod{
				Base: modifier.NewBaseWithHitlag("cryo-traveler-c6", 15*60),
				Amount: func(ai info.AttackInfo) float64 {
					if attacks.AttackTagIsStellar(ai.AttackTag) {
						return 0.05 * float64(stacks)
					}
					return 0
				},
			})
		}
	}
	c.ConsumeEnergy(6)
	c.SetCD(action.ActionBurst, 15*60)
	return action.Info{Frames: frames.NewAbilFunc(burstFrames), AnimationLength: burstFrames[action.InvalidAction], CanQueueAfter: burstFrames[action.ActionSwap], State: action.BurstState}, nil
}
