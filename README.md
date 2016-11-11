# omnissh
SSH-based Backdoor

In development.  When complete (or at least useable), it will be a single-file
backdoor which uses the SSH protocol but implements its own shell.

After reading a paper on detecting malware based on system calls, it became
apparent that backdoors should probably avoid a fork/exec for every command.
This, then, attempts to be something like a cross-platform dropbear and
busybox.

The paper I read is
[http://repository.tudelft.nl/islandora/object/uuid:c71c85bc-d742-449b-88e7-33e172392ec2?collection=education].
Aside from some glaring errors, the author did make the valid point that, with
a few exceptions, most programs make their own system calls for most things
these days.  Whoops.
