/*
 * LF: Global Fully Replicated Key/Value Store
 * Copyright (C) 2018-2019  ZeroTier, Inc.  https://www.zerotier.com/
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 *
 * --
 *
 * You can be released from the requirements of the license by purchasing
 * a commercial license. Buying such a license is mandatory as soon as you
 * develop commercial closed-source software that incorporates or links
 * directly against ZeroTier software without disclosing the source code
 * of your own application.
 */

package lf

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
)

// seededPrng is a deterministic cryptographic random Reader used to generate key pairs from specific seeds.
// This is used in various places to generate deterministic cryptographically strong random number sequences
// from known seeds. It's basically just AES-CTR. Its design should be considered a protocol constant.
type seededPrng struct {
	i uint
	c cipher.Block
	b [16]byte
	n [16]byte
}

func (s *seededPrng) seed(b []byte) {
	k := sha256.Sum256(b)
	k[8]++ // defensive precaution in case the same 'key' is used elsewhere to initialize AES from SHA256
	s.i = 16
	s.c, _ = aes.NewCipher(k[:])
	for i := range s.n {
		s.n[i] = 0
	}
}

func (s *seededPrng) Read(b []byte) (int, error) {
	for i := 0; i < len(b); i++ {
		if s.i == 16 {
			s.i = 0
			for j := 0; j < 16; j++ {
				s.n[j]++
				if s.n[j] != 0 {
					break
				}
			}
			s.c.Encrypt(s.b[:], s.n[:])
		}
		b[i] = s.b[s.i]
		s.i++
	}
	return len(b), nil
}
