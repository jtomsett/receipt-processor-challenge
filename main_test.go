package main

import (
	"math"
	"testing"
)

func TestRetailerNameCalculationSimple(t *testing.T) {
	name := "basic123"
	pts := calculatePointsForRetailerName(name)
	if pts != len(name) {
		t.Fatalf(`%s should be %d points but was %d points!`, name, len(name), pts)
	}
}

func TestRetailerNameCalculationEmpty(t *testing.T) {
	name := ""
	pts := calculatePointsForRetailerName(name)
	if pts != len(name) {
		t.Fatalf(`%s should be %d points but was %d points!`, name, len(name), pts)
	}
}

func TestRetailerNameCalculationTrailingNonAN(t *testing.T) {
	name := "basic123 %"
	pts := calculatePointsForRetailerName(name)
	expected := len(name) - 2
	if pts != expected {
		t.Fatalf(`%s should be %d points but was %d points!`, name, expected, pts)
	}
}

func TestRetailerNameCalculationStartingNonAn(t *testing.T) {
	name := "& basic123"
	pts := calculatePointsForRetailerName(name)
	expected := len(name) - 2
	if pts != expected {
		t.Fatalf(`%s should be %d points but was %d points!`, name, expected, pts)
	}
}

func TestRetailerNameCalculationInternalNonAn(t *testing.T) {
	name := "basic_ &123"
	pts := calculatePointsForRetailerName(name)
	expected := len(name) - 3
	if pts != expected {
		t.Fatalf(`%s should be %d points but was %d points!`, name, expected, pts)
	}
}

func TestRoundTotalBonusRound(t *testing.T) {
	total := "12.00"
	pts, err := calculateRoundTotalEvenBonus(total)
	if err != nil && pts != 50 {
		t.Fatalf(`%s should have been a 50 point bonus!`, total)
	}
}

func TestRoundTotalBonusNoBonus(t *testing.T) {
	total := "12.12"
	pts, err := calculateRoundTotalEvenBonus(total)
	if err != nil && pts != 0 {
		t.Fatalf(`%s should not have been a 50 point bonus!`, total)
	}
}

func TestRoundTotalBonusError(t *testing.T) {
	total := "error"
	pts, err := calculateRoundTotalEvenBonus(total)
	if err == nil && pts != 0 {
		t.Fatalf(`%s should have thrown an error!`, total)
	}
}

func TestRoundTotalBonus25(t *testing.T) {
	total := "12.75"
	pts, err := calculateRoundTotal25Bonus(total)
	if err != nil && pts != 25 {
		t.Fatalf(`%s should have been a 25 point bonus!`, total)
	}
}

func TestRoundTotalBonusNo25Bonus(t *testing.T) {
	total := "12.12"
	pts, err := calculateRoundTotal25Bonus(total)
	if err != nil && pts != 0 {
		t.Fatalf(`%s should not have been a 25 point bonus!`, total)
	}
}

func TestRoundTotalBonus25Error(t *testing.T) {
	total := "error"
	pts, err := calculateRoundTotal25Bonus(total)
	if err == nil && pts != 0 {
		t.Fatalf(`%s should have thrown an error!`, total)
	}
}

func TestItemLengthBonus(t *testing.T) {
	var items = []item{
		{
			ShortDescription: "asdf", Price: "1",
		},
		{
			ShortDescription: "asdf", Price: "1",
		},
		{
			ShortDescription: "asdf", Price: "1",
		},
		{
			ShortDescription: "asdf", Price: "1",
		},
	}
	pts := calculateItemLengthBonus(items)
	if pts != 10 {
		t.Fatalf(`%d item length should be %d points!`, len(items), 2)
	}
}

func TestItemLengthBonusNil(t *testing.T) {
	pts := calculateItemLengthBonus(nil)
	if pts != 0 {
		t.Fatalf(`nil items should be %d points!`, 0)
	}
}

