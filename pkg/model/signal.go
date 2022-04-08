package model

type Signal []byte

func (s Signal) Type() rune {
	return rune(s[0])
}

func (s Signal) Value() float64 {
	return float64(((int)(s[1]) << 8) | (int)(s[2]))
}
