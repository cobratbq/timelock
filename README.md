# Timelock

This is a experimental implementation of a custom-designed time-lock encryption mechanism.

This is very much a work in progress, and more than likely faulty, as I am not a cryptographer.

`TODO: add details, design choices.`

`TODO: current implementation is very much a toy example with no real complexity. Needs to be updated.`

## Implementation notes

These are implementation notes stating mostly arbitrary choices as this is an arbitrary experiment. Actual choices have to be made still.

Functions:

- _AEAD_(plaintext, associated data, key, __nonce__), in this case AES-256-GCM.
- _sha256_(input)
- _'🧩'_ indicates unknown "puzzle" value of certain length, with length contributing to difficulty to guess, i.e. puzzle complexity.
- `🧩input` = _sha256_(_`🧩`_, input)

Given:

Each _milestone`N`_ represents an iteration in encrypted form. Each iteration relies on the result of the previous iteration. Any number of iterations can be used.

- __input1__
- __milestone1__ = _AEAD_(`input2`, `🧩input1`, _sha256_(`🧩input1`), __nonce1__)
- __milestone2__ = _AEAD_(`input3`, `🧩input2`, _sha256_(`🧩input2`), __nonce2__)
- __milestone3__ = _AEAD_(`secretKey`, `🧩input3`, _sha256_(`🧩input3`), __nonce3__)
- __ciphertext__ = _AEAD_(`plaintext`, _NIL_, `secretKey`, __nonce4__)

Remarks:

- Rationale:
  - `TODO: write down rationale`
  - `TODO: Is there any benefit to having the actual plaintext payload independent of the chain of iterations?`
- Time-component realized through assumed average guessing time for puzzle component (`n` number of bytes of random data).
- Solution complexity / time-bounds determined by chosen _one-way function_, _number of iterations of one-way function_, _number of iterations of milestone values_.
- Massive parallelism mitigated by serializing solving capability by creating many iterations (of milestone values). Each iteration relies on the previous iteration's result for both _key material_ and _associative data_.
- Trade-off: many milestones vs large puzzle-complexity for time-bounds. Many milestones of smaller complexity will give smaller upper and lower bounds for necessary time to solve? While number of milestones ensures certain amount of time needed to solve? (Does this even make sense?)
- As chained milestones each contain random byte-array being a key, it is hard to determine by statistical analysis whether or not decryption key is correctly guessed. Relies on _authentication tag_ for confirmation of correct decryption. _Authentication tag_ is obfuscated through _associated data_. _Associated data_ is _hash input_, hence you need to know both hashed value used for decryption and original pre-image used for hashing. So 2 necessary parts for decryption:
  - symmetric key, to decrypt content.
  - associated data, to confirm of decryption result.
- Unhashed `?input` as _associative data_:
  - to prevent skipping hashing challenge and tackling decryption directly. (If you know the _key_, you also know the _associative data_.)
    > It is outside of the model how the associated-data is made known to the receiver. We do not consider the associated-data to be part of the ciphertext, though the receiver will need it in order to decrypt.  
  -- [Phillip Rogaway][AEAD-paper], 20 September 2002
  
    Seems to suggest that indeed `AD` is a necessary component for decryption.
  - to prevent _2nd pre-image_ attacks, as the original hash input needs to be known.  
- Start _sha256_ hash-content with puzzle component, one cannot reuse "prehashed input". (puzzle component is the variable)
- `O(1)` for creator, as he can choose the puzzle components.  
  `O(avg-hashing-time-to-find-puzzle-component)` for solver, as he needs to rediscover the same puzzle components, e.g. try out every byte-combination.
- Choices of one-way function, AEAD cipher, etc. are arbitrary choices atm based on what is readily available in Go standard library.
- Individual milestone inputs can be released at will, in order to progress decryption effort.
- ...

## TODO

- Actual secret data is not associated with last milestone, is there any value in doing so?
- Theoretical questions:
  - Investigate if lengths of all values in use are acceptable in real-world scenarios. Are we crossing any boundaries/limits that violate cryptography rules?
  - Investigate whether using "hash input" as associated data might reveal any information on the "hash input". (To what extent is associated data recoverable from the ciphertext?)
  - How is this approach different from chained hashes? Each milestone can be solved with a single correctly guessed solution, so in theory can be solved very quickly. How is this different for chained hashes?
  - Confirm: AEAD is guaranteed to use AD to "obfuscate" authentication tag?
- Write document explaining, detailing the mechanism.
  - look up actual terminology instead of own made-up stuff ...
  - increasing complexity with each iteration.
  - different types of hardness by choice of one-way function (cpu-bound, memory-bound, ...) (As described in [article](https://www.gwern.net/Self-decrypting-files))
- Fine-tune implementation to make it actually useable.
- See if we can prove this: https://tamarin-prover.github.io/

## References

Material for on-the-fly learning of cryptographic concepts: (seriously, I'm mostly a noob with this stuff ...)

- https://en.wikipedia.org/wiki/Authenticated_encryption
- https://blog.cryptographyengineering.com/2012/05/19/how-to-choose-authenticated-encryption/
- https://web.cs.ucdavis.edu/~rogaway/papers/ad.pdf (needed to check up on security guarantee for "associated data", i.e. needed to confirm that "associated data" is a necessary component for decryption)
- https://tools.ietf.org/html/rfc5116
- [YouTube: Time-Lock Puzzles from Randomized Encodings](https://www.youtube.com/watch?v=bRcegZugqfY)

The following references were inspiration to this custom solution.

- https://www.gwern.net/Self-decrypting-files

References still to read:

- http://people.seas.harvard.edu/~salil/research/timelock.pdf
- https://www.youtube.com/watch?v=bRcegZugqfY
- https://www.youtube.com/watch?v=fqI35RcNdn8

[AEAD-paper]: https://web.cs.ucdavis.edu/~rogaway/papers/ad.pdf "Authenticated-Encryption with Associated-Data - Phillip Rogaway"
