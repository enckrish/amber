# Amber: Scalable LLM-based Log Analysis Tool
Amber is a scalable, fault-tolerant log analysis system for use with LLM analyzers.

## Architecture
![architecture](https://github.com/enckrish/amber/assets/41233911/03cd4b35-0e57-428e-96cc-daa0a2c44e28)

Although Amber in its current form is more geared towards use with LLMs, it can be used with any kind of analyzer, though the benefits fade away. 

Amber comprises 3 systems:
1. [CLI](http://www.github.com/enckrish/amber) (this)
2. [Router](http://www.github.com/enckrish/amberine-router)
3. Analyzers, implemented by end users, and must adhere to the [protocol](https://github.com/enckrish/amber/blob/master/docs/analyzer_protocol.md). Example python3 implementation is available at [enckrish/aquamarine](https://www.github.com/enckrish/aquamarine).

The docs for using Amber can be found [here](https://github.com/enckrish/amber/blob/master/docs).

Amber CLI is written in Go, using [tview](https://github.com/rivo/tview) for the terminal UI.

The presentation on Amber can be found at https://drive.google.com/file/d/1l_A511ot6_os74rmG7KmtN4plHdor_do/view?usp=sharing.
