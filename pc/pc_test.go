package pc_test

import (
	"testing"

	"github.com/wmnsk/go-m3ua/pc"
)

// TODO: coverage...

func TestConvertPointCode(t *testing.T) {
	cases := []struct {
		name             string
		pc               *pc.PointCode
		currVar, nextVar pc.Variant
		before, after    string
	}{
		{
			name:    "1234/3-8-3 to 4-3-7",
			pc:      pc.NewPointCode(1234, pc.Variant383),
			currVar: pc.Variant383,
			nextVar: pc.Variant437,
			before:  "0-154-2",
			after:   "1-1-82",
		}, {
			name:    "0xffffffff/3-8-3 to 4-3-7",
			pc:      pc.NewPointCode(0xffffffff, pc.Variant383),
			currVar: pc.Variant383,
			nextVar: pc.Variant437,
			before:  "7-255-7",
			after:   "15-7-127",
		}, {
			name:    "0/3-8-3 to 4-3-7",
			pc:      pc.NewPointCode(0, pc.Variant383),
			currVar: pc.Variant383,
			nextVar: pc.Variant437,
			before:  "0-0-0",
			after:   "0-0-0",
		},
	}

	for _, c := range cases {
		if got, want := c.pc.String(), c.before; got != want {
			t.Errorf("NewPointCode failed. got: %s, want: %s", got, want)
		}

		if got, want := pc.NewPointCodeFrom(c.before, c.currVar).Uint32(), c.pc.Uint32(); got != want {
			t.Errorf("NewPointCodeFrom failed. got: %d, want: %d", got, want)
		}

		if got, err := c.pc.ConvertTo(c.nextVar); err != nil {
			t.Fatalf("Failed to convert %s to %s", c.pc.Variant(), c.nextVar)
		} else {
			want := c.after
			if got != want {
				t.Errorf("ConvertTo failed. got: %s, want: %s", got, want)
			}
		}

		if got, want := pc.NewPointCodeFrom(c.after, c.nextVar).Uint32(), c.pc.Uint32(); got != want {
			t.Errorf("NewPointCodeFrom failed. got: %d, want: %d", got, want)
		}
	}
}
