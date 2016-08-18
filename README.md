# WORK IN PROGRESS

The main purpose of this program is to test the implementation
and performance of the [net/styx](https://aqwari.net/net/styx)
package, which itself is a work in progress.

# BUILD

	go build

# USE

Start jsonfs on port 5640:

	./jsonfs -a localhost:5640 example.json
	
Using [Plan 9 from userspace][p9p]:

	mkdir mnt
	9 mount localhost:5640 mnt

You should see output from jsonfs, such as

	accepted connection from 127.0.0.1:32876
	65535 Tversion msize=8192 version="9P2000"
	0 Tattach fid=0 afid=4294967295 uname="droyo" aname=""

When using example.json, the tree hierarchy should look
something like this:

	├── glossary
	│   ├── title
	│   └── GlossDiv
	│       ├── title
	│       └── GlossList
	│           └── GlossEntry
	│               ├── ID
	│               ├── GlossTerm
	│               ├── Acronym
	│               └── GlossDef
	│                   ├── para
	│                   └── GlossSeeAlso
	│               ├── GlossSee

[p9p]: https://swtch.com/plan9port/
