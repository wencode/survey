package survey

import (
	"fmt"
	"testing"
	"time"
)

func TestSurveyVars(t *testing.T) {
	SetOutoutDir("./")

	var0 := NewAdder("survey_test", "sum")
	var1 := NewMiner("survey_test", "min")
	var2 := NewMaxer("survey_test", "max")

	NewWindow("survey_test",
		"sum_change_value",
		WithWindowParam(10, var0))
	NewWindow("survey_test",
		"min_change_value",
		WithWindowParam(10, var1))
	NewWindow("survey_test",
		"max_change_value",
		WithWindowParam(10, var2))

	NewUnitWindow("survey_test",
		"sum_rate",
		WithWindowParam(10, var0))
	NewUnitWindow("survey_test",
		"min_rate",
		WithWindowParam(10, var1))
	NewUnitWindow("survey_test",
		"max_rate",
		WithWindowParam(10, var2))

	for i := 0; i < 10; i++ {
		var0.Put(i)
		var1.Put(i)
		var2.Put(i)
		fmt.Print(".")
		time.Sleep(time.Second)
	}
	fmt.Println()

	<-time.After(time.Second * 1)

	Quit()

}

func TestLatency(t *testing.T) {
	SetOutoutDir("./")
	ltn := NewLatency("latency", "write")

	for i := 0; i < 10; i++ {
		lv := i * 100 * 1000
		time.Sleep(time.Microsecond * time.Duration(lv))
		ltn.Put(lv)
		fmt.Printf("%d->{%d}\t", lv, ltn.Get())
	}
	fmt.Println()

	Quit()
}
