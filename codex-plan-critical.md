# go-common-types — critical defect evidence and authorized fixes

Base: `1a139b52726a8bcbf02b06527f3fed684d5994d4`. Each of the seven critical findings has a probe in
`/Volumes/Dev/Dev/critical-review-20261116/go-common-types/probes/pN/main.go`, and its raw output is in
`.../logs/probe-pN*.log`. All were reproduced on Go 1.27.4. The original audit prohibited fixes; the follow-up request
authorizes the seven bounded fixes below. Source line references in the evidence sections refer
to the original base, not the modified files.

## Implementation checkpoint

| Finding | Bounded implementation | Shipped regression |
|---|---|---|
| CF-A | Snapshot before callbacks/formatters, decode outside locks, single atomic commit; actual BSON codec | CallbackReentrySnapshot, DecoderCallbackReentry, ReadersAndWriterProgress, JSONMarshalerReentry, BSONWireRoundTrip, BSONWriterProgress |
| CF-B | Snapshot each binary input independently; explicit self identities, no two shared locks | BinaryOperationsConcurrent, NestedReadersSelfAndBinaryProgress |
| CF-C | Ordered SafeSet values under a single shared storage mutex; index bounds and mutation in one critical section | ConcurrentIndexedOperations |
| CF-D | Copy the requested tail while holding one read lock | LastNIndependent |
| CF-E | Decode fresh state and replace data/order together; failed decode preserves state | JSONReplacementAtomicity; corrected original UnmarshalJSON test |
| CF-F | Go equality for comparable values (including pointers), DeepEqual otherwise; no formatting identity | EqualityIgnoresFormatting, JSONConcurrentReplacement |
| CF-G | Validate 0..1/non-NaN before sorting, retain fraction interpolation | PercentileRejectsInvalidFraction |

The related BSON wire-format defect is resolved as required by CF-A: the existing methods use
MongoDB Go driver v2.10.0's BSON codec, preserve typed values/order and reject malformed framing.
No service connection is made. SafeSet keeps its public `T any` API and value receivers, with
shared private storage; membership/insertion are O(n), explicitly documented. Shared.GenerateHash
is unchanged. Existing slice Difference direction and ToSlice borrowed-view behavior are retained.

Negative controls run the final regression tests against archived original production source in
an isolated scratch fixture; every finding fails for its recorded mechanism. The original weak
ordered-map JSON test now asserts expected values and rejects a stale entry. Fix verification,
committed-HEAD `make ci`, review checkpoint and pending patch release are recorded in the external
`fix-report.json`/`fix-report.md`; no push or release happens before independent parent review.
The detailed sections below preserve the original failure evidence and bounded regression plan.

The packages' READMEs state the collections are "thread-safe" / "safe concurrent access". Severity
is measured against that contract.

## CF-A: a held non-reentrant RWMutex is re-acquired, causing a permanent deadlock
- **Severity:** critical. UnmarshalBSON hangs on every call. Its write lock is never released, so
  every later call on that map blocks forever. The recursive RLock variant hangs every goroutine
  using the collection once a concurrent writer is queued.
- **Source:**
  - `safeorderedmap/safeorderedmap.go:532-537` (UnmarshalBSON Lock) calls UnmarshalJSON, which
    takes the Lock again at `:499-501`.
  - `:524-529` (MarshalBSON RLock) calls MarshalJSON, which takes the RLock again at `:485-487`.
  - `safeset/safeset.go` holds the RLock while calling RLocking methods: String `:28-37`, ToSlice
    `:50-58`, First `:94-102`, Last `:106-114`, Clone `:146-152`, Map `:175-181`, and MarshalJSON
    `:339-345` (triple nesting).
  - `safeslice/safeslice.go`: LastN `:106-114` (ToSlice) and Mode `:445-466` (Empty, Unique).
  - The self-aliased calls `x.Difference(x)` `:375-382`, `x.Subset(x)` `:391-396` and
    `x.Intersection(x)` `:410-417`.
  - Callbacks also run under collection read locks: ordered-map Each `:276-285`, slice Each
    `:254-263` and set Map `:175-182` are representative entry points. A callback that mutates
    that collection tries to acquire its write lock while retaining the read lock.
- **Triggers:**
  - `safeorderedmap.New[int]().UnmarshalBSON([]byte(`{"key":1}`))` never returns.
    This original probe uses JSON text and establishes the input-independent lock failure only;
    it is not evidence of a valid BSON round trip.
  - The follow-up real-driver call `bson.Unmarshal(validBSON, safeorderedmap.New[int32]())`
    also blocks, using the validated 14-byte document `0e000000106b6579000100000000`.
  - `safeset.New(1,2,3)` with one goroutine looping `Add`, while others loop `String()`,
    `json.Marshal(s)` or `ToSlice()`, deadlocks.
  - `s := safeslice.New(1); s.Each(func(int) { s.Add(2) })` cannot return: Each's deferred
    RUnlock only runs after the callback returns, while Add waits for that read lock to release.
    This is a control-flow proof; the analogous ordered-map Each path explicitly documents
    callbacks performing modifications. No separate critical finding is counted.