func TestItemLengthBonusSingle(t *testing.T) {
	var items = []item{
		{
			ShortDescription: "asdf", Price: "1",
		},
	}
	pts := calculateItemLengthBonus(items)
	if pts != 0 {
		t.Fatalf(`%d item length should be %d points!`, len(items), 0)
	}
}

func TestItemLengthBonusOdd(t *testing.T) {
	var items = []item{
		{
			ShortDescription: "asdf", Price: "1",
		},
		{
			ShortDescription: "asdf", Price: "1",
		},
		{
			ShortDescription: "asdf", Price: "1",
		},
	}
	pts := calculateItemLengthBonus(items)
	if pts != 5 {
		t.Fatalf(`%d item length should be %d points!`, len(items), 5)
	}
}

func TestItemLengthBonusEmpty(t *testing.T) {
	var items = []item{}
	pts := calculateItemLengthBonus(items)
	if pts != 0 {
		t.Fatalf(`%d item length should be %d points!`, len(items), 0)
	}
}

func TestItemDescriptionBonus(t *testing.T) {
	var items = []item{
		{
			ShortDescription: "asdf4_", Price: "1",
		},
		{
			ShortDescription: " asdf4_   ", Price: "1",
		},
		{
			ShortDescription: "asdf 3   ", Price: "1",
		},
		{
			ShortDescription: "asdf  3   ", Price: "1",
		},
		{
			ShortDescription: "asdf", Price: "1",
		},
	}
	var expectedPts = int(math.Ceil(.2) * 3)
	pts, err := calculateItemDescBonus(items)
	if err != nil || pts != expectedPts {
		t.Fatalf(`%d item description should provide %d points but is %d points!`, len(items), expectedPts, pts)
	}
}

func TestItemDescriptionBonusError(t *testing.T) {

	var items = []item{
		{
			ShortDescription: "asdf", Price: "1",
		},
		{
			ShortDescription: "asdf", Price: "asdf",
		},
	}
	pts, err := calculateItemDescBonus(items)
	if err == nil && pts != 0 {
		t.Fatalf(`Items should throw an error!`)
	}
}

func TestDateBonusEven(t *testing.T) {
	dateString := "2009-11-10"
	pts, err := calculateOddDayBonus(dateString)
	if err != nil || pts != 0 {
		t.Fatalf(`%s date should be 0 points!`, dateString)
	}
}

func TestDateBonusOdd(t *testing.T) {
	dateString := "2009-11-11"
	pts, err := calculateOddDayBonus(dateString)
	if err != nil || pts != 6 {
		t.Fatalf(`%s date should be 6 points!`, dateString)
	}
}

func TestDateBonusError(t *testing.T) {
	dateString := "2009-11-75"
	pts, err := calculateOddDayBonus(dateString)
	if err == nil && pts != 0 {
		t.Fatalf(`%s date should be 0 points and throw an error!`, dateString)
	}
}

func TestTimeBonus(t *testing.T) {
	time := "15:00"
	pts, err := calculateTimeOfDayBonus(time)
	if err != nil || pts != 10 {
		t.Fatalf(`%s time should be 10 points!`, time)
	}
}

func TestTimeBonusBefore(t *testing.T) {
	time := "05:00"
	pts, err := calculateTimeOfDayBonus(time)
	if err != nil || pts != 0 {
		t.Fatalf(`%s time should be 0 points!`, time)
	}
}

func TestTimeBonusAfter(t *testing.T) {
	time := "18:00"
	pts, err := calculateTimeOfDayBonus(time)
	if err != nil || pts != 0 {
		t.Fatalf(`%s time should be 0 points!`, time)
	}
}

func TestTimeBonusError(t *testing.T) {
	time := "asdf"
	pts, err := calculateTimeOfDayBonus(time)
	if err == nil || pts != 0 {
		t.Fatalf(`%s time should throw an error!`, time)
	}
}
