package relational

import (
	"github.com/jedib0t/go-pretty/v6/progress"
	"testing"
	"time"
)

func TestProgress(t *testing.T) {
	pw := progress.NewWriter()

	pw.SetAutoStop(false)
	pw.SetTrackerLength(25)
	pw.SetMessageWidth(24)
	pw.SetNumTrackersExpected(13)
	pw.SetSortBy(progress.SortByPercentDsc)
	pw.SetStyle(progress.StyleDefault)
	pw.SetTrackerPosition(progress.PositionRight)
	pw.SetUpdateFrequency(time.Millisecond * 100)
	pw.Style().Colors = progress.StyleColorsExample
	pw.Style().Options.PercentFormat = "%4.1f%%"
	pw.Style().Visibility.ETA = true
	pw.Style().Visibility.ETAOverall = true
	pw.Style().Visibility.Percentage = true
	pw.Style().Visibility.Speed = false
	pw.Style().Visibility.SpeedOverall = false
	pw.Style().Visibility.Time = true
	//pw.Style().Visibility.TrackerOverall = true
	pw.Style().Visibility.Value = true
	pw.Style().Visibility.Pinned = false

	go pw.Render()

	go func() {
		t := &progress.Tracker{
			Message: "Tracker message",
			Total:   200,
			//ExpectedDuration: time.Second * 2,
			Units: progress.UnitsDefault,
		}
		pw.AppendTracker(t)
		ticker := time.Tick(time.Second)
		for !t.IsDone() {
			select {
			case <-ticker:
				t.Increment(50)
			}
		}
	}()

	select {}
}
