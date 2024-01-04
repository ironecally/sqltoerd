# sqltoerd
Read DDL from .sql and generate ERD (Entity Relationship Diagram) based on [mermaid.ink](https://mermaid.ink) format!

This tool will read DDL (Data Definition Language) from.sql file that put under `source` folder, and parsing the result to `result` folder.

# Prerequisite
- `npm` installed (this is required to install mermaid-cli)
- installation of [mermaid-cli](https://github.com/mermaid-js/mermaid-cli) (mmdc) binary
    - verify by running ```mmdc -V```
- installation of [go1.17](https://go.dev/dl/) or later

# Usage

## With Golang Installed
1. To ensure all the required folder is ready, Run ```make init``` 
2. Ensure that `mmdc` binary installed by running ```mmdc --version```
3. Try to build the binary with ```make build```
4. To use the tools, simply use the command ```make run filename=source/dekara_dump.sql output=result/result-custom.pdf```
    - `filename` argument should refer to sql file that you're trying to convert
    - `output` argument dictates the output filename and extension, default value is `result/result.pdf`
5. You can find the result ERD in `result/` folder of the directory

## Binary Only
**Work in Progress**

