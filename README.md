taskstats [![Build Status](https://travis-ci.org/mdlayher/taskstats.svg?branch=master)](https://travis-ci.org/mdlayher/taskstats) [![GoDoc](https://godoc.org/github.com/mdlayher/taskstats?status.svg)](https://godoc.org/github.com/mdlayher/taskstats) [![Go Report Card](https://goreportcard.com/badge/github.com/mdlayher/taskstats)](https://goreportcard.com/report/github.com/mdlayher/taskstats)
=========

Package `taskstats` provides access to Linux's taskstats interface, for sending
per-task and per-process statistics from the kernel to userspace.

For more information on taskstats, please see:
  - https://www.kernel.org/doc/Documentation/accounting/taskstats.txt
  - https://www.kernel.org/doc/Documentation/accounting/taskstats-struct.txt

MIT Licensed.