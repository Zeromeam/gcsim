package characters

import (
	"math"
	"strings"
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

func TestMizukiWitchRevelationParamControlsC1AdditionalAttack(t *testing.T) {
	tests := []struct {
		name       string
		params     map[string]int
		wantDirect int
	}{
		{name: "default-enabled", params: map[string]int{}, wantDirect: 1},
		{name: "explicitly-disabled", params: map[string]int{"witch_revelation": 0}, wantDirect: 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, targets := makeCore(1)
			prof := testhelper.DefaultProfile(keys.YumemizukiMizuki, testhelper.TestWeaponKey)
			prof.Base.Cons = 1
			for key, value := range tc.params {
				prof.Params[key] = value
			}
			idx, err := c.AddChar(prof)
			if err != nil {
				t.Fatalf("adding Mizuki: %v", err)
			}
			c.Player.SetActive(idx)
			if err := c.Init(); err != nil {
				t.Fatalf("initializing core: %v", err)
			}
			c.Combat.DefaultTarget = targets[0].Key()

			directHits := 0
			c.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
				atk := args[1].(*info.AttackEvent)
				if strings.HasPrefix(atk.Info.Abil, "Twenty-Three Nights' Awaiting") {
					directHits++
				}
			}, "mizuki-witch-param-direct-hit")

			if err := c.Player.Exec(action.ActionSkill, keys.YumemizukiMizuki, nil); err != nil {
				t.Fatalf("casting Mizuki Skill: %v", err)
			}
			for range 5 {
				advanceCoreFrame(c)
			}
			atk := &info.AttackEvent{Info: info.AttackInfo{ActorIndex: idx}}
			c.Events.Emit(event.OnStellarSwirl, targets[0], atk)
			advanceCoreFrame(c)

			if directHits != tc.wantDirect {
				t.Fatalf("Witch C1 direct hits: got %d, want %d", directHits, tc.wantDirect)
			}
			wantAdditive := 5.5 * c.Player.ByIndex(idx).Stat(attributes.EM)
			if math.Abs(atk.Info.ReactionFlatDmg-wantAdditive) > 1e-9 {
				t.Fatalf("base C1 additive damage changed: got %.6f, want %.6f", atk.Info.ReactionFlatDmg, wantAdditive)
			}
		})
	}
}

func TestMizukiWitchRevelationParamControlsC2ResShred(t *testing.T) {
	tests := []struct {
		name      string
		params    map[string]int
		wantShred bool
	}{
		{name: "default-enabled", params: map[string]int{}, wantShred: true},
		{name: "explicitly-disabled", params: map[string]int{"witch_revelation": 0}, wantShred: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, targets := makeCore(1)
			prof := testhelper.DefaultProfile(keys.YumemizukiMizuki, testhelper.TestWeaponKey)
			prof.Base.Cons = 2
			for key, value := range tc.params {
				prof.Params[key] = value
			}
			idx, err := c.AddChar(prof)
			if err != nil {
				t.Fatalf("adding Mizuki: %v", err)
			}
			c.Player.SetActive(idx)
			if err := c.Init(); err != nil {
				t.Fatalf("initializing core: %v", err)
			}

			if err := c.Player.Exec(action.ActionSkill, keys.YumemizukiMizuki, nil); err != nil {
				t.Fatalf("casting Mizuki Skill: %v", err)
			}
			c.Events.Emit(event.OnEnemyHit, targets[0], &info.AttackEvent{Info: info.AttackInfo{ActorIndex: idx}})

			got := targets[0].ResistModIsActive("mizuki-c2-cryo")
			if got != tc.wantShred {
				t.Fatalf("Witch C2 Cryo RES shred active: got %v, want %v", got, tc.wantShred)
			}
		})
	}
}
