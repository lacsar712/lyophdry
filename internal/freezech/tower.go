package freezech

import (
	"context"
	"time"

	"github.com/lacsar712/lyophdry/internal/airflow"
	"github.com/lacsar712/lyophdry/internal/clock"
	"github.com/lacsar712/lyophdry/internal/config"
	"github.com/lacsar712/lyophdry/internal/model"
	"github.com/lacsar712/lyophdry/internal/shelftemp"
	"github.com/lacsar712/lyophdry/internal/store"
)

type ZoneTable struct {
	tower  model.TowerID
	zones  []model.ZoneAssignment
	plenum model.PlenumID
}

func NewZoneTable(tower model.TowerID, count int, plenum model.PlenumID, cfg config.Config) (*ZoneTable, error) {
	if count <= 0 {
		return nil, model.Wrap("zone_table", "count", model.ErrInvalidID)
	}
	z := &ZoneTable{tower: tower, plenum: plenum}
	for i := 0; i < count; i++ {
		zoneID, err := model.ParseZoneID(tower, i)
		if err != nil {
			return nil, err
		}
		z.zones = append(z.zones, model.ZoneAssignment{
			Zone: zoneID, Plenum: plenum, Enabled: true,
			Setpoint: model.AirflowSetpoint{CubicMetersPerHour: cfg.DefaultAirflowCMH, TolerancePct: cfg.AirflowTolerancePct},
			TargetMoist: cfg.TargetMoistPct,
		})
	}
	return z, nil
}

func (z *ZoneTable) Zones() []model.ZoneAssignment {
	out := make([]model.ZoneAssignment, len(z.zones))
	copy(out, z.zones)
	return out
}

func (z *ZoneTable) EnabledCount() int {
	n := 0
	for _, zone := range z.zones {
		if zone.Enabled {
			n++
		}
	}
	return n
}

func (z *ZoneTable) UpdateMoisture(zone model.ZoneID, pct float64) {
	for i := range z.zones {
		if z.zones[i].Zone == zone {
			z.zones[i].ActualMoist = pct
			return
		}
	}
}

func (z *ZoneTable) UpdateFlow(zone model.ZoneID, cmh float64) {
	for i := range z.zones {
		if z.zones[i].Zone == zone {
			z.zones[i].LastFlow = cmh
			return
		}
	}
}

type LyoPlant struct {
	cfg       config.Config
	clk       clock.Clock
	mem       *store.Memory
	plenums   *airflow.PlenumTable
	fans      *airflow.FanBank
	sensors   *shelftemp.SensorBank
	gradient  *shelftemp.GradientController
	profile   *shelftemp.ProfileManager
	holdStart time.Time
	holdDur   time.Duration
}

func NewLyoPlant(cfg config.Config, clk clock.Clock, mem *store.Memory) *LyoPlant {
	return &LyoPlant{
		cfg: cfg, clk: clk, mem: mem,
		plenums:  airflow.NewPlenumTable(),
		fans:     airflow.NewFanBank(),
		sensors:  shelftemp.NewSensorBank(),
		gradient: shelftemp.NewGradientController(cfg.MaxGradientDeltaPct),
	}
}

func (p *LyoPlant) Plenums() *airflow.PlenumTable { return p.plenums }

func (p *LyoPlant) Fans() *airflow.FanBank { return p.fans }

func (p *LyoPlant) BindAirflow(plenum model.PlenumID, sp model.AirflowSetpoint) {
	if pl, ok := p.plenums.Get(plenum); ok {
		pl.BindSetpoint(sp)
	}
}

func (p *LyoPlant) Coordinator() *Coordinator { return NewCoordinator(p.cfg, p.clk, p.fans) }

func (p *LyoPlant) PrimePlenum(ctx context.Context, plenum model.PlenumID) error {
	pl, ok := p.plenums.Get(plenum)
	if !ok {
		return model.Wrap("tower_plant", "plenum", model.ErrNotFound)
	}
	dur := time.Duration(p.cfg.PlenumPrimeSec) * time.Second
	return pl.Prime(ctx, p.clk, dur)
}

func (p *LyoPlant) ObserveFlow(plenum model.PlenumID, cmh float64) {
	if pl, ok := p.plenums.Get(plenum); ok {
		pl.ObserveFlow(cmh)
	}
}

func (p *LyoPlant) ValidateFlows(ctx context.Context) error {
	return p.plenums.ValidateAll()
}

func (p *LyoPlant) RegisterSensor(sensor *shelftemp.Sensor) {
	p.sensors.Register(sensor)
}

func (p *LyoPlant) ObserveMoisture(zone model.ZoneID, pct float64) error {
	_, err := p.sensors.ObserveZone(zone, pct, p.clk.Now())
	return err
}

func (p *LyoPlant) ValidateGradient() error {
	return p.gradient.Validate(p.sensors.Readings())
}

func (p *LyoPlant) InitProfile(zones []model.ZoneID, targets []float64) error {
	pm, err := shelftemp.NewProfileManager(zones, targets)
	if err != nil {
		return err
	}
	p.profile = pm
	return nil
}

func (p *LyoPlant) ArmMoistureHold(start time.Time, duration time.Duration, targetPct float64) {
	if p.profile == nil {
		return
	}
	p.profile.Window(start, duration, targetPct)
	p.holdStart = start
	p.holdDur = duration
}

func (p *LyoPlant) HoldActive() bool {
	if p.profile == nil {
		return false
	}
	return p.profile.HoldActive(p.clk.Now)
}

func (p *LyoPlant) ReleaseHold() {
	if p.profile != nil {
		p.profile.ReleaseHold()
	}
}

func (p *LyoPlant) AtTarget(tolerance float64) bool {
	if p.profile == nil {
		return false
	}
	return p.profile.AllAtTarget(p.sensors.Readings(), tolerance)
}

func (p *LyoPlant) GradientDelta() float64 {
	return p.gradient.Delta(p.sensors.Readings())
}

func (p *LyoPlant) GradientDeltaFor(readings []model.MoistureReading) float64 {
	return p.gradient.Delta(readings)
}

func (p *LyoPlant) SensorReadings() []model.MoistureReading {
	return p.sensors.Readings()
}

func (p *LyoPlant) Profile() *shelftemp.ProfileManager {
	return p.profile
}

func (p *LyoPlant) HoldDuration() time.Duration {
	if p.holdDur > 0 {
		return p.holdDur
	}
	return time.Minute
}
