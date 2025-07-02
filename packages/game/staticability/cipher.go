package staticability

import "deckronomicon/packages/game/mtg"

type Cipher struct {
}

func (a *Cipher) Name() string {
	return "Cipher"
}

func (a *Cipher) StaticKeyword() mtg.StaticKeyword {
	return mtg.StaticKeywordCipher
}
