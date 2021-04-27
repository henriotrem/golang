RAFT Protocol
====

I've been working couple of days following the excellent MIT 6.824 course about distributed systems, I got inspired by different implementation and incorporated my own improvements into the code - no use of sync.Mutext everything is done through channels and I tried to clean and refactor as much as I could the code. I also split the main.go into different files for readability and practicality.

An excellent website to understand the [protocol](https://raft.github.io/).

![alt text](https://miro.medium.com/max/1400/0*t1x9nNs6KIf7RSv4.png)

Raft is a consensus algorithm that is designed to be easy to understand. It's equivalent to Paxos in fault-tolerance and performance. The difference is that it's decomposed into relatively independent subproblems, and it cleanly addresses all major pieces needed for practical systems. We hope Raft will make consensus available to a wider audience, and that this wider audience will be able to develop a variety of higher quality consensus-based systems than are available today.

The use cases for such a protocol are far-reaching, such as replicated state
machines which are a key component of many distributed systems. They enable
building Consistent, Partition Tolerant (CP) systems, with limited
fault tolerance as well.

# MIT 6.824 Distributed Systems Labs

### (Updated to Spring 2021 Course Labs)

Course website: http://nil.csail.mit.edu/6.824/2020/schedule.html

- [x] Lab 1: MapReduce

- [x] Lab 2: Raft Consensus Algorithm
  - [x] Lab 2A: Raft Leader Election
  - [x] Lab 2B: Raft Log Entries Append
  - [x] Lab 2C: Raft state persistence
  
- [ ] Lab 3: Fault-tolerant Key/Value Service
  - [ ] Lab 3A: Key/value Service Without Log Compaction
  - [ ] Lab 3B: Key/value Service With Log Compaction

- [ ] Lab 4: Sharded Key/Value Service

