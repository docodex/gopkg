package snowflake_test

import (
	"errors"
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/docodex/gopkg/internal"
	"github.com/docodex/gopkg/snowflake/v2"
)

//******************************************************************************
// General Test funcs

func TestDefault(t *testing.T) {
	s := snowflake.Default()
	for range 10 {
		id, err := s.Generate()
		if err != nil {
			t.Fatalf("failed to generate id: %v", err)
		}
		fmt.Println(id)
	}
}

func TestNew(t *testing.T) {
	_, err := snowflake.New(snowflake.WithNode(0))
	if err != nil {
		t.Fatalf("error creating snowflake, %s", err)
	}
	_, err = snowflake.New(snowflake.WithNode(5000))
	if err == nil {
		t.Fatalf("no error creating snowflake, %s", err)
	}
}

// lazy check if Generate will create duplicate IDs
// would be good to later enhance this with more smarts
func TestGenerateDuplicateID(t *testing.T) {
	s, _ := snowflake.New(snowflake.WithNode(1))
	var x, y int64
	for range 1000000 {
		y, _ = s.Generate()
		if x == y {
			t.Errorf("x(%d) & y(%d) are the same", x, y)
		}
		x = y
	}
}

func TestPrintAll(t *testing.T) {
	s, err := snowflake.New(snowflake.WithNode(0))
	if err != nil {
		t.Fatalf("error creating snowflake, %s", err)
	}
	id, _ := s.Generate()
	t.Logf("Int64    : %#v", s.Decompose(id))
}

func TestGenerate(t *testing.T) {
	now := time.Now()
	s, err := snowflake.New(snowflake.WithEpoch(now))
	if err != nil {
		t.Fatalf("failed to create snowflake: %v", err)
	}

	sleepTime := int64(50)
	time.Sleep(time.Millisecond * time.Duration(sleepTime))

	id, err := s.Generate()
	if err != nil {
		t.Fatalf("failed to generate id: %v", err)
	}

	actualTime := (s.Timestamp(id) - now.UnixMilli())
	if actualTime < sleepTime || actualTime > sleepTime+1 {
		t.Errorf("unexpected time: %d", actualTime)
	}

	actualSequence := s.Sequence(id)
	if actualSequence != 0 {
		t.Errorf("unexpected sequence: %d", actualSequence)
	}

	actualNode := s.Node(id)
	if actualNode != int64(internal.Lower8BitPrivateIPv4()) {
		t.Errorf("unexpected machine: %d", actualNode)
	}

	fmt.Println("snowflake id:", id)
	fmt.Println("epoch time:", now.UnixMilli())
	fmt.Println("decompose:", s.Decompose(id))
}

func TestGenerate_InSequence(t *testing.T) {
	now := time.Now()
	s, err := snowflake.New(snowflake.WithEpoch(now))
	if err != nil {
		t.Fatalf("failed to create snowflake: %v", err)
	}

	startTime := now.UnixMilli()
	node := int64(internal.Lower8BitPrivateIPv4())

	var numID int
	var lastID int64
	var maxSeq int64

	currentTime := startTime
	for currentTime-startTime < 200 {
		id, err := s.Generate()
		if err != nil {
			t.Fatalf("failed to generate id: %v", err)
		}
		currentTime = time.Now().UnixMilli()
		numID++

		if id == lastID {
			t.Fatal("duplicated id")
		}
		if id < lastID {
			t.Fatal("must increase with time")
		}
		lastID = id

		parts := s.Decompose(id)

		actualTime := parts["time"]
		overtime := startTime + actualTime - currentTime
		if overtime > 0 {
			t.Errorf("unexpected overtime: %d", overtime)
		}

		actualSequence := parts["sequence"]
		if actualSequence > maxSeq {
			maxSeq = actualSequence
		}

		actualMachine := parts["node"]
		if actualMachine != node {
			t.Errorf("unexpected machine: %d", actualMachine)
		}
	}

	if maxSeq != (1<<11)-1 {
		t.Errorf("unexpected max sequence: %d", maxSeq)
	}
	fmt.Println("max sequence:", maxSeq)
	fmt.Println("number of id:", numID)
}

func TestGenerate_InParallel(t *testing.T) {
	s1, err := snowflake.New(snowflake.WithNode(1))
	if err != nil {
		t.Fatalf("error creating snowflake, %s", err)
	}
	s2, err := snowflake.New(snowflake.WithNode(2))
	if err != nil {
		t.Fatalf("error creating snowflake, %s", err)
	}

	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	fmt.Println("number of cpu:", numCPU)

	consumer := make(chan int64)

	const numID = 1000
	generate := func(s *snowflake.Snowflake) {
		for range numID {
			id, err := s.Generate()
			if err != nil {
				t.Logf("failed to generate id: %v", err)
				return
			}
			consumer <- id
		}
	}

	var numGenerator int
	for range numCPU / 2 {
		go generate(s1)
		go generate(s2)
		numGenerator += 2
	}

	set := make(map[int64]struct{})
	for range numID * numGenerator {
		id := <-consumer
		if _, ok := set[id]; ok {
			t.Fatal("duplicated id")
		}
		set[id] = struct{}{}
	}
	fmt.Println("number of id:", len(set))
}

