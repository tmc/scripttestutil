/*
The scripttest command assists with runninng tests against commands.

It is a cli wrapper over the rsc.io/script and rsc.io/script/scripttest packages.

Usage:

scripttest [-v] <command> [args...]

Commands:

  test, run    run scripttest files (default pattern: testdata/*.txt)
               scripttest test                # uses -p or default pattern
               scripttest test 'custom/*.txt' # overrides pattern

  scaffold     create scripttest scaffold in [dir]
               scripttest scaffold .

  infer        infer command info in [dir]
               scripttest infer .
*/
package main
