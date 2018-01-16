package mvt

import (
	"math"
	"math/bits"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/internal/mercator"
	"github.com/paulmach/orb/maptile"
	"github.com/paulmach/orb/project"
)

func newProjection(tile maptile.Tile, extent uint32) *project.Projection {
	if isPowerOfTwo(extent) {
		// powers of two extents allows for some more simplicity
		n := uint32(bits.TrailingZeros32(extent))
		z := uint32(tile.Z) + n

		minx := float64(tile.X << n)
		miny := float64(tile.Y << n)
		return &project.Projection{
			ToPlanar: func(p orb.Point) orb.Point {
				x, y := mercator.ToPlanar(p[0], p[1], z)
				return orb.Point{
					math.Floor(x - minx),
					math.Floor(y - miny),
				}
			},
			ToGeo: func(p orb.Point) orb.Point {
				lon, lat := mercator.ToGeo(p[0]+minx+0.5, p[1]+miny+0.5, z)
				return orb.Point{lon, lat}
			},
		}
	}

	// I really don't know why anyone would use a non-power of two extent,
	// but technically it is supported.
	e := float64(extent)
	z := uint32(tile.Z)

	minx := float64(tile.X)
	miny := float64(tile.Y)
	return &project.Projection{
		ToPlanar: func(p orb.Point) orb.Point {
			x, y := mercator.ToPlanar(p[0], p[1], z)
			return orb.Point{
				math.Floor((x - minx) * e),
				math.Floor((y - miny) * e),
			}
		},
		ToGeo: func(p orb.Point) orb.Point {
			lon, lat := mercator.ToGeo((p[0]/e)+minx, (p[1]/e)+miny, z)
			return orb.Point{lon, lat}
		},
	}
}

func isPowerOfTwo(n uint32) bool {
	return (n & (n - 1)) == 0
}