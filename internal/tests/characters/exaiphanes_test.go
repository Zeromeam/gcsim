package characters

import (
	"math"
	"testing"

	_ "github.com/genshinsim/gcsim/internal/characters/traveler/anemo/aether"
	_ "github.com/genshinsim/gcsim/internal/weapons/sword/exaiphanesblade"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/testhelper"
)

func addTravelerWithExaiphanes(t *testing.T, refine int) (float64, float64) {
	t.Helper()
	c, targets := makeCore(1)
	prof := testhelper.DefaultProfile(keys.AetherAnemo, keys.ExaiphanesBlade)
	prof.Weapon.Level = 90
	prof.Weapon.MaxLevel = 90
	prof.Weapon.Refine = refine
	idx, err := c.AddChar(prof)
	if err != nil {
		t.Fatalf("adding Traveler: %v", err)
	}
	c.Player.SetActive(idx)
	if err := c.Init(); err != nil {
		t.Fatalf("initializing core: %v", err)
	}

	char := c.Player.ByIndex(idx)
	cd := char.Stat(attributes.CD)
	char.Energy = 0
	c.Events.Emit(event.OnEnemyDamage, targets[0], &info.AttackEvent{Info: info.AttackInfo{ActorIndex: idx}}, 1.0, false)
	return cd, char.Energy
}

func TestExaiphanesR3TravelerBonuses(t *testing.T) {
	r1CD, r1Energy := addTravelerWithExaiphanes(t, 1)
	r3CD, r3Energy := addTravelerWithExaiphanes(t, 3)

	if math.Abs((r3CD-r1CD)-0.42) > 1e-9 {
		t.Fatalf("R3 CRIT DMG over R1: got %.6f, want 0.42", r3CD-r1CD)
	}
	if r1Energy != 3 || r3Energy != 5 {
		t.Fatalf("energy restoration: got R1 %.1f and R3 %.1f, want 3 and 5", r1Energy, r3Energy)
	}
}
