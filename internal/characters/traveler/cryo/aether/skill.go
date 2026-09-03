package aether

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var skillFrames []int

func init() {
	skillFrames = frames.InitAbilSlice(45)
	skillFrames[action.ActionSwap] = 35
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(), Abil: "Ice Fog Piercer",
		AttackTag: attacks.AttackTagElementalArt, ICDTag: attacks.ICDTagElementalArt,
		ICDGroup: attacks.ICDGroupDefault, StrikeType: attacks.StrikeTypePierce,
		Element: attributes.Cryo, Durability: 25, Mult: skillDMG[c.TalentLvlSkill()],
	}
	particles := false
	particleCB := func(a info.AttackCB) {
		if particles || a.Target.Type() != info.TargettableEnemy {
			return
		}
		particles = true
		c.Core.QueueParticle(c.Base.Key.String(), 3, attributes.Cryo, c.ParticleDelay)
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 4), 12, 12, particleCB)
	c.startFrostpierceStar()
	c.SetCD(action.ActionSkill, 15*60)
	return action.Info{Frames: frames.NewAbilFunc(skillFrames), AnimationLength: skillFrames[action.InvalidAction], CanQueueAfter: skillFrames[action.ActionSwap], State: action.SkillState}, nil
}

func (c *char) startFrostpierceStar() {
	c.starSrc = c.Core.F
	src := c.starSrc
	duration := 12 * 60
	if c.Base.Cons >= 4 {
		duration = 15 * 60
	}
	for hitmark := 2 * 60; hitmark <= duration; hitmark += 2 * 60 {
		c.QueueCharTask(func() {
			if c.starSrc != src {
				return
			}
			c.fireIceCrystal()
		}, hitmark)
	}
	c.QueueCharTask(func() {
		if c.starSrc == src {
			c.starSrc = -1
		}
	}, duration+1)
}

func (c *char) fireIceCrystal() {
	done := false
	cb := func(a info.AttackCB) {
		if done || a.Target.Type() != info.TargettableEnemy {
			return
		}
		done = true
		c.frostglow = min(c.frostglow+1, 8)
		c.applyC2()
	}
	ai := info.AttackInfo{
		ActorIndex: c.Index(), Abil: "Frostpierce Star Ice Crystal",
		AttackTag: attacks.AttackTagElementalArt, ICDTag: attacks.ICDTagElementalArt,
		ICDGroup: attacks.ICDGroupDefault, StrikeType: attacks.StrikeTypeDefault,
		Element: attributes.Cryo, Durability: 25, Mult: crystalDMG[c.TalentLvlSkill()], IsDeployable: true,
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 2.5), 0, 0, cb)
}
