package block

import (
	"io"
	"time"
	"unsafe"
)

var (
	_ = unsafe.Sizeof(0)
	_ = io.ReadFull
	_ = time.Now()
)

type BlockHead struct {
	Version    int64
	ParentHash []byte
	TreeHash   []byte
	Info       []byte
	Number     int64
	Witness    string
	Signature  []byte
	Time       int64
}

func (d *BlockHead) Size() (s uint64) {

	{
		l := uint64(len(d.ParentHash))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	{
		l := uint64(len(d.TreeHash))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	{
		l := uint64(len(d.Info))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	{
		l := uint64(len(d.Witness))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	{
		l := uint64(len(d.Signature))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	s += 24
	return
}
func (d *BlockHead) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{

		buf[0+0] = byte(d.Version >> 0)

		buf[1+0] = byte(d.Version >> 8)

		buf[2+0] = byte(d.Version >> 16)

		buf[3+0] = byte(d.Version >> 24)

		buf[4+0] = byte(d.Version >> 32)

		buf[5+0] = byte(d.Version >> 40)

		buf[6+0] = byte(d.Version >> 48)

		buf[7+0] = byte(d.Version >> 56)

	}
	{
		l := uint64(len(d.ParentHash))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+8] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+8] = byte(t)
			i++

		}
		copy(buf[i+8:], d.ParentHash)
		i += l
	}
	{
		l := uint64(len(d.TreeHash))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+8] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+8] = byte(t)
			i++

		}
		copy(buf[i+8:], d.TreeHash)
		i += l
	}
	{
		l := uint64(len(d.Info))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+8] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+8] = byte(t)
			i++

		}
		copy(buf[i+8:], d.Info)
		i += l
	}
	{

		buf[i+0+8] = byte(d.Number >> 0)

		buf[i+1+8] = byte(d.Number >> 8)

		buf[i+2+8] = byte(d.Number >> 16)

		buf[i+3+8] = byte(d.Number >> 24)

		buf[i+4+8] = byte(d.Number >> 32)

		buf[i+5+8] = byte(d.Number >> 40)

		buf[i+6+8] = byte(d.Number >> 48)

		buf[i+7+8] = byte(d.Number >> 56)

	}
	{
		l := uint64(len(d.Witness))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+16] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+16] = byte(t)
			i++

		}
		copy(buf[i+16:], d.Witness)
		i += l
	}
	{
		l := uint64(len(d.Signature))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+16] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+16] = byte(t)
			i++

		}
		copy(buf[i+16:], d.Signature)
		i += l
	}
	{

		buf[i+0+16] = byte(d.Time >> 0)

		buf[i+1+16] = byte(d.Time >> 8)

		buf[i+2+16] = byte(d.Time >> 16)

		buf[i+3+16] = byte(d.Time >> 24)

		buf[i+4+16] = byte(d.Time >> 32)

		buf[i+5+16] = byte(d.Time >> 40)

		buf[i+6+16] = byte(d.Time >> 48)

		buf[i+7+16] = byte(d.Time >> 56)

	}
	return buf[:i+24], nil
}

