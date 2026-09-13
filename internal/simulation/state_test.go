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
	if s.currDef != 0 {
		t.Errorf("buf not right")
	}
	if !slices.Equal(modulesBefore, conf.modules) {
		t.Errorf("conf modules were modified")
	}
}
func TestRegenerateEP(t *testing.T) {
	input := BotState{
		maxHP:   140,
		maxAP:   120,
		maxEP:   100,
		currHP:  100,
		currAP:  100,
		currEP:  50,
		speed:   100,
		currDef: 15,
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
		if rollInitiative(generator1, EqualSpeed) == BotABegins {
			timesBotABegins++
		} else {
			timesBotBBegins++
		}
	}
	if timesBotABegins == 0 || timesBotBBegins == 0 {
		t.Errorf("rollInitiative seems broken. BotA begins %v times, BotB begins %v times", timesBotABegins, timesBotBBegins)
	}

	gen3 := rand.New(rand.NewSource(4))
	gen4 := rand.New(rand.NewSource(4))
	for _ = range 5 {
		if in := rollInitiative(gen3, BotAFaster); in != BotABegins {
			t.Errorf("rollInitiavtive got BotAFaster, gave %v", in)
		}
	}
	var got [5]int
	var want [5]int

	for n := range 5 {
		got[n] = gen3.Int()
		want[n] = gen4.Int()
	}
	if got != want {
		t.Errorf("rollInitiative used generator with BotAFaster. got %v, want %v", got, want)
	}

	gen3 = rand.New(rand.NewSource(4))
	gen4 = rand.New(rand.NewSource(4))
	for _ = range 5 {
		if in := rollInitiative(gen3, BotBFaster); in != BotBBegins {
			t.Errorf("rollInitiavtive got BotBFaster, gave %v", in)
		}
	}
	for n := range 5 {
		got[n] = gen3.Int()
		want[n] = gen4.Int()
	}
	if got != want {
		t.Errorf("rollInitiative used generator with BotBFaster. got %v, want %v", got, want)
	}
}

func TestApplyDamage(t *testing.T) {
	type testCase struct {
		name    string
		startAP int
		startHP int
		dmg     int
		wantAP  int
		wantHP  int
	}
	cases := []testCase{
		{name: "shieldOverflow", startAP: 5, startHP: 40, dmg: 20, wantAP: 0, wantHP: 25},
		{name: "ko", startAP: 5, startHP: 40, dmg: 60, wantAP: 0, wantHP: 0},
		{name: "shieldStable", startAP: 30, startHP: 40, dmg: 20, wantAP: 10, wantHP: 40},
		{name: "noDmg", startAP: 30, startHP: 30, dmg: 0, wantAP: 30, wantHP: 30},
		{name: "dmgEqualsAP", startAP: 30, startHP: 30, dmg: 30, wantAP: 0, wantHP: 30},
	}
	for _, c := range cases {
		got := BotState{
			maxHP:   140,
			maxAP:   120,
			maxEP:   100,
			currHP:  c.startHP,
			currAP:  c.startAP,
			currEP:  50,
			speed:   100,
			currDef: 15,
		}
		want := got
		want.currHP = c.wantHP
		want.currAP = c.wantAP
		got.applyDamage(c.dmg)
		if got != want {
			t.Errorf("test '%s': got %+v, want %+v", c.name, got, want)
		}
	}
}
func TestApplyDefense(t *testing.T) {
	type testCase struct {
		name    string
		currDef int
		dmg     int
		effDmg  int
	}

	cases := []testCase{
		{name: "0dmg", dmg: 0, currDef: 10, effDmg: 0},
		{name: "dmg=def", dmg: 10, currDef: 10, effDmg: 0},
		{name: "dmg>def", dmg: 20, currDef: 15, effDmg: 5},
		{name: "dmg<def", dmg: 10, currDef: 15, effDmg: 0},
		{name: "0def", dmg: 20, currDef: 0, effDmg: 20},
	}

	for _, c := range cases {
		got := BotState{
			maxHP:   140,
			maxAP:   120,
			maxEP:   100,
			currHP:  100,
			currAP:  100,
			currEP:  50,
			speed:   100,
			currDef: c.currDef,
		}
		want := got
		want.currDef = 0
		if gotEffDmg := got.applyDefense(c.dmg); gotEffDmg != c.effDmg {
			t.Errorf("case '%v' - got effDmg: %v, want effDmg: %v", c.name, gotEffDmg, c.effDmg)
		}
		if got != want {
			t.Errorf("case '%v' - got: %+v, want: %+v", c.name, got, want)
		}
	}
}

func TestDefendedDamage(t *testing.T) {
	type testCase struct {
		name     string
		startHP  int
		startAP  int
		startDef int
		dmg      int
		wantHP   int
		wantAP   int
	}

	cases := []testCase{
		{name: "dmgSpillover", startHP: 40, startAP: 5, startDef: 15, dmg: 25, wantHP: 35, wantAP: 0},
	}
	for _, c := range cases {
		got := BotState{
			maxHP:   140,
			maxAP:   120,
			maxEP:   100,
			currHP:  c.startHP,
			currAP:  c.startAP,
			currEP:  50,
			speed:   100,
			currDef: c.startDef,
		}
		want := got
		want.currAP = c.wantAP
		want.currDef = 0
		want.currHP = c.wantHP
		got.applyDamage(got.applyDefense(c.dmg))

		if got != want {
			t.Errorf("case '%v' - got: %+v, want: %+v", c.name, got, want)
		}
	}
}

