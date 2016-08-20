# WORK IN PROGRESS

The main purpose of this program is to test the implementation
and performance of the [net/styx](https://aqwari.net/net/styx)
package, which itself is a work in progress.

# BUILD

	go build

# USE

Start jsonfs on port 5640:

	./jsonfs -a localhost:5640 example.json

Install the `ixpc` client from [libixp][ixp]:

	sudo apt-get install libixp-dev

Try listing something

	$ export IXP_ADDRESS='tcp!localhost!5640'
	$ ixpc ls /
	apiVersion
	data/
	$ ixpc read apiVersion
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

	bootes=; ixpc ls data/items
	accessControl/
	aspectRatio
	category
	commentCount
	content/
	description
	duration
	favoriteCount
	id
	player/
	rating
	ratingCount
	status/
	tags
	thumbnail/
	title
	updated
	uploaded
	uploader
	viewCount

[ixp]: https://bitbucket.org/kmaglione/libixp
