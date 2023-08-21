# Amber: Scalable LLM-based Log Analysis Tool

Amber is a scalable, fault-tolerant log analysis system for use LLM analyzers.

Although Amber in its current form is more geared towards use with LLMs, it can be used with any kind of analyzer, though the benefits fade away. 

Amber comprises 3 systems:
1. [CLI](http://www.github.com/enckrish/amber) (this)
2. [Router](http://www.github.com/enckrish/amberine-router)
3. Analyzers (custom, should adhere docs/analyzer_protocol.md)

Amber CLI is written in Go, using [tview](https://github.com/rivo/tview) for the terminal UI.