# aoc-go-template
Go template project for advent of code is intended for use with gonew see https://go.dev/blog/gonew

Will automatically download input and puzzle instructions. 
Need session Cookie to work here are 3 ways to provide the value for the session cookie.

1. In an ADVENT_OF_CODE_SESSION environment variable.
2. In a file called .adventofcode.session (note the dot) in your home directory /home/alice on Linux (not tested for other operating systems).
3. In a file called adventofcode.session (no dot) in your user's config directory /home/alice/.config on Linux (not tested for other operating systems).

inspired by https://github.com/scarvalhojr/aoc-cli