func TestGetDefenseBuff(t *testing.T) {
	got := BotState{
		maxHP:   140,
		maxAP:   120,
		maxEP:   100,
		currHP:  100,
		currAP:  100,
		currEP:  50,
		speed:   100,
		currDef: 0,
	}
	want := got
	got.getDefenseBuff(rand.New(rand.NewSource(69)))
	if got.currDef < 10 || got.currDef > 20 {
		t.Errorf("defBuf is not 10-20. got: %v", got.currDef)
	}
	want.currDef = got.currDef
	if got != want {
		t.Errorf("getDefenseBuff, no Buff - got: %+v, want: %+v", got, want)
	}

	gen1 := rand.New(rand.NewSource(5))
	gen2 := rand.New(rand.NewSource(5))

	got = BotState{
		maxHP:   140,
		maxAP:   120,
		maxEP:   100,
		currHP:  100,
		currAP:  100,
		currEP:  50,
		speed:   100,
		currDef: 15,
	}
	want = got
	got.getDefenseBuff(gen1)
	if got != want {
		t.Errorf("getDefenseBuff, active Buff - got: %+v, want: %+v", got, want)
	}
	var valsGen1 [5]int
	var valsGen2 [5]int

	for n := range 5 {
		valsGen1[n] = gen1.Int()
		valsGen2[n] = gen2.Int()
	}
	if valsGen1 != valsGen2 {
		t.Errorf("random numbers differ. gen1: %v, gen2: %v", valsGen1, valsGen2)
	}
}

func TestAttack(t *testing.T) {
	defenderGot := BotState{
		maxHP:   140,
		maxAP:   120,
		maxEP:   100,
		currHP:  100,
		currAP:  0,
		currEP:  50,
		speed:   100,
		currDef: 0,
	}
	defenderWant := defenderGot
	gen := rand.New(rand.NewSource(69))

	attack(&defenderGot, gen)
	if defenderGot.currHP < 70 || defenderGot.currHP > 80 {
		t.Errorf("defenders hp should be 70 - 80. got: %v", defenderGot.currHP)
	}
	defenderWant.currHP = defenderGot.currHP
	if defenderGot != defenderWant {
		t.Errorf("states changed. got %+v, want %+v", defenderGot, defenderWant)
	}
	defenderGot = BotState{
		maxHP:   140,
		maxAP:   120,
		maxEP:   100,
		currHP:  100,
		currAP:  0,
		currEP:  50,
		speed:   100,
		currDef: 0,
	}
	attackerGot := BotState{
		maxHP:   140,
		maxAP:   120,
		maxEP:   100,
		currHP:  100,
		currAP:  0,
		currEP:  20,
		speed:   100,
		currDef: 0,
	}
	defenderWant = defenderGot
	attackerWant := attackerGot
	attackerWant.currEP = 0

	enoughEP := strongAttack(&attackerGot, &defenderGot, gen)
	if !enoughEP {
		t.Errorf("strongAttack returned false, want true")
	}
	if defenderGot.currHP < 55 || defenderGot.currHP > 65 {
		t.Errorf("strongAttack's dmg range not right. hp: %v, should be 55-65", defenderGot.currHP)
	}
	defenderWant.currHP = defenderGot.currHP
	if attackerGot != attackerWant || defenderGot != defenderWant {
		t.Errorf("botstates changed. attackerGot: %+v, attackerWant: %+v, defenderGot: %+v, defenderWant: %+v",
			attackerGot, attackerWant, defenderGot, defenderWant)
	}
	defenderGot = BotState{
		maxHP:   140,
		maxAP:   120,
		maxEP:   100,
		currHP:  100,
		currAP:  0,
		currEP:  50,
		speed:   100,
		currDef: 0,
	}
	attackerGot = BotState{
		maxHP:   140,
		maxAP:   120,
		maxEP:   100,
		currHP:  100,
		currAP:  0,
		currEP:  19,
		speed:   100,
		currDef: 0,
	}
	defenderWant = defenderGot
	attackerWant = attackerGot

	gen = rand.New(rand.NewSource(69))
	gen2 := rand.New(rand.NewSource(69))

	enoughEP = strongAttack(&attackerGot, &defenderGot, gen)
	if enoughEP {
		t.Errorf("strongAttack with 19EP. got true, want false")
	}
	if attackerGot != attackerWant || defenderGot != defenderWant {
		t.Errorf("attackerGot: %+v, attackerWant: %+v, defenderGot: %+v, defenderWant: %+v",
			attackerGot, attackerWant, defenderGot, defenderWant)
	}
	for _ = range 10 {
		if gen.Int() != gen2.Int() {
			t.Errorf("RNG changed when it shouldn't have")
		}
	}
}
