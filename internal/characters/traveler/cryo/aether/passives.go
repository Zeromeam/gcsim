package aether

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

const (
	radianceKey    = "cryo-traveler-radiance-stellar-swirl"
	icepointICD    = "cryo-traveler-icepoint-icd"
	c1EnergyICD    = "cryo-traveler-c1-energy-icd"
	c2BuffKey      = "cryo-traveler-c2-em"
	c2UpgradeKey   = "cryo-traveler-c2-em-upgraded"
	stellarBaseKey = "cryo-traveler-stellar-base-dmg"
)

func (c *char) initStellar() {
	c.Core.Flags.Custom[reactable.StellarSwirlEnableKey] = 1
	c.Core.Events.Subscribe(event.OnStellarSwirl, func(args ...any) {
		c.AddStatus(radianceKey, 8*60, true)
		c.upgradeC2ForActive(args[1].(*info.AttackEvent))
	}, "cryo-traveler-radiance")

	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		switch atk.Info.AttackTag {
		case attacks.AttackTagDirectStellarConduct, attacks.AttackTagDirectStellarSwirl:
			atk.Info.BaseDmgBonus += c.stellarBaseBonus()
		}
		if attacks.AttackTagIsStellar(atk.Info.AttackTag) {
			c.upgradeC2ForActive(atk)
		}
	}, stellarBaseKey)

	c.Core.Events.Subscribe(event.OnLunarReactionAttack, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.AttackTag == attacks.AttackTagReactionStellarSwirl {
			atk.Info.BaseDmgBonus += c.stellarBaseBonus()
		}
	}, stellarBaseKey+"-reaction")
}

func (c *char) stellarBaseBonus() float64 {
	return min(c.TotalAtk()/100.0*0.0035, 0.07)
}

func (c *char) initA4() {
	if c.Base.Ascension < 4 {
		return
	}
	m := make([]float64, attributes.EndStatType)
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("cryo-traveler-a4", -1),
		AffectedStat: attributes.EM,
		Extra:        true,
		Amount: func() []float64 {
			atk := c.SelectStat(true, attributes.BaseATK, attributes.ATKP, attributes.ATK).TotalATK()
			m[attributes.EM] = min(atk*0.08, 160)
			return m
		},
	})
}

func (c *char) initIcepoint() {
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if !attacks.AttackTagIsStellar(atk.Info.AttackTag) || c.StatusIsActive(icepointICD) {
			return
		}
		c.icepoint = min(c.icepoint+1, 3)
		c.AddStatus(icepointICD, 2*60, true)
	}, "cryo-traveler-icepoint")
}

func (c *char) initConstellations() {
	if c.Base.Cons < 1 {
		return
	}
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != c.Index() || !attacks.AttackTagIsStellar(atk.Info.AttackTag) || c.StatusIsActive(c1EnergyICD) {
			return
		}
		c.AddStatus(c1EnergyICD, 30, true)
		c.AddEnergy("cryo-traveler-c1", 5)
	}, "cryo-traveler-c1")
}

func (c *char) applyC2() {
	if c.Base.Cons < 2 {
		return
	}
	target := c.Core.Player.ActiveChar()
	m := make([]float64, attributes.EndStatType)
	target.AddStatus(c2BuffKey, 5*60, true)
	target.DeleteStatus(c2UpgradeKey)
	target.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag(c2BuffKey, 5*60),
		AffectedStat: attributes.EM,
		Amount: func() []float64 {
			m[attributes.EM] = 60
			if target.StatusIsActive(c2UpgradeKey) {
				m[attributes.EM] = 120
			}
			return m
		},
	})
}

func (c *char) upgradeC2ForActive(atk *info.AttackEvent) {
	if c.Base.Cons < 2 || atk.Info.ActorIndex != c.Core.Player.Active() {
		return
	}
	target := c.Core.Player.ActiveChar()
	if !target.StatusIsActive(c2BuffKey) {
		return
	}
	target.AddStatus(c2UpgradeKey, target.StatusDuration(c2BuffKey), true)
}

func (c *char) isStellarSwirl() bool {
	return c.StatusIsActive(radianceKey)
}
