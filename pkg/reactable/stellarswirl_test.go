package reactable

import (
	"math"
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func TestStellarSwirlStackProgression(t *testing.T) {
	c := testCore()
	target := addTargetToCore(c)

	for want := 1; want <= sswMaxStacks; want++ {
		target.addSSwStack()
		if got := target.sswStacks(); got != want {
			t.Fatalf("after trigger %d: got %d stacks, want %d", want, got, want)
		}
	}

	target.addSSwStack()
	if got := target.sswStacks(); got != sswMaxStacks {
		t.Fatalf("stack count exceeded cap: got %d, want %d", got, sswMaxStacks)
	}
}

func TestStellarSwirlCarriesAdditiveDmgIntoEveryContribution(t *testing.T) {
	c, targets := testCoreWithTrgs(1)
	if err := c.Init(); err != nil {
		t.Fatalf("initializing core: %v", err)
	}

	const additive = 1234.5
	seen := 0
	individualDmg := 0.0
	c.Events.Subscribe(event.OnLunarReactionAttack, func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		if atk.Info.AttackTag != attacks.AttackTagReactionStellarSwirl {
			return
		}
		seen++
		if got := atk.Info.FlatDmg; got != additive {
			t.Fatalf("individual contribution additive damage: got %.3f, want %.3f", got, additive)
		}
		// Eliminate the random crit branch and calculate the same individual
		// contribution independently for the final combination assertion.
		atk.Snapshot.Stats[attributes.CR] = 0
		char := c.Player.ByIndex(atk.Info.ActorIndex)
		individualDmg = combat.CalcLunarReactionDmg(
			char.Base.Level,
			char.ReactBonus(atk.Info),
			atk.Info,
			atk.Snapshot.Stats[attributes.EM],
		)
	}, "stellar-swirl-additive-test")

	var contributors [info.MaxChars]bool
	contributors[0] = true
	ai := info.AttackInfo{
		ActorIndex:       0,
		AttackTag:        attacks.AttackTagReactionStellarSwirl,
		FlatDmg:          additive,
		IgnoreDefPercent: 1,
	}
	ap := combat.NewSingleTargetHit(targets[0].Key())
	combined, _ := targets[0].calcStellarSwirlDmg(targets[0], ai, ap, contributors, 0.75)

	if seen != 1 {
		t.Fatalf("contribution events: got %d, want 1", seen)
	}
	want := individualDmg * sswContributorMult[0]
	if math.Abs(combined.FlatDmg-want) > 1e-9 {
		t.Fatalf("combined damage reapplied additive term: got %.6f, want %.6f", combined.FlatDmg, want)
	}
}
