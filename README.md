# WORK IN PROGRESS

This is a demo/toy program to test the implementation and performance of
the [net/styx](https://aqwari.net/net/styx) package, which itself is a
work in progress. If you are having trouble getting things to work, ensure
you have checked out the latest version of the `aqwari.net/net/styx`
package.

# BUILD

	go build

# USE

Start jsonfs on port 5640:

	./jsonfs -a localhost:5640 example.json

Using plan9port's `9pfuse` utility, mount the fs:

	9pfuse localhost:5640 /mnt/jsonfs

If you have a recent (2.6+) linux kernel, you can
mount using the kernel's `v9fs` implementation.
Unfortunately you'll need root access to do so
without modifying `/etc/fstab`:

	sudo modprobe 9p
	sudo mount -t 9p -o \
		tcp,name=`whoami`,uname=`whoami`,port=5640 \
		127.0.0.1 /mnt/jsonfs
	
Try looking around

	$ ls /mnt/jsonfs
	apiVersion data
	$ cat /mnt/jsonfs/apiVersion
	2.0

You should see output from jsonfs, such as

	accepted connection from 127.0.0.1:36602
	→ 65535 Tversion msize=8192 version="9P2000"
	← 65535 Rversion msize=8192 version="9P2000"
	→ 000 Tattach fid=1 afid=NOFID uname="droyo" aname=""
	← 000 Rattach qid="type=128 ver=0 path=1"
	→ 000 Twalk fid=1 newfid=2 "apiVersion"
	← 000 Rwalk wqid="type=0 ver=0 path=2"
	→ 000 Topen fid=2 mode=0
	← 000 Ropen qid="type=0 ver=0 path=2" iounit=0
	→ 000 Tread fid=2 offset=0 count=8168
	← 000 Rread count=3
	→ 000 Tread fid=2 offset=3 count=8168
	← 000 Rread count=0

When using example.json, the tree hierarchy should look
something like this:

	$ tree /mnt/jsonfs
	/mnt/jsonfs
	├── apiVersion
	└── data
	    ├── items
	    │   └── 0
	    │       ├── accessControl
	    │       │   ├── comment
	    │       │   ├── commentVote
	    │       │   ├── embed
	    │       │   ├── list
	    │       │   ├── rate
	    │       │   ├── syndicate
	    │       │   └── videoRespond
	    │       ├── aspectRatio
	    │       ├── category
	    │       ├── commentCount
	    │       ├── content
	    │       │   ├── 1
	    │       │   ├── 5
	    │       │   └── 6
	    │       ├── description
	    │       ├── duration
	    │       ├── favoriteCount
	    │       ├── id
	    │       ├── player
	    │       │   └── default
	    │       ├── rating
	    │       ├── ratingCount
	    │       ├── status
	    │       │   ├── reason
	    │       │   └── value
	    │       ├── tags
	    │       │   ├── 0
	    │       │   ├── 1
	    │       │   └── 2
	    │       ├── thumbnail
	    │       │   ├── default
	    │       │   └── hqDefault
	    │       ├── title
	    │       ├── updated
	    │       ├── uploaded
	    │       ├── uploader
	    │       └── viewCount
	    ├── itemsPerPage
	    ├── startIndex
	    ├── totalItems
	    └── updated
