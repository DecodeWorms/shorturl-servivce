package idgenerator

import (
	"github.com/oklog/ulid"
	"io"
	"math/rand"
	"strings"
	"sync"
	"time"
)

//go:generate mockgen -source=id_generator.go -destination=../../mocks/id_generator_mock.go -package=mocks
type IdGenerator interface {
	Generate() string
}

type ulIdGenerator struct {
	entropy io.Reader
	mutex   sync.Mutex
}

func New() IdGenerator {
	// Create a properly seeded random source
	source := rand.NewSource(time.Now().UnixNano())
	randomizer := rand.New(source)

	// Create the monotonic entropy source
	entropy := ulid.Monotonic(randomizer, 0)

	return &ulIdGenerator{
		entropy: entropy,
	}
}

func (generator *ulIdGenerator) Generate() string {
	generator.mutex.Lock()
	defer generator.mutex.Unlock()

	// Get current timestamp and generate ULID
	timestamp := ulid.Timestamp(time.Now())
	id := ulid.MustNew(timestamp, generator.entropy)

	return strings.ToLower(id.String())
}