func TestGenerate_ReturnsError(t *testing.T) {
	s, err := snowflake.New(
		snowflake.WithEpoch(time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)),
		snowflake.WithNode(1),
	)
	if err != nil {
		t.Fatalf("failed to create snowflake: %v", err)
	}
	_, err = s.Generate()
	if err != nil {
		t.Logf("failed generate id: %v", err)
	} else {
		t.Fatalf("generate id err: %v", err)
	}
}

func TestComposeAndDecompose(t *testing.T) {
	now := time.Now()
	s, err := snowflake.New(snowflake.WithEpoch(now))
	if err != nil {
		t.Fatalf("failed to create snowflake: %v", err)
	}

	testCases := []struct {
		name     string
		time     time.Time
		node     int64
		sequence int64
	}{
		{
			name:     "zero values",
			time:     now,
			sequence: 0,
			node:     0,
		},
		{
			name:     "max sequence",
			time:     now,
			node:     0,
			sequence: 1<<11 - 1,
		},
		{
			name:     "max machine id",
			time:     now,
			node:     1<<10 - 1,
			sequence: 0,
		},
		{
			name:     "future time",
			time:     now.Add(time.Hour),
			sequence: 0,
			node:     0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id, err := s.Compose(tc.time, tc.node, tc.sequence)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			parts := s.Decompose(id)

			// Verify time part
			expectedTime := tc.time.UnixMilli()
			if parts["timestamp"] != expectedTime {
				t.Errorf("time mismatch: got %d, want %d", parts["time"], expectedTime)
			}

			// Verify sequence part
			if parts["sequence"] != int64(tc.sequence) {
				t.Errorf("sequence mismatch: got %d, want %d", parts["sequence"], tc.sequence)
			}

			// Verify machine id part
			if parts["node"] != int64(tc.node) {
				t.Errorf("node id mismatch: got %d, want %d", parts["node"], tc.node)
			}

			// Verify id part
			if parts["id"] != id {
				t.Errorf("id mismatch: got %d, want %d", parts["id"], id)
			}
		})
	}
}

const year = time.Duration(365*24) * time.Hour

func TestCompose_ReturnsError(t *testing.T) {
	now := time.Now()
	s, err := snowflake.New(snowflake.WithEpoch(now))
	if err != nil {
		t.Fatalf("failed to create snowflake: %v", err)
	}

	testCases := []struct {
		name     string
		time     time.Time
		sequence int64
		node     int64
		err      error
	}{
		{
			name:     "start time ahead",
			time:     now.Add(-time.Second),
			node:     0,
			sequence: 0,
			err:      snowflake.ErrOverTimeLimit,
		},
		{
			name:     "over time limit",
			time:     now.Add(time.Duration(175) * year),
			node:     0,
			sequence: 0,
			err:      snowflake.ErrOverTimeLimit,
		},
		{
			name:     "invalid sequence",
			time:     now,
			node:     0,
			sequence: 1 << 11,
			err:      snowflake.ErrOverSequenceLimit,
		},
		{
			name:     "invalid machine id",
			time:     now,
			node:     1 << 10,
			sequence: 0,
			err:      snowflake.ErrOverNodeLimit,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := s.Compose(tc.time, tc.node, tc.sequence)
			if !errors.Is(err, tc.err) {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// ****************************************************************************
// Benchmark Methods

func BenchmarkGenerate(b *testing.B) {
	s, _ := snowflake.New(snowflake.WithNode(1))

	b.ReportAllocs()

	for b.Loop() {
		_, _ = s.Generate()
	}
}

func BenchmarkGenerate2(b *testing.B) {
	s, _ := snowflake.New(
		snowflake.WithNodeBits(9),
		snowflake.WithSequenceBits(12),
		snowflake.WithNode(1),
	)

	b.ReportAllocs()

	for b.Loop() {
		_, _ = s.Generate()
	}
}

func BenchmarkGenerateMaxSequence(b *testing.B) {
	s, err := snowflake.New(
		snowflake.WithNodeBits(1),
		snowflake.WithSequenceBits(26),
		snowflake.WithNode(1),
	)
	if err != nil {
		b.Fatalf("error creating snowflake, %s", err)
	}

	b.ReportAllocs()

	for b.Loop() {
		_, _ = s.Generate()
	}
}
