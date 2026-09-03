package combat

import (
	"math"
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func TestStellarSwirlReactionBaseFormula(t *testing.T) {
	atk := info.AttackInfo{
		AttackTag:    attacks.AttackTagReactionStellarSwirl,
		BaseDmgBonus: 0.14,
		Mult:         3,
		FlatDmg:      250,
	}
	got := CalcLunarReactionDmg(90, 0.5, atk, 1000)
	want := 3*(1+6.0*1000/(2000+1000)+0.5)*CalcReactionBaseDmg(90)*1.14 + 250
	if math.Abs(got-want) > 1e-9 {
		t.Fatalf("unexpected Stellar Swirl contribution: got %.12f, want %.12f", got, want)
	}
}
