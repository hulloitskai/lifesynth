# lifesynth

_A system for synthesizing and contextualizing personal data._

## Tools

### `journalmeta`

`journalmeta` is a tool for extracting metadata from my Markdown journal
entries.

Meta blocks are either code fences that are tagged with `meta`, like this:

    ```meta
    > 9:21 AM
    @ Strange Love Coffee
    ! 6/10
    # work coding
    ```

Or they can be standalone bits of inline code, in a form of shorthand, like
this:

`> 5:00 PM ! 0.6 @ [home]`

As of the current schema (`v1`), there are four parseable tokens:

- `>` – time
- `@` – location
- `!` – [valence](<https://en.wikipedia.org/wiki/Valence_(psychology)>)
- `#` – tags (levelled)

#### Usage:

`journalmeta` is invoked as follows:

```bash
journalmeta <path to markdown file>
```
