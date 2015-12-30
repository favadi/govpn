/*
GoVPN -- simple secure free software virtual private network daemon
Copyright (C) 2014-2016 Sergey Matveev <stargrave@stargrave.org>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package govpn

import (
	"testing"
	"time"
)

var (
	peer       *Peer
	plaintext  []byte
	ciphertext []byte
	peerId     *PeerId
	conf       *PeerConf
)

type Dummy struct {
	dst *[]byte
}

func (d Dummy) Write(b []byte) (int, error) {
	if d.dst != nil {
		*d.dst = b
	}
	return len(b), nil
}

func init() {
	MTU = 1500
	id := new([IDSize]byte)
	peerId := PeerId(*id)
	conf = &PeerConf{
		Id:      &peerId,
		Timeout: time.Second * time.Duration(TimeoutDefault),
		Noise:   false,
		CPR:     0,
	}
	peer = newPeer(true, "foo", Dummy{&ciphertext}, conf, new([SSize]byte))
	plaintext = make([]byte, 789)
}

func BenchmarkEnc(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		peer.EthProcess(plaintext)
	}
}

func BenchmarkDec(b *testing.B) {
	peer = newPeer(true, "foo", Dummy{&ciphertext}, conf, new([SSize]byte))
	peer.EthProcess(plaintext)
	peer = newPeer(true, "foo", Dummy{nil}, conf, new([SSize]byte))
	orig := make([]byte, len(ciphertext))
	copy(orig, ciphertext)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		peer.nonceBucket0 = make(map[uint64]struct{}, 1)
		peer.nonceBucket1 = make(map[uint64]struct{}, 1)
		copy(ciphertext, orig)
		if !peer.PktProcess(ciphertext, Dummy{nil}, true) {
			b.Fail()
		}
	}
}
