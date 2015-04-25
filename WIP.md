Working In Progress
======================

Major Components
----------------

- Parser
- Server
- Client

Parser
-------------

Server
-------------

### Commands

- `check`
- `parse`
- `complete`
- `compile`
- `minify`
- `optimize`

#### `check`

The check command checks the syntax and reports error.

`check` also implictly parse the file and create a parse tree in the memory.

Protocol:

    check {path}

Response:

    {line}:{type}:{suggest}:{message}

Sample Response:

    23:error:append_comma:Missing comma at the end

#### `parse`

Parse a scss file and save the syntax tree, symbol table in the memory.

`parse /path/to/file.scss`

#### Client


- [ ] Supporting SASS comments
- [ ] Supporting pseudo elements:  ::first-line Pseudo-element
- [ ] Supporting more length unit:

    http://www.w3.org/TR/2011/REC-CSS2-20110607/syndata.html#lengths

    in: inches — 1in is equal to 2.54cm.
    cm: centimeters
    mm: millimeters
    pt: points — the points used by CSS are equal to 1/72nd of 1in.
    pc: picas — 1pc is equal to 12pt.
    px: pixel units — 1px is equal to 0.75pt.

- [ ] unicode range



http://sass-lang.com/documentation/file.SASS_REFERENCE.html


Reference
-------------

sass reference
http://sass-lang.com/documentation/file.SASS_REFERENCE.html

ambiguity
https://www.facebook.com/cindylinz/posts/10202186527405801?hc_location=ufi


