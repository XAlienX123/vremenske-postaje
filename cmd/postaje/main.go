// Command postaje demonstrates the vreme package with a few weather stations.
//
// Usage:
//
//	postaje [--minTemp -50] [--maxTemp 50] [--stMeritev 3]
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/XAlienX123/vremenske-postaje/vreme"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "postaje",
		Usage: "vremenske postaje in poročila o temperaturah",
		Flags: []cli.Flag{
			&cli.FloatFlag{Name: "minTemp", Value: -50, Usage: "najnižja dovoljena temperatura meritve (°C)"},
			&cli.FloatFlag{Name: "maxTemp", Value: 50, Usage: "najvišja dovoljena temperatura meritve (°C)"},
			&cli.IntFlag{Name: "stMeritev", Value: 3, Usage: "najmanjše število meritev za izračun statistike"},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			run(cmd.Float("minTemp"), cmd.Float("maxTemp"), int(cmd.Int("stMeritev")))
			return nil
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(minTemp, maxTemp float64, stMeritev int) {
	postaje := map[string]vreme.Postaja{}
	vreme.DodajPostajo(postaje, "LJ01", "Ljubljana Bežigrad", "Ljubljana")
	vreme.DodajPostajo(postaje, "MB01", "Maribor Tabor", "Maribor")
	vreme.DodajPostajo(postaje, "KP01", "Koper Luka", "Koper")
	vreme.DodajPostajo(postaje, "KR01", "Kredarica", "Triglav")

	zacetne := map[string][]float64{
		"LJ01": {18.2, 21.5, 19.8, 16.5},
		"MB01": {12.4},
		"KP01": {24.5, 27, 26.1, 28.4},
		"KR01": {2, -1.5, 3.5},
	}
	for oznaka, meritve := range zacetne {
		for _, t := range meritve {
			vreme.DodajMeritev(postaje, oznaka, t, minTemp, maxTemp)
		}
	}

	vreme.IzpisMeritev(postaje)
	fmt.Println()

	vreme.DodajMeritev(postaje, "MB01", 14.1, minTemp, maxTemp)
	vreme.DodajMeritev(postaje, "MB01", 13.3, minTemp, maxTemp)
	vreme.DodajMeritev(postaje, "LJ01", 75, minTemp, maxTemp)
	vreme.DodajMeritev(postaje, "XX99", 20, minTemp, maxTemp)
	fmt.Println()

	vreme.IzpisMeritev(postaje)
	fmt.Println()
	vreme.IzpisPorocila(postaje, stMeritev)

	vreme.DodajPostajo(postaje, "NG01", "Nova Gorica", "Nova Gorica")
	vreme.DodajMeritev(postaje, "NG01", 20, minTemp, maxTemp)
	fmt.Println()
	vreme.IzpisPorocila(postaje, stMeritev)
}
