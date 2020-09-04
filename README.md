# clay-code-challenge
Backend code challenge for Clay

## The API 

The API features a single query endpoint mounted at the root, which takes two parameters -- `mode` and `username`, and returns the latest change diff in HTML form. 

_For example this query_

`http://localhost:8080/?mode=biography;username=Achariam`

_Yields this:_

<body><span>Building http://clay.earth. Previously product lead at @GoldmanSachs, Forbes 30 Under 30, @Harvard CS.</span><ins style="background:#e6ffe6;"> Native New Yorker</ins></body>

## How it Works

The program here has two components:

1. a monitor loop that runs in the background, polling the underlying twitter endpoint and writes new values into the database
2. An HTTP server that reads from the database and serves the diff HTML you see in the browser

Golang was a perfect choice here, as it provides the ability to do asynchronous polling and the ability to serve HTTP requests out of the box. NodeJS was a possibility as well, but Golang solves many of the problems that one encounters using NodeJS at production scale (in my experience).

More importantly, golang had a port of the `diff-patch-match` module that was designed for google docs. In fact the creator of the original module cited the module used here when dismissing a request to design one himself.

## Features, Improvements, and Scaling

### Making the interface more intuitive
The first improvement I would make here is to make `time` a query parameter. As of now, the diff returned is the latest change, but the ability to rewind back in time would be nice. Perhaps two time parameters `begin` and `end` which allow the user to bookend the query. Incorporating this is as simply as adding to the `WHERE` clause of the query in `getLastTwoEvents`

As well, I might add a homepage with a UI, move the query endpoint to a `/diff` extension, and improve the styling of the HTML with a CSS framework such as bulma or bootstrap

### Scaling to Higher Load

How to scale this product depends entirely on the needs that evolve after deploying, but hypothetically let's say the following loose criteria are in place:

1. The diff API needs to be available nearly all the time.
2. Moderate delays in accuracy of the data are tolerable.

What other behavior can we expect at scale? Twitter itself deals with highly 'bursty' reads, e.g. a celebrity tweets something controversial and suddenly demand spikes to read their profile. I imagine we would deal with this ourselves. 

Here are some strategies I'd take when scaling this system:

#### Break out reading from writing

The API server needs to be available nearly 24/7, so we need to have more than one EC2 instance sitting behind a load balancer. AWS reserves the right to terminate any EC2 instance at any time, for any reason. On top of that, Golang's garbage collector can take out a server for a couple seconds, or a hard disk can fail, or an entire availability zone can go down (if unlikely, it happens).

However, this requires us to break out the monitor loop into it's own worker pool -- otherwise we have a large number of unnecessary workers polling the same data.

In other words, we break out the system that reads the data from the system that writes the data -- a very common pattern. This is often done at the database level as well, i.e. one database that receives writes only and one database that receives reads. 

As mentioned earlier, a nice side-effect of this stategy is improved ability to serve bursty reads i.e. a spike in demand for one user's data. The read pipeline can handle bursty data more effecitvely; the read-only database can cache frequently accessed data more effectively than it would if it had to handle both reads and writes.

#### Deleting old data

Postgres databases I've worked with tend to have performance issues when dealing with tables of about `1 billion` or more rows. Index cruft begins to develop, and queries take longer times.

At that point, we'd start to think about a partition schema which allows for bulk deletes and prevents cruft. Older data is probably less frequently accessed, and so partitioning time into months and periodically deleting the oldest month, for example, might be helpful. 

#### Security

Given more time, I'd tighten up a couple security holes. In particular, SQL injection isn't prevented against here.

Access control and authentication is also completely absent -- any user can see anyones diffs, and it's not obvious that we want that out of the box. 
