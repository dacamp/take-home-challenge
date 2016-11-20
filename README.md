Challenge
=========

# Problem Space
A highly available distributed database of grow-only counters [[1]](https://gist.github.com/sargun/f7300f622181c543349a59e02c166391).

This is a timed challenge, albeit self-governed.  I began researching the problem space on 07:00 PST Friday Nov 18th and began writing the README around 19:30 on Saturday Nov 19th.  Code writing began around 21:00.

# Priorities
As accuracy and availability are top priorities, while factoring in the commutative nature of an grow-only unsigned integer dataset, this package intends to focus on a strong eventual consistency CRDT model.  Achieving this requires a comprise to consistency, satisfying CAP theorem.

The remaining two priorities focus on reliability and resilency and can be iterated upon during the lifetype of this project given a solid initial design.  I'd prefer to focus on metrics and tracing in the intial commits of this service, but I've decided to minimize the number of external packages to better present my skills as a software engineer.  I'll make notes of these comprimises where appropriate, regardless of implementation.

* [P0] To be accurate for the consistent endpoint. That endpoint should always return values that are accurate given a point in time after all messages have finished processing in the system.
* [P1] To be available in light of network or process failure.
* [P2] To be crash-tolerant.

## Accuracy
The consistent endpoint `GET /counter/:name:/consistent_value` must be accurate.

> It may be inconsistent, and therefore there is another endpoint, /counter/:name:/consistent_value. This endpoint is allowed to block for a period of time, and we will not send any requests until it responds. If network conditions prohibit it, you may return an error. The expectation is that after we heal the network, you will respond with the current value of the counter within a reasonable amount of time.

I've understood this to mean that requests to this endpoint should communicate with all known peers to validate each node's value.   I plan to cycle through the list of known peers, make an RPC call to retreive the value of `:name:` and then to merge them into an accurate presentation of the value of `:name:`.  As mentioned, no additional requests will be sent until this endpoint responds, so inflight writes that have not properly been replicated on all nodes maybe ignored.

> If network conditions prohibit it, you may return an error.

**NOTE:** given the unreliability of networks, a default timeout (for all requests) will be 5 seconds.  I'd like this value to be flexible.  The `CHALLENGE_TIMEOUT=5` environment variable has been exposed and accepts a positive integer value interpreted as seconds.

# Approach
A collection of peers (nodes) will each carry a replica of the shared dataset and communicate via a broadcast protocol.  The models must allow for idempotent requests, while handling any out of order or incomplete requests.

## Broadcast
Given the timeline, external package handicap, and my skill set, I'll opt for a rudamentry message broadcasting system.

Given _n_ peers, an incoming request to peer `n[0]` will be concurrently broadcast updates out to all peers `n[1,n-1]`.  This has a pretty large and blatent failure domain.  These limitations are self imposed, mostly due to my lack of experience using gossip based protocols.  I may deviate and try to write some type of per counter gossip.

Should the resources be available, I'd prefer to use a gossip protocol, such as the opensource  [hashicorp/memberlist](https://github.com/hashicorp/memberlist) project.


## Model
```go
type Counter interface {
	ID() string # unique string identifer

    ConsistentValue() (uint64, error)
	IncVal(uint64)
	Merge(uint64)
	Value() uint64
}

type gCounter struct {
	id      string
	counter uint64

	peers *PeerList
}

func (g *gCounter) ConsistentValue() (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CHALLENGE_TIMEOUT * time.Second)

	// Even though ctx should have expired already, it is good
	// practice to call its cancelation function in any case.
	// Failure to do so may keep the context and its parent alive
	// longer than necessary.
	defer cancel()

	// probably a channel that closes once all peers have been heard
	// from or a timeout occurs
	return peers.GetPeerValues(ctx, g)
}

func (g *gCounter) IncVal(i uint64) {
	peers.Broadcast(atomic.AddInt64(&g.counter, i))
}

// Merge will only assign the new value if it's greater than the current counter value
func (g *gCounter) Merge(v uint64) {
	if g.counter < v {
		atomic.StoreUint64(&g.counter, v)
	}
}

func (g *gCounter) Value() uint64 {
	return g.counter // no need for this to be an atomic operation
}
```


# References
* https://arxiv.org/pdf/0907.0929v1.pdf
* http://hal.upmc.fr/inria-00555588/document
* https://medium.com/@istanbul_techie/a-look-at-conflict-free-replicated-data-types-crdt-221a5f629e7e
