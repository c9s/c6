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




