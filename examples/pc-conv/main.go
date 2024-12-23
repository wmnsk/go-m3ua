package main

import (
	"flag"
	"log"
	"strings"

	"github.com/dmisol/go-m3ua/pc"
)

func main() {
	var (
		raw     = flag.Int("raw", -1, "Raw PC in integer(decimal or hex).")
		str     = flag.String("str", "", "Formatted PC in string, splitted with \"-\".")
		variant = flag.String("variant", "3-8-3", "Variant of PC to convert to/to be converted from.")
	)
	flag.Parse()

	if *raw < 0 && *str == "" {
		log.Fatal("Invalid command-line flag: -raw or -str is required.")
	}

	if *raw > 0 {
		p := pc.NewPointCode(uint32(*raw), pc.Variant(*variant))
		log.Printf("PC successfully converted.\n\tRaw: %d, Formatted: %s, Variant: %s", p.Uint32(), p.String(), p.Variant())
	}

	if *str != "" {
		if len(strings.Split(*str, "-")) < 2 {
			log.Fatalf("Invalid formatted PC given: %s, should be splitted with \"-\".", *str)
		}
		p := pc.NewPointCodeFrom(*str, pc.Variant(*variant))
		log.Printf("PC successfully converted.\n\tRaw: %d, Formatted: %s, Variant: %s", p.Uint32(), p.String(), p.Variant())
	}
}
