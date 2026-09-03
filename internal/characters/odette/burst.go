package odette

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var burstFrames []int

func init() {
	burstFrames = frames.InitAbilSlice(120)
	burstFrames[action.ActionSwap] = 105
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(), AttackTag: attacks.AttackTagElementalBurst,
		ICDTag: attacks.ICDTagElementalBurst, ICDGroup: attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault, Element: attributes.Cryo, Durability: 25,
	}
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 7)
	for i, hitmark := range []int{36, 48, 60} {
		ai.Abil = "Presto: Bluebird Finale Slash"
		ai.Mult = burstSlashDMG[c.TalentLvlBurst()]
		if i > 0 {
			ai.Durability = 0
		}
		c.Core.QueueAttack(ai, ap, 0, hitmark)
	}
	ai.Abil = "Presto: Bluebird Finale Final Slash"
	ai.Mult = burstFinalDMG[c.TalentLvlBurst()]
	ai.Durability = 0
	c.Core.QueueAttack(ai, ap, 0, 72)

	c.AddStatus(snowSwanDreamKey, 20*60, true)
	// Burst summons the Double to Odette's side and refreshes it. Starting a
	// new summon generation invalidates the old scheduled attacks, restarts the
	// measured attack sequence, and grants fresh Marvelous Splendor stacks.
	c.summonDouble()
	c.AddStatus(codaWindowKey, 6*60, true)
	c.ConsumeEnergy(6)
	c.SetCD(action.ActionBurst, 15*60)
	return action.Info{Frames: frames.NewAbilFunc(burstFrames), AnimationLength: burstFrames[action.InvalidAction], CanQueueAfter: burstFrames[action.ActionSwap], State: action.BurstState}, nil
}
