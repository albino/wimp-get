# wimp-get

Tool to download FLACs from WiMP/Tidal.

## Binaries

Windows binaries are available in the [releases](https://git.fuwafuwa.moe/albino/wimp-get/releases) tab.

## Compiling

You need Go >= 1.8, and suitable ffmpeg, sox and mktorrent binaries. Once you've got that, just run `go build -ldflags "-s -w"`, and a `wimp-get` binary should be generated.

## Configuration

If your copies of ffmpeg, sox or mktorrent are outwith your path, put their paths in magic.json. Extract the `sessionId` of your WiMP session using your browser's developer tools, and insert that into magic.json too.

## Help!

If it's not working, come and ask for help on IRC. Connect to `irc.rizon.net`, type `/join #wimp-get` and get my attention by saying my name (albino).

## License info

wimp-get copyright (C) 2017 Al Beano  
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License version 2 (see the `LICENSE` file), or, at your option, any later version published by the Free Software Foundation.
