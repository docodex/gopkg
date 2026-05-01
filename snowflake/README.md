# snowflake

snowflake provides a very simple Twitter snowflake generator.

## ID Format

By default, the ID format follows the original Twitter snowflake format.

* The ID as a whole is a 63 bit integer stored in an int64
* 42 bits are used to store a timestamp with millisecond precision, using a custom epoch(about 139 years from the epoch time).
* 10 bits are used to store a node id - a range from 0 through 1023.
* 11 bits are used to store a sequence number - a range from 0 through 2047.

## Custom Format

The number of bits used for the node id and the sequence number are compile-time constants
(`nodeBits` and `sequenceBits` in `snowflake.go`). Changing the layout requires editing those
constants in the source and rebuilding; there is no runtime option for it.

Remember that there is a maximum of (63 - timeBits) bits available that can be shared between these two values.

## Custom Epoch

By default, this package uses the epoch of "2026-01-01 00:00:00 +0000 UTC". You can set your own epoch value by setting `epochTimestamp`.

## Custom Time Unit

By default, this package uses the time unit of 1 ms. You can set your own time unit by setting `timeUnit`.

## How it Works

Each time you generate an ID, it works, like this.

* A timestamp with `timeUnit` millisecond(s) precision is stored using `timeBits` bits of the ID.
* Then the node id is added in subsequent bits.
* Then the sequence number is added, starting at 0 and incrementing for each ID generated in the same millisecond. If you generate enough IDs in the same millisecond that the sequence would roll over, the caller blocks inside `Generate()` (busy-waiting on the monotonic clock) until the next millisecond.

The default Twitter format shown below.

```
+---------------------------------------------------------------------------+
| 1 Bit Unused | 42 Bit Timestamp | 10 Bit Node ID | 11 Bit Sequence Number |
+---------------------------------------------------------------------------+
```

Using the default settings, this allows for 2048 unique IDs to be generated every millisecond, per node id.

## Usage

Import the package into your project then construct a new `Snowflake` instance with
a unique node id. The default settings permit a node id range from 0 to 1023. If you
have set a custom `nodeBits` value, you will need to calculate what your node id
range will be. Call the `Generate()` method on the instance to generate and return a
unique snowflake ID.

Keep in mind that every `Snowflake` instance - including instances running on
different servers - must be configured with a unique node id. If two instances share
the same node id the generator cannot guarantee unique IDs across them.

**Example Program:**

```
package main

import (
	"fmt"

	"github.com/docodex/gopkg/snowflake"
)

func main() {
	// Create a new Snowflake instance with a node id of 1
	s, err := snowflake.New(snowflake.WithNode(1))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Generate a snowflake ID.
	id, err := s.Generate()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print out the ID.
	fmt.Printf("ID: %d\n", id)
	// Print out the ID's timestamp
	fmt.Printf("ID Time: %d\n", snowflake.Timestamp(id))
	// Print out the ID's node id
	fmt.Printf("ID Node: %d\n", snowflake.Node(id))
	// Print out the ID's sequence number
	fmt.Printf("ID Sequence: %d\n", snowflake.Sequence(id))
}
```

## Performance

With default settings, this snowflake generator should be sufficiently fast
enough on most systems to generate 2048 unique ID's per millisecond. This is the
maximum that the snowflake ID format supports. That is, around 488-489
nanoseconds per operation. While set the sequence to 12, that is around
243-244 nanoseconds per operation. That would be around 31-32 nanoseconds
per operation if set the sequence bits to 26.

Since the snowflake generator is single threaded the primary limitation will be
the maximum speed of a single processor on your system.

To benchmark the generator on your system run the following command inside the
snowflake package directory.

```sh
go test -run=^$ -bench=.
```
