# What?

Elect leader. Once it's elected as the leader, it remains being leader unless network faults happens.

# Why?

There are some existing libraries for leader election but they seem to be slow to react to change, and are in favor of resigning after job done. I find it the same as a distributed lock. In that case, time not synced will be a problem.

# How to use

##### ENV:

envs are declared in [env.go](env.go)

##### Code:

```
//
```

# Run

```
chmod 765 test.sh ./examples/run.sh
```

Run examples

```
./examples/run.sh
```

Run tests

```
./test.sh
```
