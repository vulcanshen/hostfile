# hostfile

> Cross-platform hosts file manager CLI.
> Manages entries in a dedicated block without touching hand-written content.
> More information: <https://github.com/vulcanshen/hostfile>.

- Initialize hostfile management (backs up original as "origin"):

`hostfile init`

- Add domains to an IP address:

`hostfile add {{ip}} {{domain1 domain2 ...}}`

- Show all managed entries:

`hostfile show`

- Search for an IP or domain:

`hostfile search {{ip|domain}}`

- Disable/enable an entry without deleting it:

`hostfile disable {{ip|domain}}`
`hostfile enable {{ip|domain}}`

- Save and load snapshots:

`hostfile save {{name}}`
`hostfile load {{name}}`

- Import entries from a file (replace or merge):

`hostfile apply {{path/to/file}}`
`hostfile merge {{path/to/file}}`

- Remove an entry:

`hostfile remove {{ip|domain}}`
