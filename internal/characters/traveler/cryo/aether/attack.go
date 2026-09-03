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
	attackFrames   [][]int
	attackHitmarks = []int{13, 13, 16, 30, 25}
	attackMults    = [][]float64{attack1, attack2, attack3, attack4, attack5}
)

func init() {
	next := []int{23, 32, 40, 49, 81}
	ca := []int{34, 34, 44, 55, 500}
	attackFrames = make([][]int, 5)
	for i := range attackFrames {
		attackFrames[i] = frames.InitNormalCancelSlice(attackHitmarks[i], next[i])
		attackFrames[i][action.ActionAttack] = next[i]
		attackFrames[i][action.ActionCharge] = ca[i]
	}
}

func (c *char) Attack(p map[string]int) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(), Abil: fmt.Sprintf("Normal %d", c.NormalCounter+1),
		AttackTag: attacks.AttackTagNormal, ICDTag: attacks.ICDTagNormalAttack,
		ICDGroup: attacks.ICDGroupDefault, StrikeType: attacks.StrikeTypeSlash,
		Element: attributes.Physical, Durability: 25,
		Mult: attackMults[c.NormalCounter][c.TalentLvlAttack()],
	}
	hitmark := attackHitmarks[c.NormalCounter]
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 2.5), hitmark, hitmark)
	defer c.AdvanceNormalIndex()
	return action.Info{Frames: frames.NewAttackFunc(c.Character, attackFrames), AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction], CanQueueAfter: hitmark, State: action.NormalAttackState}, nil
}
