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

var (
	chargeFrames   [][]int
	chargeHitmarks = [][]int{
		{10, 21}, // Aether
		{14, 25}, // Lumine
	}
)

func init() {
	chargeFrames = make([][]int, 2)
	chargeFrames[0] = frames.InitAbilSlice(55)
	chargeFrames[0][action.ActionSkill] = 37
	chargeFrames[0][action.ActionBurst] = 36
	chargeFrames[0][action.ActionDash] = 21
	chargeFrames[0][action.ActionJump] = 21
	chargeFrames[0][action.ActionSwap] = 44

	chargeFrames[1] = frames.InitAbilSlice(58)
	chargeFrames[1][action.ActionSkill] = 34
	chargeFrames[1][action.ActionBurst] = 35
	chargeFrames[1][action.ActionDash] = 25
	chargeFrames[1][action.ActionJump] = 25
	chargeFrames[1][action.ActionSwap] = 25
}

func (c *char) ChargeAttack(p map[string]int) (action.Info, error) {
	special := c.icepoint >= 3 && !c.StatusIsActive("cryo-traveler-freezing-ice-icd")
	if special {
		c.icepoint = 0
		c.frostglow = min(c.frostglow+2, 8)
		c.AddStatus("cryo-traveler-freezing-ice-icd", 15*60, true)
	}

	chargeMults := [][]float64{charge1, charge2}
	if c.gender == 1 {
		chargeMults[1] = lumineCharge2
	}
	for i, mults := range chargeMults {
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
		hitmark := chargeHitmarks[c.gender][i]
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 2.5), hitmark, hitmark)
	}
	c.ResetNormalCounter()
	return action.Info{Frames: frames.NewAbilFunc(chargeFrames[c.gender]), AnimationLength: chargeFrames[c.gender][action.InvalidAction], CanQueueAfter: chargeHitmarks[c.gender][1], State: action.ChargeAttackState}, nil
}
