# rckv

Simple key-value database implemented as an HTTP server.

## Next steps

- Add append-only disk-based storage
- Add compaction
- Update storage format to something that allows more efficient lookups, e.g.
  SSTables
- Increase write throughput by reworking design to reduce write lock contention

## Other ideas

- Signal to the client whether a set operation updated or created a record
- Allow setting multiple keys at once
- Add deletes
