package characters

import (
	"strings"
	"testing"

	_ "github.com/genshinsim/gcsim/internal/characters/odette"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/testhelper"
)

// Snow Swan's Dream is a temporary Burst buff. Its permanent reaction-bonus
// listener must not make the status appear active before the Burst is cast.
func TestOdetteSnowSwanDreamRequiresBurst(t *testing.T) {
	c, trg := makeCore(1)
	prof := testhelper.DefaultProfile(keys.Odette, testhelper.TestWeaponKey)
	idx, err := c.AddChar(prof)
	if err != nil {
		t.Fatalf("adding Odette: %v", err)
	}
	c.Player.SetActive(idx)
	if err := c.Init(); err != nil {
		t.Fatalf("initializing core: %v", err)
	}
	c.Combat.DefaultTarget = trg[0].Key()

	ai := info.AttackInfo{AttackTag: attacks.AttackTagDirectStellarSwirl}
	char := c.Player.ByIndex(idx)
	if got := char.ReactBonus(ai); got != 0 {
		t.Fatalf("Snow Swan's Dream active before Burst: got reaction bonus %.6f", got)
	}

	if err := c.Player.Exec(action.ActionBurst, keys.Odette, nil); err != nil {
		t.Fatalf("casting Odette Burst: %v", err)
	}
	if got := char.ReactBonus(ai); got <= 0 {
		t.Fatalf("Snow Swan's Dream inactive after Burst: got reaction bonus %.6f", got)
	}
}

func TestOdetteDanceDoubleCompletesTenHitSequence(t *testing.T) {
	c, trg := makeCore(1)
	prof := testhelper.DefaultProfile(keys.Odette, testhelper.TestWeaponKey)
	idx, err := c.AddChar(prof)
	if err != nil {
		t.Fatalf("adding Odette: %v", err)
	}
	c.Player.SetActive(idx)
	if err := c.Init(); err != nil {
		t.Fatalf("initializing core: %v", err)
	}
	c.Combat.DefaultTarget = trg[0].Key()

	hits := 0
	c.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if strings.HasSuffix(atk.Info.Abil, "Dance Move") {
			hits++
		}
	}, "odette-double-hit-count")

	if err := c.Player.Exec(action.ActionSkill, keys.Odette, nil); err != nil {
		t.Fatalf("casting Odette Skill: %v", err)
	}
	for range 22 * 60 {
		advanceCoreFrame(c)
	}
	if hits != 10 {
		t.Fatalf("Dance Double normal hits: got %d, want 10", hits)
	}
}

func TestOdetteCodaStellarDancePersistsAfterRadiance(t *testing.T) {
	c, trg := makeCore(1)
	prof := testhelper.DefaultProfile(keys.Odette, testhelper.TestWeaponKey)
	idx, err := c.AddChar(prof)
	if err != nil {
		t.Fatalf("adding Odette: %v", err)
	}
	c.Player.SetActive(idx)
	if err := c.Init(); err != nil {
		t.Fatalf("initializing core: %v", err)
	}
	c.Combat.DefaultTarget = trg[0].Key()

	// A Stellar Swirl puts Odette in Radiance for eight seconds.
	c.Events.Emit(event.OnStellarSwirl, trg[0], &info.AttackEvent{Info: info.AttackInfo{ActorIndex: idx}})
	if err := c.Player.Exec(action.ActionSkill, keys.Odette, nil); err != nil {
		t.Fatalf("casting Odette Skill: %v", err)
	}
	for range 120 {
		advanceCoreFrame(c)
	}
	if err := c.Player.Exec(action.ActionSkill, keys.Odette, nil); err != nil {
		t.Fatalf("casting Odette Coda: %v", err)
	}

	lateStellarHits := 0
	c.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if c.F > 8*60 && atk.Info.AttackTag == attacks.AttackTagDirectStellarSwirl && strings.Contains(atk.Info.Abil, "Dance Move") {
			lateStellarHits++
		}
	}, "odette-coda-stellar-persistence")

	for range 13 * 60 {
		advanceCoreFrame(c)
	}
	if lateStellarHits == 0 {
		t.Fatal("Dance Double lost its Coda Stellar attacks when Radiance expired")
	}
}
