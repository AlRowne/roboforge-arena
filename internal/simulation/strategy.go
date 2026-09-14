package simulation

type Condition int

const (
	DefenderHPAtMost25Percent Condition = iota
	AttackerAPZero
	AttackerEPAtLeast20
)

type Action int

const (
	Attack Action = iota
	Defend
	StrongAttack
)

type Rule struct {
	cond   Condition
	action Action
}

func condActionOK(cond Condition, action Action, attacker, defender BotState) bool {
	var condOK bool
	var actionOK bool

	switch cond {
	case DefenderHPAtMost25Percent:
		condOK = defender.currHP <= defender.maxHP/4
	case AttackerAPZero:
		condOK = attacker.currAP == 0
	case AttackerEPAtLeast20:
		condOK = attacker.currEP >= 20
	default:
		condOK = false
	}

	switch action {
	case Attack:
		actionOK = true
	case Defend:
		actionOK = attacker.currDef == 0
	case StrongAttack:
		actionOK = attacker.currEP >= 20
	default:
		actionOK = false
	}

	return condOK && actionOK
}

func chooseAction(rules []Rule, attacker, defender BotState) (Action, int) {
	for n, rule := range rules {
		if condActionOK(rule.cond, rule.action, attacker, defender) {
			return rule.action, n + 1
		}
	}
	return Attack, 0
}
