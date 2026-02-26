package nmea

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

var vhw = []struct {
	name string
	raw  string
	err  string
	msg  VHW
}{
	{
		name: "good sentence",
		raw:  "$VWVHW,45.0,T,43.0,M,3.5,N,6.4,K*56",
		msg: VHW{
			TrueHeading:            Float64{Value: 45.0, Valid: true},
			MagneticHeading:        Float64{Value: 43.0, Valid: true},
			SpeedThroughWaterKnots: Float64{Value: 3.5, Valid: true},
			SpeedThroughWaterKPH:   Float64{Value: 6.4, Valid: true},
		},
	},
	{
		name: "partial sentence",
		raw:  "$INVHW,187.9,T,,,19.6,N,36.3,K*3E",
		msg: VHW{
			TrueHeading:            Float64{Value: 187.9, Valid: true},
			MagneticHeading:        Float64{Value: 0, Valid: false},
			SpeedThroughWaterKnots: Float64{Value: 19.6, Valid: true},
			SpeedThroughWaterKPH:   Float64{Value: 36.3, Valid: true},
		},
	},
	{
		name: "bad sentence",
		raw:  "$VWVHW,T,45.0,43.0,M,3.5,N,6.4,K*56",
		err:  "nmea: VWVHW invalid true heading: T",
	},
}

func TestVHW(t *testing.T) {
	for _, tt := range vhw {
		t.Run(tt.name, func(t *testing.T) {
			m, err := Parse(tt.raw)
			if tt.err != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				vhw := m.(VHW)
				vhw.BaseSentence = BaseSentence{}
				assert.Equal(t, tt.msg, vhw)
			}
		})
	}
}

func TestVHW_NaNForEmptyFloat(t *testing.T) {
	// Test that SentenceParser with NaNForEmptyFloat returns NaN for empty float fields
	p := SentenceParser{
		NaNForEmptyFloat: true,
	}

	t.Run("partial sentence with NaN", func(t *testing.T) {
		m, err := p.Parse("$INVHW,187.9,T,,,19.6,N,36.3,K*3E")
		assert.NoError(t, err)
		vhw := m.(VHW)

		assert.Equal(t, 187.9, vhw.TrueHeading.Value)
		assert.True(t, vhw.TrueHeading.Valid)

		assert.True(t, math.IsNaN(vhw.MagneticHeading.Value), "empty MagneticHeading should be NaN")
		assert.False(t, vhw.MagneticHeading.Valid)

		assert.Equal(t, 19.6, vhw.SpeedThroughWaterKnots.Value)
		assert.True(t, vhw.SpeedThroughWaterKnots.Valid)

		assert.Equal(t, 36.3, vhw.SpeedThroughWaterKPH.Value)
		assert.True(t, vhw.SpeedThroughWaterKPH.Valid)
	})

	t.Run("full sentence unaffected", func(t *testing.T) {
		m, err := p.Parse("$VWVHW,45.0,T,43.0,M,3.5,N,6.4,K*56")
		assert.NoError(t, err)
		vhw := m.(VHW)

		assert.Equal(t, 45.0, vhw.TrueHeading.Value)
		assert.True(t, vhw.TrueHeading.Valid)

		assert.Equal(t, 43.0, vhw.MagneticHeading.Value)
		assert.True(t, vhw.MagneticHeading.Valid)

		assert.Equal(t, 3.5, vhw.SpeedThroughWaterKnots.Value)
		assert.True(t, vhw.SpeedThroughWaterKnots.Valid)

		assert.Equal(t, 6.4, vhw.SpeedThroughWaterKPH.Value)
		assert.True(t, vhw.SpeedThroughWaterKPH.Valid)
	})
}
