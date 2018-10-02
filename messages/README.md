# Messages
Persistent packet storage

# Message order preservation

#### Description: 
 Preserve the original ordering in which messages were published to the wait queue while keeping memory usage and lookup complexity to a minimum.

#### Use Case:
Suppose clientA publishes a set of messages {message1, message2,  ...., messagen} to a particular topic T which is subscribed by clientB. When clientB is offline and has a subscription with QoS > 0, messages are stored for redelivery whenever clientB becomes online. Due to the message sequence number ( Message ID ), messages can become out of order which violates consistency rule. Therefore it is crucial to have a mechanism to preserve message order, it should satisfy rule for two stages. The first stage is the partial ordering of unacknowledged messages, and the second stage is total order preservation of messages with the higher quality of service.
 
#### Change(s):
- Improve message storage by implementing it with Ordered Sets as opposed to the current implementation which uses Hash Sets.

# Test units
check to ensure that it passes all test cases and view the code coverage in `html` mode.

**NOTE:** `msgidfull` is the build tag for concurrent test case of `GetNewID(uuid.UUID)` which requires `65535` ( `0xFFFF` )insertions. It must return `0` when all `65535` slots are occupied.

```bash
$ go test -v -coverprofile=cover.out . && go tool cover -html=cover.out
$ go test -v -tags msgidfull -coverprofile=cover.out . && go tool cover -html=cover.out
```

