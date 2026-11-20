// Package vreme provides weather stations that store temperature measurements
// and print measurement lists and statistical reports.
//
// Stations are kept in a map keyed by station code. Create stations with
// [DodajPostajo], add measurements with [DodajMeritev] and print them with
// [IzpisMeritev] and [IzpisPorocila].
package vreme

import (
	"fmt"
	"maps"
	"slices"
)

// Postaja is a weather station with a name, a location and a list of
// temperature measurements in degrees Celsius.
// Its fields are unexported, so stations are created with [DodajPostajo].
type Postaja struct {
	ime     string
	kraj    string
	meritve []float64
}

// DodajPostajo adds a new station with the given code, name and location to
// postaje. It prints a message and does nothing if the code already exists.
func DodajPostajo(postaje map[string]Postaja, oznaka, ime, kraj string) {
	if _, ok := postaje[oznaka]; ok {
		fmt.Printf("Postaja %q že obstaja.\n", oznaka)
		return
	}
	postaje[oznaka] = Postaja{ime: ime, kraj: kraj}
}

// DodajMeritev adds the temperature measurement to the station with the given
// code. The measurement is rejected with a message if the station does not
// exist or the temperature is outside the range [minTemp, maxTemp].
func DodajMeritev(postaje map[string]Postaja, oznaka string, temperatura, minTemp, maxTemp float64) {
	p, ok := postaje[oznaka]
	if !ok {
		fmt.Printf("Postaja %q ne obstaja.\n", oznaka)
		return
	}
	if temperatura < minTemp || temperatura > maxTemp {
		fmt.Printf("Neveljavna meritev %.1f °C za postajo %q (dovoljeno od %.1f do %.1f °C).\n", temperatura, oznaka, minTemp, maxTemp)
		return
	}
	p.meritve = append(p.meritve, temperatura)
	postaje[oznaka] = p
}

// statistika returns the minimum, maximum and average temperature of the
// station. ok is false if the station does not exist or has fewer than
// stMeritev measurements.
func statistika(postaje map[string]Postaja, oznaka string, stMeritev int) (najnizja, najvisja, povprecna float64, ok bool) {
	p, obstaja := postaje[oznaka]
	if !obstaja || len(p.meritve) == 0 || len(p.meritve) < stMeritev {
		return
	}
	najnizja, najvisja = p.meritve[0], p.meritve[0]
	vsota := 0.0
	for _, m := range p.meritve {
		najnizja = min(najnizja, m)
		najvisja = max(najvisja, m)
		vsota += m
	}
	return najnizja, najvisja, vsota / float64(len(p.meritve)), true
}

// IzpisMeritev prints all stations, sorted by code, with their measurements.
func IzpisMeritev(postaje map[string]Postaja) {
	fmt.Println("MERITVE:")
	for _, oznaka := range slices.Sorted(maps.Keys(postaje)) {
		p := postaje[oznaka]
		fmt.Printf("%s - %s (%s): %v\n", oznaka, p.ime, p.kraj, p.meritve)
	}
}

// IzpisPorocila prints the minimum, maximum and average temperature of each
// station, sorted by code, with a description of the average (Vroče, Toplo,
// Hladno or Mrzlo). Stations with fewer than stMeritev measurements are
// reported as having too few measurements.
func IzpisPorocila(postaje map[string]Postaja, stMeritev int) {
	for _, oznaka := range slices.Sorted(maps.Keys(postaje)) {
		ime := postaje[oznaka].ime
		mn, mx, avg, ok := statistika(postaje, oznaka, stMeritev)
		if !ok {
			fmt.Printf("%s: premalo meritev\n", ime)
			continue
		}
		var opis string
		switch {
		case avg >= 25:
			opis = "Vroče"
		case avg >= 15:
			opis = "Toplo"
		case avg >= 5:
			opis = "Hladno"
		default:
			opis = "Mrzlo"
		}
		fmt.Printf("%s: min %.1f °C, max %.1f °C, povprečje %.1f °C -> %s\n", ime, mn, mx, avg, opis)
	}
}
