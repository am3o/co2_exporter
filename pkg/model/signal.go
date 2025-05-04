package model

type Signal []byte

func (s Signal) Type() rune {
	if len(s) == 0 {
		return 0
	}

	return rune(s[0])
}

func (s Signal) Value() float64 {
	if len(s) <= 2 {
		return 0
	}
	return float64(((int)(s[1]) << 8) | (int)(s[2]))
}
