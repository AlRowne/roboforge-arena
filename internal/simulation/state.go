package simulation

import "math/rand"

type Chassis struct {
	hp    int
	ap    int
	ep    int
	speed int
}

type Module struct {
	hp    int
	ap    int
	speed int
}

type BotConfig struct {
	chassis Chassis
	modules []Module
}

type BotState struct {
	maxHP   int
	maxAP   int
	maxEP   int
	currHP  int
	currAP  int
	currEP  int
	currDef int
	speed   int
}

type SpeedComparison int

const (
	EqualSpeed SpeedComparison = iota
	BotAFaster
	BotBFaster
)

type Initiative int

const (
	BotABegins Initiative = iota
	BotBBegins
)

func addMods(m []Module) Module {
	var bonuses Module
	for _, mod := range m {
		bonuses.hp += mod.hp
		bonuses.ap += mod.ap
		bonuses.speed += mod.speed
	}
	return bonuses
}

func createBotstate(conf BotConfig) BotState {
	compMod := addMods(conf.modules)
	state := BotState{

		maxHP:   conf.chassis.hp + compMod.hp,
		maxAP:   conf.chassis.ap + compMod.ap,
		maxEP:   conf.chassis.ep,
		speed:   conf.chassis.speed + compMod.speed,
		currHP:  conf.chassis.hp + compMod.hp,
		currAP:  conf.chassis.ap + compMod.ap,
		currEP:  conf.chassis.ep / 2,
		currDef: 0,
	}
	return state
}

func regEP(state BotState) (BotState, int) {
	gainedEP := 10
	if state.currEP+10 > state.maxEP {
		gainedEP = state.maxEP - state.currEP
		state.currEP = state.maxEP
		return state, gainedEP
	}
	state.currEP += 10
	return state, gainedEP
}

func compareSpeed(botA BotState, botB BotState) SpeedComparison {
	if botA.speed > botB.speed {
		return BotAFaster
	}
	if botB.speed > botA.speed {
		return BotBFaster
	}
	return EqualSpeed
}

func rollInitiative(gen *rand.Rand, sc SpeedComparison) Initiative {
	if sc == EqualSpeed {
		if gen.Intn(2) == 0 {
			return BotABegins
		}
		return BotBBegins
	}
	if sc == BotAFaster {
		return BotABegins
	}
	return BotBBegins
}

func attack(defender *BotState, gen *rand.Rand) {
	dmg := gen.Intn(11) + 20
	dmg = defender.applyDefense(dmg)
	defender.applyDamage(dmg)
}

func strongAttack(attacker *BotState, defender *BotState, gen *rand.Rand) bool {
	if attacker.currEP < 20 {
		return false
	}
	attacker.currEP -= 20
	dmg := gen.Intn(11) + 35
	dmg = defender.applyDefense(dmg)
	defender.applyDamage(dmg)
	return true
}

func (state *BotState) applyDamage(dmg int) {
	dmg = state.currAP - dmg // negative values means HP dmg, 0 and positve values represent new currAP
	if dmg < 0 {
		state.currAP = 0
		state.currHP += dmg
		if state.currHP < 0 {
			state.currHP = 0
		}
	} else {
		state.currAP = dmg
	}
}

func (state *BotState) applyDefense(dmg int) int {
	effDmg := dmg - state.currDef
	state.currDef = 0
	return max(0, effDmg)
}

func (state *BotState) getDefenseBuff(gen *rand.Rand) {
	if state.currDef > 0 {
		return
	}
	def := gen.Intn(11) + 10
	state.currDef = def
}
