This package exists to retrieve Go's signal trampoline function.

This function is called "runtime·sigtramp" on all platforms.
We can't access it from our package because it's not exported.
So instead we extract it using Go assembly.

It should be theoretically possible to extract it on every combination
of platform and OS using Go assembly suitable for each platform/OS,
however since go-libedit only needs it on darwin, we implement it only
there.
