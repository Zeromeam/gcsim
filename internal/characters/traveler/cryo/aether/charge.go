package aether

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var chargeFrames []int

func init() {
	chargeFrames = frames.InitAbilSlice(55)
	chargeFrames[action.ActionSkill] = 37
	chargeFrames[action.ActionBurst] = 36
	chargeFrames[action.ActionDash] = 21
	chargeFrames[action.ActionJump] = 21
	chargeFrames[action.ActionSwap] = 44
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	special := c.icepoint >= 3 && !c.StatusIsActive("cryo-traveler-freezing-ice-icd")
	if special {
		c.icepoint = 0
		c.frostglow = min(c.frostglow+2, 8)
		c.AddStatus("cryo-traveler-freezing-ice-icd", 15*60, true)
	}

	for i, mults := range [][]float64{charge1, charge2} {
		ai := info.AttackInfo{
			ActorIndex: c.Index(), Abil: fmt.Sprintf("Charged Attack %d", i+1),
			AttackTag: attacks.AttackTagExtra, ICDTag: attacks.ICDTagNormalAttack,
			ICDGroup: attacks.ICDGroupDefault, StrikeType: attacks.StrikeTypeSlash,
			Element: attributes.Physical, Durability: 25, Mult: mults[c.TalentLvlAttack()],
		}
		if special {
			ai.Abil = fmt.Sprintf("Charged Attack: Freezing Ice %d", i+1)
			ai.Element = attributes.Cryo
			ai.Mult += 1.40
			if c.isStellarSwirl() {
				ai.AttackTag = attacks.AttackTagDirectStellarSwirl
				ai.ICDTag = attacks.ICDTagNone
				ai.Durability = 0
				ai.IgnoreDefPercent = 1
			}
		}
		hitmark := []int{10, 21}[i]
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 2.5), hitmark, hitmark)
	}
	c.ResetNormalCounter()
	return action.Info{Frames: frames.NewAbilFunc(chargeFrames), AnimationLength: chargeFrames[action.InvalidAction], CanQueueAfter: 21, State: action.ChargeAttackState}, nil
}
