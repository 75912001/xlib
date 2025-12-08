package hashring

import "crypto/md5"

func hashDigest(data string) [md5.Size]byte {
	return md5.Sum([]byte(data))
}

func hashVal(bKey []byte) uint32 {
	return (uint32(bKey[3]) << 24) |
		(uint32(bKey[2]) << 16) |
		(uint32(bKey[1]) << 8) |
		(uint32(bKey[0]))
}