func (d *BlockHead) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{

		d.Version = 0 | (int64(buf[i+0+0]) << 0) | (int64(buf[i+1+0]) << 8) | (int64(buf[i+2+0]) << 16) | (int64(buf[i+3+0]) << 24) | (int64(buf[i+4+0]) << 32) | (int64(buf[i+5+0]) << 40) | (int64(buf[i+6+0]) << 48) | (int64(buf[i+7+0]) << 56)

	}
	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+8] & 0x7F)
			for buf[i+8]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+8]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		if uint64(cap(d.ParentHash)) >= l {
			d.ParentHash = d.ParentHash[:l]
		} else {
			d.ParentHash = make([]byte, l)
		}
		copy(d.ParentHash, buf[i+8:])
		i += l
	}
	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+8] & 0x7F)
			for buf[i+8]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+8]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		if uint64(cap(d.TreeHash)) >= l {
			d.TreeHash = d.TreeHash[:l]
		} else {
			d.TreeHash = make([]byte, l)
		}
		copy(d.TreeHash, buf[i+8:])
		i += l
	}
	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+8] & 0x7F)
			for buf[i+8]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+8]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		if uint64(cap(d.Info)) >= l {
			d.Info = d.Info[:l]
		} else {
			d.Info = make([]byte, l)
		}
		copy(d.Info, buf[i+8:])
		i += l
	}
	{

		d.Number = 0 | (int64(buf[i+0+8]) << 0) | (int64(buf[i+1+8]) << 8) | (int64(buf[i+2+8]) << 16) | (int64(buf[i+3+8]) << 24) | (int64(buf[i+4+8]) << 32) | (int64(buf[i+5+8]) << 40) | (int64(buf[i+6+8]) << 48) | (int64(buf[i+7+8]) << 56)

	}
	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+16] & 0x7F)
			for buf[i+16]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+16]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		d.Witness = string(buf[i+16 : i+16+l])
		i += l
	}
	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+16] & 0x7F)
			for buf[i+16]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+16]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		if uint64(cap(d.Signature)) >= l {
			d.Signature = d.Signature[:l]
		} else {
			d.Signature = make([]byte, l)
		}
		copy(d.Signature, buf[i+16:])
		i += l
	}
	{

		d.Time = 0 | (int64(buf[i+0+16]) << 0) | (int64(buf[i+1+16]) << 8) | (int64(buf[i+2+16]) << 16) | (int64(buf[i+3+16]) << 24) | (int64(buf[i+4+16]) << 32) | (int64(buf[i+5+16]) << 40) | (int64(buf[i+6+16]) << 48) | (int64(buf[i+7+16]) << 56)

	}
	return i + 24, nil
}

type BlockRaw struct {
	Head    BlockHead
	Content [][]byte
}

func (d *BlockRaw) Size() (s uint64) {

	{
		s += d.Head.Size()
	}
	{
		l := uint64(len(d.Content))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}

		for k0 := range d.Content {

			{
				l := uint64(len(d.Content[k0]))

				{

					t := l
					for t >= 0x80 {
						t >>= 7
						s++
					}
					s++

				}
				s += l
			}

		}

	}
	return
}
func (d *BlockRaw) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{
		nbuf, err := d.Head.Marshal(buf[0:])
		if err != nil {
			return nil, err
		}
		i += uint64(len(nbuf))
	}
	{
		l := uint64(len(d.Content))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+0] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+0] = byte(t)
			i++

		}
		for k0 := range d.Content {

			{
				l := uint64(len(d.Content[k0]))

				{

					t := uint64(l)

					for t >= 0x80 {
						buf[i+0] = byte(t) | 0x80
						t >>= 7
						i++
					}
					buf[i+0] = byte(t)
					i++

				}
				copy(buf[i+0:], d.Content[k0])
				i += l
			}

		}
	}
	return buf[:i+0], nil
}

func (d *BlockRaw) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{
		ni, err := d.Head.Unmarshal(buf[i+0:])
		if err != nil {
			return 0, err
		}
		i += ni
	}
	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+0] & 0x7F)
			for buf[i+0]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+0]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		if uint64(cap(d.Content)) >= l {
			d.Content = d.Content[:l]
		} else {
			d.Content = make([][]byte, l)
		}
		for k0 := range d.Content {

			{
				l := uint64(0)

				{

					bs := uint8(7)
					t := uint64(buf[i+0] & 0x7F)
					for buf[i+0]&0x80 == 0x80 {
						i++
						t |= uint64(buf[i+0]&0x7F) << bs
						bs += 7
					}
					i++

					l = t

				}
				if uint64(cap(d.Content[k0])) >= l {
					d.Content[k0] = d.Content[k0][:l]
				} else {
					d.Content[k0] = make([]byte, l)
				}
				copy(d.Content[k0], buf[i+0:])
				i += l
			}

		}
	}
	return i + 0, nil
}
