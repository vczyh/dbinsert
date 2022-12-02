package relation

import (
	"github.com/jedib0t/go-pretty/v6/progress"
)

type Progress struct {
	tables        []*Table
	dbs           map[string][]*Table
	dbTrackers    map[string]*progress.Tracker
	tableTrackers map[*Table]*progress.Tracker
	pw            progress.Writer
}

func NewProgress(tables []*Table) (*Progress, error) {
	p := &Progress{
		tables:        tables,
		dbs:           make(map[string][]*Table),
		dbTrackers:    make(map[string]*progress.Tracker),
		tableTrackers: make(map[*Table]*progress.Tracker),
	}
	for _, table := range p.tables {
		p.dbs[table.Database] = append(p.dbs[table.Database], table)
	}

	return p, nil
}

func (p *Progress) Render() {
	pw := progress.NewWriter()
	pw.SetAutoStop(true)
	pw.SetSortBy(progress.SortByPercent)
	pw.SetStyle(progress.StyleDefault)
	pw.Style().Colors = progress.StyleColorsExample
	pw.SetTrackerPosition(progress.PositionLeft)
	pw.SetPinnedMessages("All Databases Progress")

	dbSize := make(map[string]int)
	for _, table := range p.tables {
		dbSize[table.Database] += table.Size
	}

	for db, tables := range p.dbs {
		var dbSize int
		for _, table := range tables {
			tableTracker := &progress.Tracker{
				Message: table.Name,
				Total:   int64(table.Size),
				Units:   progress.UnitsDefault,
			}
			p.tableTrackers[table] = tableTracker
			dbSize += table.Size
		}

		dbTracker := &progress.Tracker{
			Message: db,
			Total:   int64(dbSize),
			//ExpectedDuration: time.Second * 2,
			Units: progress.UnitsDefault,
		}
		p.dbTrackers[db] = dbTracker
	}

	for _, tracker := range p.dbTrackers {
		pw.AppendTracker(tracker)
	}
	p.pw = pw
	p.pw.Render()
}

func (p *Progress) Stop() {
	p.pw.Stop()
}

func (p *Progress) Increment(table *Table, inc int) {
	//tabletTracker := p.tableTrackers[table]
	dbTracker := p.dbTrackers[table.Database]
	//tabletTracker.Increment(int64(inc))
	dbTracker.Increment(int64(inc))
}

func (p *Progress) Ended() bool {
	return !p.pw.IsRenderInProgress()
}
