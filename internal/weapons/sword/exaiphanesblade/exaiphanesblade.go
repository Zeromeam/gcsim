package exaiphanesblade

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/catalog"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	critDMGPerResonance = 0.06
	defaultResonances   = 7
	procICD             = 5 * 60
	procICDKey          = "exaiphanes-blade-icd"
	attackBuffKey       = "exaiphanes-blade-atk"
	critBuffKey         = "exaiphanes-blade-crit-dmg"
)

var energyRestore = []float64{3, 3, 5, 5, 5}

type Weapon struct{ Index int }

func (w *Weapon) SetIndex(idx int) { w.Index = idx }
func (w *Weapon) Init() error      { return nil }

func NewWeapon(c *core.Core, char *character.CharWrapper, p info.WeaponProfile) (info.Weapon, error) {
	w := &Weapon{}
	// Traveler's Path only functions when the Traveler equips the weapon.
	if catalog.CharacterMap[char.Base.Key].Region != model.AssocType_ASSOC_TYPE_MAINACTOR {
		return w, nil
	}

	// From R2 onward, the blade grants 6% CRIT DMG for every Element the
	// Traveler has resonated with. A Cryo Traveler has necessarily resonated
	// with all seven currently available Elements. The override keeps the
	// weapon usable for deliberately restricted account scenarios.
	resonances := defaultResonances
	if v, ok := p.Params["resonances"]; ok {
		resonances = min(max(v, 0), defaultResonances)
	}
	if p.Refine >= 2 {
		critBuff := make([]float64, attributes.EndStatType)
		critBuff[attributes.CD] = critDMGPerResonance * float64(resonances)
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase(critBuffKey, -1),
			AffectedStat: attributes.CD,
			Amount:       func() []float64 { return critBuff },
		})
	}

	buff := make([]float64, attributes.EndStatType)
	buff[attributes.ATKP] = passiveATK[p.Refine-1]

	c.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != char.Index() || char.StatusIsActive(procICDKey) {
			return
		}
		char.AddStatus(procICDKey, procICD, true)
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag(attackBuffKey, int(passiveDuration*60)),
			AffectedStat: attributes.ATKP,
			Amount:       func() []float64 { return buff },
		})
		char.AddEnergy("exaiphanes-blade", energyRestore[p.Refine-1])
	}, fmt.Sprintf("exaiphanes-blade-%d", char.Index()))

	return w, nil
}
