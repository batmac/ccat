package pcg

type PCG32 struct {
	state uint64
	inc   uint64
}

func (pcg *PCG32) Next() uint32 {
	oldstate := pcg.state
	pcg.state = oldstate*6364136223846793005 + (pcg.inc | 1)
	xorshifted := uint32(((oldstate >> 18) ^ oldstate) >> 27)
	rot := uint32(oldstate >> 59)
	return (xorshifted >> rot) | (xorshifted << ((-rot) & 31))
}

func NewPCG32(initState uint64, initSeq uint64) *PCG32 {
	pcg := &PCG32{0, (initSeq << 1) | 1}
	pcg.Next()
	pcg.state += initState
	pcg.Next()
	return pcg
}
