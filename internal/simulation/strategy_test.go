package simulation

import "testing"

func TestStrategy(t *testing.T) {
	attacker := BotState{
		maxHP:   100,
		maxAP:   120,
		maxEP:   100,
		currHP:  25,
		currAP:  0,
		currEP:  20,
		speed:   100,
		currDef: 0,
	}
	defender := attacker
	action, n := chooseAction([]Rule{
		Rule{
			cond:   DefenderHPAtMost25Percent,
			action: StrongAttack,
		},
		Rule{
			cond:   DefenderHPAtMost25Percent,
			action: StrongAttack,
		},
		Rule{
			cond:   DefenderHPAtMost25Percent,
			action: StrongAttack,
		},
	}, attacker, defender)
	if action != StrongAttack || n != 1 {
		t.Errorf("action should be StrongAttack, is: %v, rule should be 1, is: %v", action, n)

	}

}
