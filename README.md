# Timelock

This is a experimental implementation of a custom-designed time-lock encryption mechanism.

This is very much a work in progress, and more than likely faulty, as I am not a cryptographer.

`TODO: add details, design choices.`

`TODO: current implementation is very much a toy example with no real complexity. Needs to be updated.`

## Implementation notes

These are implementation notes stating mostly arbitrary choices as this is an arbitrary experiment. Actual choices have to be made still.

Functions:

- _AEAD_(plaintext, associated data, key, __nonce__)
- _sha256_(input)
- _`puzzle`_ indicates unknown "puzzle" value of certain length, with length contributing to difficulty to guess, i.e. puzzle complexity.
- `?input` = _sha256_(_`puzzle`_, input)

Given:

Each _interim`N`_ represents an iteration. Each iteration relies on the result of the previous iteration.

- __input1__
- __interim1__ = _AEAD_(`input2`, `?input1`, _sha256_(`?input1`), __nonce1__)
- __interim2__ = _AEAD_(`input3`, `?input2`, _sha256_(`?input2`), __nonce2__)
- __interim3__ = _AEAD_(`secretKey`, `?input3`, _sha256_(`?input3`), __nonce3__)
- __ciphertext__ = _AEAD_(__plaintext__, _NIL_, `secretKey`, __nonce4__)

Remarks:

- Solution complexity / time-bounds determined by chosen _one-way function_, _number of iterations of one-way function_, _number of iterations of interim values_.
- Massive parallelism mitigated by linearizing solving capability by creating many iterations (of interim values). Each iteration relies on the previous iteration's result for both _key material_ and _associative data_.
- Unhashed `?input` as _associative data_ to prevent skipping hashing challenge and tackling decryption directly. (If you know the _key_, you also know the _associative data_.)  
  > It is outside of the model how the associated-dataHis made known tothe receiver.  We do not consider the associated-data to be part of the ciphertext, though the receiver will need it in order to decrypt.

  Seems to suggest that indeed `AD` is a necessary component for decryption.
- Start _sha256_ hash-content with puzzle component, one cannot reuse "prehashed input". (puzzle component is the variable)
- `O(1)` for creator, as he can define puzzle components. `O(avg-hashing-time-to-find-puzzle-component)` for solver, as he needs to try out every byte-combination.
- Choices of one-way function, AEAD cipher, etc. are arbitrary choices atm based on what is readily available in Go standard library.
- ...

## TODO

- Write document explaining, detailing the mechanism.
  - look up actual terminology instead of own made-up stuff ...
- Fine-tune implementation to make it actually useable.
- Investigate if lengths of all values in use are acceptable in real-world scenarios. Are we crossing any boundaries/limits that violate cryptography rules?
- See if we can prove this: https://tamarin-prover.github.io/

## References

Material for on-the-fly learning of cryptographic concepts: (seriously, I'm mostly a noob with this stuff ...)

- https://en.wikipedia.org/wiki/Authenticated_encryption
- https://blog.cryptographyengineering.com/2012/05/19/how-to-choose-authenticated-encryption/
- https://web.cs.ucdavis.edu/~rogaway/papers/ad.pdf (needed to check up on security guarantee for "associated data", i.e. needed to confirm that "associated data" is a necessary component for decryption)

The following references were inspiration to this custom solution.

- https://www.gwern.net/Self-decrypting-files

References still to read:

- http://people.seas.harvard.edu/~salil/research/timelock.pdf
