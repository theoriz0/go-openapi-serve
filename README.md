# Go Openapi Serve

Serve openapi yaml file from a http server, by using [kin-openapi](https://github.com/getkin/kin-openapi).

## Usage

```bash
go-openapi-serve [port] ([path1] [path2]...)
```

If the path points to a folder, the yaml files in that folder will be served. However, it currently doesn't 
If the path points to a file, the corresponding file will be served.

For now, this tool not support:

- Absolute path to file or folder.
- `..` in relative path to file or folder.
- Scan yaml files in folder recursively.

## Examples

Serve yaml files under `./docs`

```bash
go-openapi-serve 8080
```

Serve yaml files under `./docs` and yaml file `./someapi.yaml` through port 8080:

```bash
go-openapi-serve 8080 docs someapi.yaml
```