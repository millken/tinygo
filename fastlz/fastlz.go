package fastlz

// based on https://github.com/Vectorized/solady/blob/5315d937d79b335c668896d7533ac603adac5315/js/solady.js

func FlzDecompress(ib []byte) []byte {
	j := uint32(0)
	ob := make([]byte, 0, 2*len(ib))
	for i := uint32(0); i < uint32(len(ib)); {
		t := uint32(ib[i]) >> 5
		if t == 0 {
			i++
			for j = uint32(ib[i-1]) + 1; j > 0; j-- {
				ob = append(ob, ib[i])
				i++
			}
		} else {
			if t < 7 {
				j = 2 + t
				t = 1
			} else {
				j = 9 + uint32(ib[i+1])
				t = 0
			}
			f := 256*(uint32(ib[i])&31) + uint32(ib[i+2-t])
			i = i + 3 - t
			r := uint32(len(ob)) - f - 1
			for ; j > 0; j-- {
				ob = append(ob, ob[r])
				r++
			}
		}
	}
	return ob
}

func FlzCompress(ib []byte) []byte {
	ob := make([]byte, 0, len(ib))
	app := func(b ...byte) {
		ob = append(ob, b...)
	}
	ht := make([]uint32, 8192)
	u24 := func(i uint32) uint32 {
		return uint32(ib[i]) | (uint32(ib[i+1]) << 8) | (uint32(ib[i+2]) << 16)
	}
	cmp := func(p uint32, q uint32, e uint32) uint32 {
		l := uint32(0)
		for e -= q; l < e; l++ {
			if ib[p+l] != ib[q+l] {
				e = 0
			}
		}
		return l
	}
	literals := func(r uint32, s uint32) {
		for ; r >= 0x20; r -= 0x20 {
			app(31)
			app(ib[s : s+0x20]...)
			s += 0x20
		}
		if r != 0 {
			app(byte(r - 1))
			app(ib[s : s+r]...)
		}
	}
	match := func(l uint32, d uint32) {
		for d--; l >= 263; l -= 262 {
			app(byte(224 + (d >> 8)))
			app(253)
			app(byte(d & 0xff))
		}
		if l >= 7 {
			app(byte(224 + (d >> 8)))
			app(byte(l - 7))
			app(byte(d & 0xff))
		} else {
			app(byte((l << 5) + (d >> 8)))
			app(byte(d & 0xff))
		}
	}
	hash := func(v uint32) uint32 {
		return ((2654435769 * v) >> 19) & 0x1fff
	}
	setNextHash := func(ip uint32) uint32 {
		ht[hash(u24(ip))] = ip
		return ip + 1
	}
	a := uint32(0)
	ipLimit := uint32(len(ib)) - 13
	if len(ib) < 13 {
		ipLimit = 0
	}
	for ip := a + 2; ip < ipLimit; {
		r := uint32(0)
		d := uint32(0)
		for {
			s := u24(ip)
			h := hash(s)
			r = ht[h]
			ht[h] = ip
			d = ip - r
			if ip >= ipLimit {
				break
			}
			ip++
			if d <= 0x1fff && s == u24(r) {
				break
			}
		}
		if ip >= ipLimit {
			break
		}
		ip--
		if ip > a {
			literals(ip-a, a)
		}
		l := cmp(r+3, ip+3, ipLimit+9)
		match(l, d)
		ip = setNextHash(setNextHash(ip + l))
		a = ip
	}
	literals(uint32(len(ib))-a, a)
	return ob
}
