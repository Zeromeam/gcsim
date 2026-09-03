package odette

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

const (
	radianceKey         = "odette-radiance-stellar-swirl"
	doubleKey           = "odette-solo-dance-double"
	codaWindowKey       = "odette-coda-window"
	snowSwanDreamKey    = "odette-snow-swan-dream"
	snowSwanDreamModKey = "odette-snow-swan-dream-react-bonus"
	stellarBaseKey      = "odette-stellar-base-dmg"
	stellarElevateKey   = "odette-stellar-elevation"
	particleICDKey      = "odette-particle-icd"
	particleICD         = 12 * 60
)

type char struct {
	*tmpl.Character
	doubleSrc     int
	doubleStellar bool
	selfSplendor  int
	teamSplendor  int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{doubleSrc: -1}
	c.Character = tmpl.NewWithWrapper(s, w)
	c.EnergyMax = 60
	c.NormalHitNum = 5
	c.SkillCon = 3
	c.BurstCon = 5
	w.Character = &c
	return nil
}

func (c *char) Init() error {
	c.Core.Flags.Custom[reactable.StellarSwirlEnableKey] = 1
	c.Core.Events.Subscribe(event.OnStellarSwirl, func(args ...any) {
		c.AddStatus(radianceKey, 8*60, true)
	}, "odette-radiance")
	c.initStellarModifiers()
	c.initSplendor()
	return nil
}

func (c *char) ActionReady(a action.Action, p map[string]int) (bool, action.Failure) {
	if a == action.ActionSkill && c.StatusIsActive(codaWindowKey) {
		return true, action.NoFailure
	}
	return c.Character.ActionReady(a, p)
}

func (c *char) stellarBaseBonus() float64 {
	return min(c.TotalAtk()/100.0*0.007, 0.14)
}

func (c *char) stellarElevation() float64 {
	return min(max(c.TotalAtk()-1000, 0)*0.00015, 0.30)
}

func (c *char) initStellarModifiers() {
	applyDirect := func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		switch atk.Info.AttackTag {
		case attacks.AttackTagDirectStellarConduct, attacks.AttackTagDirectStellarSwirl:
		default:
			return
		}
		atk.Info.BaseDmgBonus += c.stellarBaseBonus()
		if atk.Info.ActorIndex == c.Index() && c.Base.Ascension >= 4 {
			atk.Info.Elevation += c.stellarElevation()
		}
	}
	c.Core.Events.Subscribe(event.OnEnemyHit, applyDirect, stellarBaseKey)

	c.Core.Events.Subscribe(event.OnLunarReactionAttack, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.AttackTag != attacks.AttackTagReactionStellarSwirl {
			return
		}
		atk.Info.BaseDmgBonus += c.stellarBaseBonus()
		if atk.Info.ActorIndex == c.Index() && c.Base.Ascension >= 4 {
			atk.Info.Elevation += c.stellarElevation()
		}
	}, stellarElevateKey)

	c.AddReactBonusMod(character.ReactBonusMod{
		// Keep the permanent listener's key distinct from the temporary Burst
		// status. StatusIsActive searches all modifiers by key, so sharing the
		// key would make Snow Swan's Dream appear active before Burst was cast.
		Base: modifier.NewBase(snowSwanDreamModKey, -1),
		Amount: func(ai info.AttackInfo) float64 {
			if !c.StatusIsActive(snowSwanDreamKey) || !attacks.AttackTagIsStellar(ai.AttackTag) {
				return 0
			}
			return burstStellarBonus[c.TalentLvlBurst()]
		},
	})
}

func (c *char) initSplendor() {
	for _, teammate := range c.Core.Player.Chars() {
		teammate := teammate
		teammate.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBase("odette-marvelous-splendor", -1),
			Amount: func(ai info.AttackInfo) float64 {
				if !c.StatusIsActive(doubleKey) || !attacks.AttackTagIsStellar(ai.AttackTag) {
					return 0
				}
				if teammate.Index() == c.Index() {
					return 0.15 * float64(c.selfSplendor)
				}
				return 0.15 * float64(c.teamSplendor)
			},
		})
	}
}

func (c *char) isStellarSwirl() bool {
	return c.StatusIsActive(radianceKey)
}

// Odette generates five Cryo particles when either version of her Skill hits,
// with one shared 12 second generation cooldown.
func (c *char) particleCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy || c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, particleICD, true)
	c.Core.QueueParticle(c.Base.Key.String(), 5, c.Base.Element, c.ParticleDelay)
}