- **Evidence:**
  - `probe-p1.log`: UnmarshalBSON did not return within 2s, and a later `Get` from another
    goroutine was also blocked.
  - `logs/bson-wire-contract.log` and `bson-contract/main.go` in the same external scratch:
    a driver-generated document passes `bson.Raw.Validate`, but unmarshaling it into this map
    blocks beyond 300ms. Source proves the second write Lock precedes decoding for every input.
  - `probe-p3.log`: no progress after about 1s. The stacks show the reader holding the RLock at
    `safeset.go:52` and waiting on `Size()` (`safeorderedmap.go:165` via `safeset.go:55`), while the
    writer waits in `Lock` (`safeorderedmap.go:43` via `safeset.go:68`).
- **Root cause:** public methods that take the lock call other public methods that take the same
  `sync.RWMutex`. Go's RWMutex does not support recursive locking: a queued writer blocks new
  readers.
- **Bounded fix:**
  - Add unexported lock-free helpers (`valuesLocked`, `sizeLocked`, `unmarshalJSONLocked`), and
    call only those while the lock is held.
  - BSON entry points must acquire their mutex exactly once. Removing an outer lock addresses
    the recursive-lock failure only; delegating to the JSON methods does NOT implement BSON.
    Resolving BSON support also requires the separate wire-format work below (SL-BSON1).
  - Binary operations should short-circuit `other == s` before taking a second lock.
  - Invoke caller callbacks over a snapshot after releasing internal locks; define and test
    snapshot iteration semantics so reentrant mutations do not deadlock or alter traversal.
  - No public API changes.
- **Regression and negative control:**
  - Generate actual BSON bytes with an independent BSON driver, validate the bytes, and call
    both direct UnmarshalBSON and driver `bson.Unmarshal(bytes, map)` under a 2s deadline. Use
    the 14-byte int32 document above plus the valid 5-byte empty BSON document; JSON text is
    not a substitute. Verify decoded contents as well as completion.
  - Separately validate MarshalBSON output with `bson.Raw.Validate` and decode it with the
    independent driver; require typed values/contents to round-trip through the public driver
    interfaces. A lock-only change is expected to fail these wire-format assertions until
    SL-BSON1 is resolved. Passing the watchdog alone cannot close BSON support.
  - Add a watchdog test (writer plus readers of String, MarshalJSON, ToSlice, First, Last, Clone,
    Map, LastN, Mode, and self-aliased binary operations) that fails on no progress. Run it with
    `-race`.
  - Negative control: put one nested call back, and confirm the watchdog test fails.

### Related remaining serialization defect SL-BSON1: BSON methods use JSON wire encoding

- **Classification:** confirmed correctness/serialization defect, separate from CF-A's critical
  deadlock and not counted as an eighth critical finding. The ordinary driver reports an
  encoding error; persistent data loss or a security consequence has not been established.
