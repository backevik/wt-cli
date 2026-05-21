package namegen

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	mathrand "math/rand/v2"
	"strings"
)

// words is the combined word list from all four categories.
// Preserved exactly from the original worktree-functions.zsh (including duplicates
// to maintain the same probability distribution).
var words = []string{
	// Technology
	"algorithm", "binary", "cache", "compiler", "daemon", "debug", "encrypt", "firmware", "gateway",
	"hardware", "indexed", "kernel", "lambda", "memory", "network", "optimize", "parallel", "query",
	"runtime", "script", "terminal", "upload", "vector", "widget", "cipher", "decode", "filter",
	"handler", "matrix", "packet", "router", "signal", "thread", "zenith", "aurora", "blossom",
	"coral", "delta", "ember", "forest", "glacier", "harbor", "island", "jungle", "canyon", "meadow",
	// Nature
	"nebula", "ocean", "prairie", "quartz", "river", "summit", "tundra", "valley", "willow", "breeze",
	"eclipse", "jungle", "crystal", "dune", "ember", "fjord", "grove", "haven", "inlet", "canyon",
	"lagoon", "mesa", "oasis", "peak", "ridge", "shore", "tide", "wave", "crest", "foam",
	"glade", "heath", "islet", "knoll", "marsh", "plain", "reef", "shelf", "slope", "vista",
	// Actions
	"buzzing", "climbing", "dancing", "fluttering", "gliding", "hovering", "imagining", "juggling",
	"knitting", "laughing", "musing", "navigating", "orbiting", "painting", "querying", "rippling",
	"singing", "tinkering", "unfolding", "vibrating", "wandering", "exploring", "yearning", "zooming",
	"ascending", "bouncing", "cascading", "drifting", "echoing", "floating", "galloping", "humming",
	"iterating", "jolting", "kindling", "leaping", "mending", "nurturing", "observing", "perching",
	// Adjectives
	"abundant", "accurate", "agile", "bright", "calm", "clever", "cozy", "crisp", "dynamic", "elegant",
	"fluffy", "gentle", "happy", "infinite", "jolly", "kinetic", "logical", "mystic", "nimble", "optimal",
	"parsed", "peaceful", "quantum", "radiant", "serene", "thermal", "unified", "virtual", "warm", "wise",
	"zealous", "bold", "daring", "eager", "friendly", "graceful", "honest", "innovative", "joyful", "keen",
	"lively", "mellow", "noble", "optimistic", "playful", "quick", "resilient", "sturdy", "tranquil", "upbeat",
}

func newRand() *mathrand.Rand {
	var seed [32]byte
	if _, err := cryptorand.Read(seed[:]); err != nil {
		// fallback: derive seed from two crypto/rand uint64s
		var b [8]byte
		_, _ = cryptorand.Read(b[:])
		s := binary.LittleEndian.Uint64(b[:])
		return mathrand.New(mathrand.NewPCG(s, s^0xdeadbeefcafe1234))
	}
	return mathrand.New(mathrand.NewChaCha8(seed))
}

// RandomBranchName generates a branch name in the format "{username}/{word1}-{word2}-{word3}".
func RandomBranchName(username string) string {
	r := newRand()
	w1 := words[r.IntN(len(words))]
	w2 := words[r.IntN(len(words))]
	w3 := words[r.IntN(len(words))]
	return username + "/" + strings.Join([]string{w1, w2, w3}, "-")
}
