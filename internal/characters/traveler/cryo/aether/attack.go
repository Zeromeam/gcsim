package aether

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	attackFrames   [][][]int
	attackHitmarks = [][]int{
		{13, 13, 16, 30, 25}, // Aether
		{16, 10, 19, 23, 14}, // Lumine
	}
	attackHitlag = [][]float64{
		{0.03, 0.03, 0.06, 0.09, 0.12}, // Aether
		{0.03, 0.03, 0.06, 0.06, 0.10}, // Lumine
	}
	attackMults = [][]float64{attack1, attack2, attack3, attack4, attack5}
)

func init() {
	next := [][]int{
		{23, 32, 40, 49, 81}, // Aether
		{24, 21, 27, 38, 64}, // Lumine
	}
	ca := [][]int{
		{34, 34, 44, 55, 500}, // Aether
		{32, 23, 39, 45, 500}, // Lumine
	}
	attackFrames = make([][][]int, 2)
	for gender := range attackFrames {
		attackFrames[gender] = make([][]int, 5)
		for i := range attackFrames[gender] {
			attackFrames[gender][i] = frames.InitNormalCancelSlice(attackHitmarks[gender][i], next[gender][i])
			attackFrames[gender][i][action.ActionAttack] = next[gender][i]
			attackFrames[gender][i][action.ActionCharge] = ca[gender][i]
		}
	}
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(), Abil: fmt.Sprintf("Normal %d", c.NormalCounter+1),
		AttackTag: attacks.AttackTagNormal, ICDTag: attacks.ICDTagNormalAttack,
		ICDGroup: attacks.ICDGroupDefault, StrikeType: attacks.StrikeTypeSlash,
		Element: attributes.Physical, Durability: 25,
		Mult:               attackMults[c.NormalCounter][c.TalentLvlAttack()],
		HitlagFactor:       0.01,
		HitlagHaltFrames:   attackHitlag[c.gender][c.NormalCounter] * 60,
		CanBeDefenseHalted: true,
	}
	hitmark := attackHitmarks[c.gender][c.NormalCounter]
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 2.5), hitmark, hitmark)
	defer c.AdvanceNormalIndex()
	return action.Info{Frames: frames.NewAttackFunc(c.Character, attackFrames[c.gender]), AnimationLength: attackFrames[c.gender][c.NormalCounter][action.InvalidAction], CanQueueAfter: hitmark, State: action.NormalAttackState}, nil
}
