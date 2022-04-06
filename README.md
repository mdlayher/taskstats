# taskstats [![Test Status](https://github.com/mdlayher/taskstats/workflows/Test/badge.svg)](https://github.com/mdlayher/taskstats/actions) [![Go Reference](https://pkg.go.dev/badge/github.com/mdlayher/taskstats.svg)](https://pkg.go.dev/github.com/mdlayher/taskstats)  [![Go Report Card](https://goreportcard.com/badge/github.com/mdlayher/taskstats)](https://goreportcard.com/report/github.com/mdlayher/taskstats)

Package `taskstats` provides access to Linux's taskstats interface, for sending
per-task, per-process, and cgroup statistics from the kernel to userspace. MIT
Licensed.

For more information on taskstats, please see:
  - https://www.kernel.org/doc/Documentation/accounting/cgroupstats.txt
  - https://www.kernel.org/doc/Documentation/accounting/taskstats.txt
  - https://www.kernel.org/doc/Documentation/accounting/taskstats-struct.txt
  - https://andrestc.com/post/linux-delay-accounting/

## Notes

* When instrumenting Go programs, use either the `taskstats.Self()` or
  `taskstats.TGID()` method.  Using the `PID()` method on multithreaded
  programs, including Go programs, will produce inaccurate results.

* Access to taskstats requires that the application have at least `CAP_NET_RAW`
  capability (see
  [capabilities(7)](http://man7.org/linux/man-pages/man7/capabilities.7.html)).
  Otherwise, the application must be run as root.

* If running the application in a container (e.g. via Docker), it cannot be run
  in a network namespace -- usually this means that host networking must be
  used.
