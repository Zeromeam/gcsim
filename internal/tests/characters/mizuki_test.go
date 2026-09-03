package characters

import (
	"math"
	"testing"

	_ "github.com/genshinsim/gcsim/internal/characters/mizuki"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/testhelper"
)

func TestMizukiC1UsesStellarSwirlMultiplier(t *testing.T) {
	c, targets := makeCore(1)
	prof := testhelper.DefaultProfile(keys.YumemizukiMizuki, testhelper.TestWeaponKey)
	prof.Base.Cons = 1
	idx, err := c.AddChar(prof)
	if err != nil {
		t.Fatalf("adding Mizuki: %v", err)
	}
	c.Player.SetActive(idx)
	if err := c.Init(); err != nil {
		t.Fatalf("initializing core: %v", err)
	}
	c.Combat.DefaultTarget = targets[0].Key()

	if err := c.Player.Exec(action.ActionSkill, keys.YumemizukiMizuki, nil); err != nil {
		t.Fatalf("casting Mizuki Skill: %v", err)
	}
	// C1 applies its mark a few frames after Dreamdrifter begins and records
	// Mizuki's EM at that moment.
	for range 5 {
		advanceCoreFrame(c)
	}

	atk := &info.AttackEvent{Info: info.AttackInfo{ActorIndex: idx}}
	c.Events.Emit(event.OnStellarSwirl, targets[0], atk)

	want := 5.5 * c.Player.ByIndex(idx).Stat(attributes.EM)
	if math.Abs(atk.Info.ReactionFlatDmg-want) > 1e-9 {
		t.Fatalf("Mizuki C1 Stellar Swirl additive damage: got %.6f, want %.6f", atk.Info.ReactionFlatDmg, want)
	}
}
