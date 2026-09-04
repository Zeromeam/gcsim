package characters

import (
	"math"
	"testing"

	_ "github.com/genshinsim/gcsim/internal/characters/traveler/cryo/lumine"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/testhelper"
)

func TestLumineCryoUsesFemaleChargedAttackMultiplier(t *testing.T) {
	c, trg := makeCore(1)
	prof := testhelper.DefaultProfile(keys.LumineCryo, testhelper.TestWeaponKey)
	idx, err := c.AddChar(prof)
	if err != nil {
		t.Fatalf("adding Cryo Lumine: %v", err)
	}
	c.Player.SetActive(idx)
	if err := c.Init(); err != nil {
		t.Fatalf("initializing core: %v", err)
	}
	c.Combat.DefaultTarget = trg[0].Key()

	secondHitMult := 0.0
	c.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.Abil == "Charged Attack 2" {
			secondHitMult = atk.Info.Mult
		}
	}, "lumine-cryo-charge-mult")

	if err := c.Player.Exec(action.ActionCharge, keys.LumineCryo, nil); err != nil {
		t.Fatalf("using Cryo Lumine Charged Attack: %v", err)
	}
	for range 26 {
		advanceCoreFrame(c)
	}
	if math.Abs(secondHitMult-0.7224) > 1e-9 {
		t.Fatalf("Cryo Lumine Charged Attack second multiplier: got %.6f, want 0.7224", secondHitMult)
	}
}

func TestLumineCryoUsesFemaleNormalHitmark(t *testing.T) {
	c, trg := makeCore(1)
	prof := testhelper.DefaultProfile(keys.LumineCryo, testhelper.TestWeaponKey)
	idx, err := c.AddChar(prof)
	if err != nil {
		t.Fatalf("adding Cryo Lumine: %v", err)
	}
	c.Player.SetActive(idx)
	if err := c.Init(); err != nil {
		t.Fatalf("initializing core: %v", err)
	}
	c.Combat.DefaultTarget = trg[0].Key()

	wantFrame := c.F + 16
	gotFrame := -1
	c.Events.Subscribe(event.OnEnemyHit, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.Abil == "Normal 1" {
			gotFrame = c.F
		}
	}, "lumine-cryo-normal-hitmark")

	if err := c.Player.Exec(action.ActionAttack, keys.LumineCryo, nil); err != nil {
		t.Fatalf("using Cryo Lumine Normal Attack: %v", err)
	}
	for range 17 {
		advanceCoreFrame(c)
	}
	if gotFrame != wantFrame {
		t.Fatalf("Cryo Lumine Normal Attack hitmark: got frame %d, want %d", gotFrame, wantFrame)
	}
}
