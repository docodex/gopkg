# snowflake

snowflake provides a very simple Twitter snowflake generator.

Version v3 offers support to customize `Snowflake.nodeBits` and `Snowflake.sequenceBits` without changing the source code. In particular, v3 offers support for `Snowflake.timeUnit`.

## ID Format

By default, the ID format follows the original Twitter snowflake format.

* The ID as a whole is a 63 bit integer stored in an int64
* 39 bits are used to store a timestamp in unit of 10 msec, using a custom epoch(about 170 years from the epoch time).
* 11 bits are used to store a node id - a range from 0 through 2047.
* 13 bits are used to store a sequence number - a range from 0 through 8191.

## Custom Format

You can alter the number of bits used for the node id and sequence number by setting the `Snowflake.nodeBits` and `Snowflake.sequenceBits` values.

Remember that there is a maximum of (63 - timeBits) bits available that can be shared between these two values.

## Custom Epoch

By default, this package uses the Epoch of "2025-01-01 00:00:00 +0000 UTC". You can set your own epoch value by setting `Snowflake.epoch` to use as the epoch.

## How it Works

Each time you generate an ID, it works, like this.

* A timestamp with `Snowflake.timeUnit` precision is stored using `Snowflake.timeBits` bits of the ID.
* Then the node id is added in subsequent bits.
* Then the sequence number is added, starting at 0 and incrementing for each ID generated in the same time unit. If you generate enough IDs in the same time unit that the sequence would roll over or overfill then the generate function will pause until the next time unit.

The default Twitter format shown below.

```
+---------------------------------------------------------------------------+
| 1 Bit Unused | 39 Bit Timestamp | 11 Bit Node ID | 13 Bit Sequence Number |
+---------------------------------------------------------------------------+
```

Using the default settings, this allows for 8192 unique IDs to be generated every 10 msec, per node id.

## Usage

Import the package into your project then construct a new snowflake Node using a
unique node id. The default settings permit a node id range from 0 to 2047.
If you have set a custom NodeBits value, you will need to calculate what your
node id range will be. With the node object call the `Generate()` method to
generate and return a unique snowflake ID.

Keep in mind that each node you create must have a unique node id, even
across multiple servers.  If you do not keep node ids unique the generator
cannot guarantee unique IDs across all nodes.

**Example Program:**

```
package main

import (
	"fmt"

	"github.com/docodex/gopkg/snowflake/v3"
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
	fmt.Println("ID: %d", id)
	// Print out the ID's timestamp
	fmt.Println("ID Time: %d", s.Timestamp(id))
	// Print out the ID's node id
	fmt.Println("ID Node: %d", s.Node(id))
	// Print out the ID's sequence number
	fmt.Println("ID Sequence: %d", s.Sequence(id))
}
```

## Performance

With default settings, this snowflake generator should be sufficiently fast
enough on most systems to generate 8192 unique ID's per 10 msec. This is the
maximum that the snowflake ID format supports. That is, around 1213-1214
nanoseconds per operation. While set the time unit to 1 msec, that is around
122-123 nanoseconds per operation. That would be around 31-32 nanoseconds
per operation if set the sequence bits to 26.

Since the snowflake generator is single threaded the primary limitation will be
the maximum speed of a single processor on your system.

To benchmark the generator on your system run the following command inside the
snowflake package directory.

```sh
go test -run=^$ -bench=.
```
