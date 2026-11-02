package resourceassert

import "strconv"

func (p *PostgresMirrorResourceAssert) HasRefreshIntervalSeconds(seconds int) *PostgresMirrorResourceAssert {
	p.ValueSet("refresh_interval.#", "1")
	p.ValueSet("refresh_interval.0.seconds", strconv.Itoa(seconds))
	p.ValueSet("refresh_interval.0.minutes", "0")
	p.ValueSet("refresh_interval.0.hours", "0")
	return p
}

func (p *PostgresMirrorResourceAssert) HasRefreshIntervalMinutes(minutes int) *PostgresMirrorResourceAssert {
	p.ValueSet("refresh_interval.#", "1")
	p.ValueSet("refresh_interval.0.minutes", strconv.Itoa(minutes))
	p.ValueSet("refresh_interval.0.seconds", "0")
	p.ValueSet("refresh_interval.0.hours", "0")
	return p
}
