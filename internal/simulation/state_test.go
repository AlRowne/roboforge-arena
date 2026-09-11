package simulation

import (
	"math/rand"
	"slices"
	"testing"
)

func TestCreateBotState(t *testing.T) {
	chaA := Chassis{
		hp:    120,
		ap:    100,
		ep:    100,
		speed: 100,
	}
	modA := Module{
		hp:    0,
		ap:    0,
		speed: 20,
	}
	modC := Module{
		hp:    20,
		ap:    0,
		speed: 0,
	}

	mods := []Module{modA, modC}
	conf := BotConfig{
		chassis: chaA,
		modules: mods,
	}
	modulesBefore := slices.Clone(conf.modules)

	s := createBotstate(conf)

	if s.maxHP != 140 {
		t.Errorf("maxHP not right")
	}
	if s.maxAP != 100 {
		t.Errorf("maxAP not right")
	}
	if s.maxEP != 100 {
		t.Errorf("maxEP not right")
	}
	if s.currHP != 140 {
		t.Errorf("currHP not right")
	}
	if s.currAP != 100 {
		t.Errorf("currHP not right")
	}
	if s.currEP != 50 {
		t.Errorf("currEP not right")
	}
	if s.speed != 120 {
		t.Errorf("speed not right")
	}
	if s.defBuf != false {
		t.Errorf("buf not right")
	}
	if !slices.Equal(modulesBefore, conf.modules) {
		t.Errorf("conf modules were modified")
	}
}
func TestRegenerateEP(t *testing.T) {
	input := BotState{
		maxHP:  140,
		maxAP:  120,
		maxEP:  100,
		currHP: 100,
		currAP: 100,
		currEP: 50,
		speed:  100,
		defBuf: true,
	}
	want := input
	want.currEP = 60
	got, amountReg := regEP(input)
	if got.currEP != 60 || amountReg != 10 {
		t.Errorf("currEP should be 60, is: %v. amountReg should be 10, is: %v", got.currEP, amountReg)
	}
	if got != want {
		t.Errorf("state has been changed. want: %+v\n got: %+v", want, got)
	}
	nearMax := BotState{currEP: 95, maxEP: 100}
	nearMax, amountReg = regEP(nearMax)
	if nearMax.currEP != 100 || amountReg != 5 {
		t.Errorf("currEP should be 100, is: %v. amountReg should be 5, is: %v", nearMax.currEP, amountReg)
	}
	nearMax = BotState{currEP: 115, maxEP: 120}
	nearMax, amountReg = regEP(nearMax)
	if nearMax.currEP != 120 || amountReg != 5 {
		t.Errorf("currEP should be 120, is: %v. amountReg should be 5, is: %v", nearMax.currEP, amountReg)
	}
	nearMax = BotState{currEP: 120, maxEP: 120}
	nearMax, amountReg = regEP(nearMax)
	if nearMax.currEP != 120 || amountReg != 0 {
		t.Errorf("currEP should be 120, is: %v. amountReg should be 0, is: %v", nearMax.currEP, amountReg)
	}
}
func TestCompareSpeed(t *testing.T) {
	speedCompare := compareSpeed(BotState{speed: 120}, BotState{speed: 100})
	if speedCompare != BotAFaster {
		t.Errorf("botA should be faster, got %v", speedCompare)
	}
	speedCompare = compareSpeed(BotState{speed: 100}, BotState{speed: 120})
	if speedCompare != BotBFaster {
		t.Errorf("botB should be faster, got %v", speedCompare)
	}
	speedCompare = compareSpeed(BotState{speed: 100}, BotState{speed: 100})
	if speedCompare != EqualSpeed {
		t.Errorf("speed should be equal, got %v", speedCompare)
	}
}
func TestRandom(t *testing.T) {
	generator1 := rand.New(rand.NewSource(5))
	generator2 := rand.New(rand.NewSource(5))

	for n := range 5 {
		i := rollInitiative(generator1, EqualSpeed)
		j := rollInitiative(generator2, EqualSpeed)
		if i != j {
			t.Errorf("%v. numbers are unequal. gen1: %v, gen2: %v", n, i, j)
		}
	}

	timesBotABegins := 0
	timesBotBBegins := 0

	for _ = range 100 {
		if rollInitiative(generator1, EqualSpeed) == 0 {
			timesBotABegins++
		} else {
			timesBotBBegins++
		}
	}
	if timesBotABegins == 0 || timesBotBBegins == 0 {
		t.Errorf("rollInitiative seems broken. BotA begins %v times, BotB begins %v times", timesBotABegins, timesBotBBegins)
	}
}