- **Claimed contract:** source comments at `safeorderedmap/safeorderedmap.go:523` and `:531`
  claim bson.Marshaler/bson.Unmarshaler. The cached independent MongoDB Go driver v1.17.10
  defines those interfaces in `bson/marshal.go:23-31` and `bson/unmarshal.go:16-27`; both require
  actual BSON documents. The [BSON specification](https://bsonspec.org/spec.html) requires a
  little-endian document length prefix and a trailing zero byte, unlike JSON text.
- **Source/root cause:** MarshalBSON `:524-528` returns MarshalJSON bytes unchanged;
  UnmarshalBSON `:532-536` delegates to UnmarshalJSON after the nested lock. These are separate
  wire-format and locking mechanisms.
- **Evidence:** external `bson-contract/main.go`, with raw output `logs/bson-wire-contract.log`:
  - Driver-generated `{key: int32(1)}` is valid 14-byte BSON:
    `0e000000106b6579000100000000`, and `bson.Raw.Validate` succeeds (positive control).
  - `New[int32]().Add("key",1).MarshalBSON()` returns nil error and nine JSON bytes
    `7b226b6579223a317d` (`{"key":1}`). The BSON validator reads a length of 1701519995 and
    rejects the nine-byte buffer. `bson.Marshal(map)` and independent BSON decoding also error.
  - Passing the valid BSON bytes to the existing UnmarshalJSON delegate errors on byte `0x0e`.
    This is an explicit control-flow projection of retaining that decoder after lock removal,
    not an implemented fix or a claim that current UnmarshalBSON reaches the decoder.
- **Bounded future work:** choose an actual BSON codec/representation that preserves the public
  methods' claimed document contract and generic values. Encode an ordinary snapshot, not the
  same receiver through its own BSON interface (which would recurse). Decode actual BSON into
  temporary state and commit under one lock, preserving ordered-map invariants and error
  atomicity. The authorized CF-A follow-up uses the v2.10.0 BSON codec for this work; no database
  service or unrelated client integration is added.
- **Regression/negative controls:** use independently generated valid BSON and validate output
  framing with `bson.Raw.Validate`, then exercise public driver marshal/unmarshal with typed
  contents, empty documents and supported nested/binary values. Reject malformed/truncated
  BSON without corrupting state or hanging. Restoring JSON-only encoding must fail the BSON
  validator/driver tests; removing nested locks alone must still fail serialization assertions.
  These assertions are implemented by the authorized CF-A follow-up; lock removal alone
  is not treated as BSON correctness.

## CF-B: binary set operations read `other` without its lock, causing a fatal concurrent map access
- **Severity:** critical. The runtime throws `fatal error: concurrent map read and map write`,
  which cannot be recovered and crashes the whole process.
- **Source:**
  - `safeorderedmap/safeorderedmap.go`: Union `:417-420`, Difference `:434-435`, Subset `:449-450`,
    Intersection `:471-472`, and Superset `:459-460`, which reads the receiver unlocked through
    `other.Subset(m)`.
  - `safeslice/safeslice.go`: Union `:364` reads `other.data` without a lock.
- **Trigger:** `a.Union(b)` or `a.Intersection(b)` in a loop while another goroutine runs
  `b.Add(k, v)` and `b.Delete(k)`.
- **Evidence:**
  - `probe-p2-norace.log` (built without `-race`): `fatal error: concurrent map read and map write`
    at `safeorderedmap.go:472`, with a deferred recover present and still bypassed.
  - The earlier race-build probe `codex-map-race-probe.log` reported two DATA RACE reports at
    `:419` against `:50`.
- **Root cause:** only the receiver's lock is held; `other.data` and `other.order` are read raw.
- **Bounded fix:**
  - Take a snapshot of `other` under `other.RLock()` (keys and values copied), release that lock,
    then lock `m`. This avoids lock-order inversion.
  - Self-aliasing is handled by the CF-A short-circuit.
- **Regression and negative control:**
  - Add a `-race` test with a concurrent `other.Add`/`Delete` alongside each binary operation.
  - Negative control: remove the snapshot, and confirm the race detector fails.

## CF-C: SafeSet.Get and Delete check, then act, across separate locks, causing an index-out-of-range panic
- **Severity:** critical. Normal concurrent use (`Get` alongside `Delete`) crashes the process.
  `Delete(i)` can also remove a different element than the caller indexed.
- **Source:** `safeset/safeset.go:74-80` (Get) and `:83-89` (Delete). `Size()` and `Values()` are
  two separate lock acquisitions.
- **Trigger:** one goroutine loops `s.Add(i%8); s.Delete(0)`, while another loops
  `n := s.Size(); s.Get(n-1)`.
- **Evidence:** `probe-p4.log`: `panic: runtime error: index out of range [0] with length 0` at
  `safeset.go:79`.
- **Root cause:** a time-of-check/time-of-use gap between the two lock acquisitions.
- **Bounded fix:**
  - Implement Get with `s.data.GetByIndex(index)`, which is atomic.
  - Add an atomic index-based delete primitive on the ordered map, done under one write lock
    (resolve the key from `order[i]`, then delete). SafeSet is a different package, so an
    unexported map method is inaccessible: use a backward-compatible exported primitive or
    refactor ownership without manually holding a lock while calling another locking method.
- **Regression and negative control:**
  - Add a concurrent Get/Delete churn test (no panic, run with `-race`).
  - Negative control: restore the two-step version, and confirm a panic.

## CF-D: SafeSlice.LastN's "new slice" shares storage with the original, so it is silently corrupted
- **Severity:** critical (data integrity). The returned SafeSlice changes when the original
  changes. The two instances also have different mutexes, so concurrent use is a data race.
- **Source:** `safeslice/safeslice.go:114-121`. `New[T](x...)` at `:514-517` adopts `x`, a
  subslice of `s.data`. The in-place `Delete` at `:68` shifts the shared array.
- **Trigger:** `s := New[int](); s.Add(1).Add(2).Add(3); t := s.LastN(2); s.Delete(0)`. After
  this, `t` is `[3 3]`.
- **Evidence:** `probe-p5.log`: `before: [2 3], after: [3 3]`.
- **Root cause:** a variadic slice is adopted without copying.
- **Bounded fix:**
  - Copy in LastN with `append([]T(nil), sl[len(sl)-n:]...)`.
  - Preserve existing result conventions for empty/zero requests. Negative-count validation is
    a separate compatibility decision, outside the proven shared-storage fix.
- **Regression and negative control:**
  - Add a test that mutates the original after LastN and asserts the tail is unchanged.
  - Negative control: remove the copy, and confirm the test fails.

## CF-E: SafeOrderedMap.UnmarshalJSON on a reused map keeps stale entries, so later Adds are lost
- **Severity:** critical (data integrity, silent loss). After decoding into a non-empty map, a
  stale key stays in `data` but leaves `order`. A later `Add(staleKey, v)` is then skipped by
  Keys, Values, Size, Each and MarshalJSON, so re-serialization drops it.
- **Source:** `safeorderedmap/safeorderedmap.go:508-518` resets `order` but not `data`. Add at
  `:46-48` only appends to `order` for keys absent from `data`.
- **Trigger:**
  - `m := New[int]().Add("x",1); json.Unmarshal([]byte(`{"y":2}`), m); m.Add("x",3)`.
  - `json.Marshal(m)` then yields `{"y":2}`.
- **Evidence:** `probe-p6.log`: `Get(x)=3,true Size=1 Keys=[y] Marshal={"y":2}`.
- **Root cause:** the decode reuses the old `data` map.
- **Bounded fix:** in UnmarshalJSON, replace `m.data` with a fresh map built from `temp`.
- **Regression and negative control:**
  - Add a P6-shaped test asserting `Keys`, `Size` and the marshal output contain `x`.
  - Negative control: restore map reuse, and confirm the test fails.

## CF-F: SafeSet element identity is sha256(fmt "%v"), so distinct values are merged and Contains gives false positives
- **Severity:** critical (data integrity, silent loss). The README promises set semantics for
  comparable values, but `==`-distinct values whose `%v` text is equal collapse into one element.
- **Source:** `shared/shared.go:10-13`, used by Add (`safeset.go:67-68`), Contains (`:129-130`)
  and Delete (`:88`).
- **Triggers:**
  - `safeset.New(Rec{"Mary Ann","Lee"}, Rec{"Mary","Ann Lee"}).Size() == 1`.
  - `safeset.New[any](1, "1")` keeps only `"1"`.
- **Evidence:** `probe-p8.log`. The negative control `New(1,2).Size()==2` passes.
- **Root cause:** `%v` is not an injective encoding (unquoted strings, no type tag).
- **Bounded fix:**
  - Define intended equality first, including pointer identity, dynamic interface types and
    non-comparable values accepted by the existing `T any` API. Preserve that public type
    constraint; any proposed constraint restriction requires a separate migration decision.
  - Enforce equality/collision resolution when inserting and looking up values; a formatted
    hash may select a bucket but must never be the sole identity check. Do not trust String,
    GoString or other formatting hooks to establish equality.
  - Keep this correction internal to SafeSet unless a compatibility review explicitly permits
    changing the separately exported shared.GenerateHash result contract.
- **Regression and negative control:**
  - Add a P8-shaped test (two distinct Recs give Size 2; `1` and `"1"` give Size 2), plus two
    distinct values of one GoStringer type whose GoString methods both return the same text.
    They collide even with `%T|%#v`; both values must survive and Contains/Delete must respect
    the defined equality. Include true duplicates and pointer/non-comparable cases.
  - Negative control: restore formatting-only identity (including `%T|%#v`) and confirm the
    corresponding collision/bypass tests fail.

## CF-G: statistical.Percentile panics on its documented 0..100 input
- **Severity:** critical (crash on documented use). The doc comment and the README say
  "percentage (between 0 and 100)", but the code treats `p` as a fraction. Documented values
  such as `p=90` with five inputs index out of range. Not every `p > 1` panics: with five inputs,
  `p=1.01` reaches the last-index early return instead.
- **Source:** `statistical/statistical.go:113-136` (rank `= (n-1)*p` at `:124`, then `s[idx]` at
  `:134`). README table row "Percentile".
- **Trigger:** `statistical.Percentile([]int{1,2,3,4,5}, 90)`.
- **Evidence:** `probe-p7.log`: `index out of range [360] with length 5`.
- **Root cause:** the doc and the code disagree, and `p` is not validated (negative values and
  NaN also index out of range).
- **Bounded fix:**
  - Keep the fraction semantics, which the existing test and Example use (`0.5`), so current
    callers are not broken.
  - Return an error for `p < 0`, `p > 1` or NaN.
  - Correct the doc comment and the README to say 0..1.
- **Regression and negative control:**
  - Add table tests: `p` of 90, -0.1 and NaN return an error; 0 returns the min; 1 returns the max.
  - Negative control: remove validation, and confirm the panic returns.
