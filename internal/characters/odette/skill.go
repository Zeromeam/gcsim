package odette

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	skillFrames []int
	codaFrames  []int
)

const (
	skillHitmark   = 42
	doubleDuration = 20 * 60
	// The summon appears roughly one second into the cast. Its measured attack
	// sequence is timed from action start, so keep the runtime status alive for
	// that appearance delay as well as the listed 20 second field duration.
	doubleLifetime = doubleDuration + 60
	// The special Skill can be queued as the initial Skill hit lands. Its
	// three duet hits occur at roughly 1/3-second intervals, followed by the
	// final Coda at 80f. These timings are visible in the 7.0 dummy showcase:
	// initial hit at 2.31s, duet damage through 3.30s, final hit at 3.63s.
	codaDuration = 80
)

var codaDotHitmarks = []int{20, 40, 60}

// Measured from the summon action start. The live 7.0 sequence alternates
// Plume/Wing and ends just before the 20 second summon expires.
var doubleHitmarks = []int{168, 282, 408, 516, 642, 756, 885, 996, 1122, 1236}

func init() {
	skillFrames = frames.InitAbilSlice(120)
	// Odette can immediately enter Coda when the summoning hit lands.
	skillFrames[action.ActionSkill] = skillHitmark
	skillFrames[action.ActionSwap] = 120
	codaFrames = frames.InitAbilSlice(codaDuration)
	codaFrames[action.ActionSwap] = codaDuration
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(codaWindowKey) || p["second"] != 0 {
		return c.coda()
	}

	ai := info.AttackInfo{
		ActorIndex: c.Index(), Abil: "Adagio: Phantom Night Dancers",
		AttackTag: attacks.AttackTagElementalArt, ICDTag: attacks.ICDTagElementalArt,
		ICDGroup: attacks.ICDGroupDefault, StrikeType: attacks.StrikeTypeDefault,
		Element: attributes.Cryo, Durability: 25, Mult: skillDMG[c.TalentLvlSkill()],
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5), 0, skillHitmark, c.particleCB)
	c.summonDouble()
	c.AddStatus(codaWindowKey, 6*60, true)
	c.SetCD(action.ActionSkill, 15*60)
	return action.Info{Frames: frames.NewAbilFunc(skillFrames), AnimationLength: skillFrames[action.InvalidAction], CanQueueAfter: skillHitmark, State: action.SkillState}, nil
}

func (c *char) coda() (action.Info, error) {
	c.DeleteStatus(codaWindowKey)
	c.doubleStellar = c.StatusIsActive(doubleKey) && c.isStellarSwirl()

	dot := info.AttackInfo{
		ActorIndex: c.Index(), Abil: "Coda at Dawn's Tolling DoT",
		AttackTag: attacks.AttackTagElementalArt, ICDTag: attacks.ICDTagElementalArt,
		ICDGroup: attacks.ICDGroupDefault, StrikeType: attacks.StrikeTypeDefault,
		Element: attributes.Cryo, Durability: 10, Mult: codaDotDMG[c.TalentLvlSkill()],
	}
	for _, hitmark := range codaDotHitmarks {
		c.Core.QueueAttack(dot, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5), hitmark, hitmark, c.particleCB)
	}

	c.QueueCharTask(func() {
		tag := attacks.AttackTagDirectStellarConduct
		mult := codaConductDMG[c.TalentLvlSkill()]
		if c.isStellarSwirl() {
			tag = attacks.AttackTagDirectStellarSwirl
			mult = codaSwirlDMG[c.TalentLvlSkill()]
		}
		ai := info.AttackInfo{
			ActorIndex: c.Index(), Abil: "Coda at Dawn's Tolling",
			AttackTag: tag, ICDTag: attacks.ICDTagNone, ICDGroup: attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault, Element: attributes.Cryo,
			Mult: mult, IgnoreDefPercent: 1,
		}
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 6), 0, 0)
	}, codaDuration)

	return action.Info{Frames: frames.NewAbilFunc(codaFrames), AnimationLength: codaFrames[action.InvalidAction], CanQueueAfter: codaFrames[action.ActionSwap], State: action.SkillState}, nil
}

func (c *char) summonDouble() {
	c.doubleSrc = c.Core.F
	src := c.doubleSrc
	c.doubleStellar = false
	c.AddStatus(doubleKey, doubleLifetime, true)
	c.selfSplendor = 0
	c.teamSplendor = 0
	if c.Base.Ascension >= 1 {
		c.selfSplendor = 4
	}

	for i, hitmark := range doubleHitmarks {
		plume := i%2 == 0
		c.QueueCharTask(func() {
			if c.doubleSrc != src || !c.StatusIsActive(doubleKey) {
				return
			}
			c.doubleAttack(plume)
		}, hitmark)
	}

	for hitmark := 60; hitmark <= doubleDuration; hitmark += 60 {
		c.QueueCharTask(func() {
			if c.doubleSrc != src || !c.StatusIsActive(doubleKey) || c.Core.Player.Active() == c.Index() || c.selfSplendor == 0 {
				return
			}
			c.selfSplendor--
			c.teamSplendor = min(c.teamSplendor+1, 4)
		}, hitmark)
	}
	c.QueueCharTask(func() {
		if c.doubleSrc != src {
			return
		}
		c.doubleSrc = -1
		c.selfSplendor = 0
		c.teamSplendor = 0
		c.DeleteStatus(doubleKey)
	}, doubleLifetime+1)
}

func (c *char) doubleAttack(plume bool) {
	name := "Plume Dance Move"
	mult := plumeDMG[c.TalentLvlSkill()]
	directMult := plumeSwirlDMG[c.TalentLvlSkill()]
	if !plume {
		name = "Wing Dance Move"
		mult = wingDMG[c.TalentLvlSkill()]
		directMult = wingSwirlDMG[c.TalentLvlSkill()]
	}
	ai := info.AttackInfo{
		ActorIndex: c.Index(), Abil: name, AttackTag: attacks.AttackTagElementalArt,
		ICDTag: attacks.ICDTagNone, ICDGroup: attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault, Element: attributes.Cryo,
		Durability: 25, Mult: mult, IsDeployable: true,
	}
	// Each periodic Dance Double attack applies Cryo. Sharing standard Skill
	// ICD with the summon and Coda under-applies Cryo and breaks the continuous
	// off-field application visible in the live rotation.
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 5)
	c.Core.QueueAttack(ai, ap, 0, 0)
	// Coda snapshots whether the Dance Double was converted while Odette was
	// in Radiance. The additional Stellar attack then remains part of this
	// summon's dance moves; it does not require Radiance to remain active.
	if c.doubleStellar {
		ai.Abil = name + " (Stellar Swirl)"
		ai.AttackTag = attacks.AttackTagDirectStellarSwirl
		ai.ICDTag = attacks.ICDTagNone
		ai.Durability = 0
		ai.Mult = directMult
		ai.IgnoreDefPercent = 1
		c.Core.QueueAttack(ai, ap, 0, 0)
	}
}
