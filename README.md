# godepgraph

godepgraph is a program for generating a dependency graph of Go packages.

[![Go](https://github.com/kisielk/godepgraph/actions/workflows/go.yml/badge.svg)](https://github.com/kisielk/godepgraph/actions/workflows/go.yml)

## Install

    go install github.com/kisielk/godepgraph@latest

## Use

For basic usage, just give the package path of interest as the first
argument:

    godepgraph github.com/kisielk/godepgraph
    
If you intend to graph a go mod project, your package should be passed as a relative path:

    godepgraph ./pkg/api

The default output is a graph in [Graphviz][graphviz] dot format. If you have the
graphviz tools installed you can render it by piping the output to dot:

    godepgraph github.com/kisielk/godepgraph | dot -Tpng -o godepgraph.png

You can also generate [Mermaid](https://mermaid.js.org/) graphs:

    godepgraph -format mermaid github.com/kisielk/godepgraph > graph.mmd

By default godepgraph will display packages in the standard library in the
graph, though it will not delve in to their dependencies.

## Colors

godepgraph uses a simple color scheme to denote different types of packages:

  * *green*: a package that is part of the Go standard library, installed in `$GOROOT`.
  * *blue*: a regular Go package found in `$GOPATH`.
  * *yellow*: a vendored Go package found in `$GOPATH`.
  * *orange*: a package found in `$GOPATH` that uses cgo by importing the special package "C".

## Ignoring Imports

### The Go Standard Library

If you want to ignore standard library packages entirely, use the -s flag:

    godepgraph -s github.com/kisielk/godepgraph

### Vendored Libraries

If you want to ignore vendored packages entirely, use the -novendor flag:

    godepgraph -novendor github.com/something/else

### By Name

Import paths can be ignored in a comma-separated list passed to the -i flag:

    godepgraph -i github.com/foo/bar,github.com/baz/blah github.com/something/else

The packages and their imports will be excluded from the graph, unless the imports
are also imported by another package which is not excluded.

### By Prefix

Import paths can also be ignored by prefix. The -p flag takes a comma-separated
list of prefixes:

    godepgraph -p github.com,launchpad.net bitbucket.org/foo/bar

## Example

Here's some example output for godepgraph itself:

Using `godepgraph -format mermaid github.com/kisielk/godepgraph`:

```mermaid
flowchart TD

classDef goroot fill:#1D4,color:white
classDef cgofiles fill:#D52,color:white
classDef vendored fill:#D90,color:white
classDef buildErrs fill:#C10,color:white

1[flag]
click 1 href "https://godoc.org/flag"
class 1 goroot

2[fmt]
click 2 href "https://godoc.org/fmt"
class 2 goroot

3[github.com/kisielk/godepgraph]
click 3 href "https://godoc.org/github.com/kisielk/godepgraph"
3 --> 1
3 --> 2
3 --> 4
3 --> 5
3 --> 6
3 --> 7
3 --> 8

4[go/build]
click 4 href "https://godoc.org/go/build"
class 4 goroot

5[log]
click 5 href "https://godoc.org/log"
class 5 goroot

6[os]
click 6 href "https://godoc.org/os"
class 6 goroot

7[sort]
click 7 href "https://godoc.org/sort"
class 7 goroot

8[strings]
click 8 href "https://godoc.org/strings"
class 8 goroot
```

Using `godepgraph -format graphviz github.com/kisielk/godepgraph`:

![Example output](example.png)

[graphviz]: http://graphviz.org

Here's the output of `godepgraph -format mermaid github.com/kisielk/errcheck`:

```mermaid
flowchart TD

classDef goroot fill:#1D4,color:white
classDef cgofiles fill:#D52,color:white
classDef vendored fill:#D90,color:white
classDef buildErrs fill:#C10,color:white

1[bufio]
click 1 href "https://godoc.org/bufio"
class 1 goroot

2[bytes]
click 2 href "https://godoc.org/bytes"
class 2 goroot

3[context]
click 3 href "https://godoc.org/context"
class 3 goroot

4[crypto/md5]
click 4 href "https://godoc.org/crypto/md5"
class 4 goroot

5[encoding/binary]
click 5 href "https://godoc.org/encoding/binary"
class 5 goroot

6[encoding/json]
click 6 href "https://godoc.org/encoding/json"
class 6 goroot

7[errors]
click 7 href "https://godoc.org/errors"
class 7 goroot

8[flag]
click 8 href "https://godoc.org/flag"
class 8 goroot

9[fmt]
click 9 href "https://godoc.org/fmt"
class 9 goroot

10[github.com/kisielk/errcheck]
click 10 href "https://godoc.org/github.com/kisielk/errcheck"
10 --> 8
10 --> 9
10 --> 11
10 --> 12
10 --> 13
10 --> 14
10 --> 15
10 --> 16
10 --> 17
10 --> 18

11[github.com/kisielk/errcheck/errcheck]
click 11 href "https://godoc.org/github.com/kisielk/errcheck/errcheck"
11 --> 1
11 --> 2
11 --> 7
11 --> 9
11 --> 19
11 --> 20
11 --> 21
11 --> 22
11 --> 12
11 --> 13
11 --> 23
11 --> 15
11 --> 24
11 --> 17

19[go/ast]
click 19 href "https://godoc.org/go/ast"
class 19 goroot

25[go/build]
click 25 href "https://godoc.org/go/build"
class 25 goroot

26[go/constant]
click 26 href "https://godoc.org/go/constant"
class 26 goroot

27[go/parser]
click 27 href "https://godoc.org/go/parser"
class 27 goroot

28[go/scanner]
click 28 href "https://godoc.org/go/scanner"
class 28 goroot

20[go/token]
click 20 href "https://godoc.org/go/token"
class 20 goroot

21[go/types]
click 21 href "https://godoc.org/go/types"
class 21 goroot

29[golang.org/x/mod/semver]
click 29 href "https://godoc.org/golang.org/x/mod/semver"
29 --> 24

30[golang.org/x/sync/errgroup]
click 30 href "https://godoc.org/golang.org/x/sync/errgroup"
30 --> 3
30 --> 9
30 --> 18

22[golang.org/x/tools/go/analysis]
click 22 href "https://godoc.org/golang.org/x/tools/go/analysis"
22 --> 8
22 --> 9
22 --> 19
22 --> 20
22 --> 21
22 --> 23
22 --> 17
22 --> 31

32[golang.org/x/tools/go/gcexportdata]
click 32 href "https://godoc.org/golang.org/x/tools/go/gcexportdata"
32 --> 1
32 --> 2
32 --> 6
32 --> 9
32 --> 20
32 --> 21
32 --> 33
32 --> 34
32 --> 13
32 --> 35

12[golang.org/x/tools/go/packages]
click 12 href "https://godoc.org/golang.org/x/tools/go/packages"
12 --> 2
12 --> 3
12 --> 6
12 --> 7
12 --> 9
12 --> 19
12 --> 27
12 --> 28
12 --> 20
12 --> 21
12 --> 30
12 --> 32
12 --> 36
12 --> 37
12 --> 38
12 --> 39
12 --> 13
12 --> 35
12 --> 40
12 --> 14
12 --> 23
12 --> 16
12 --> 41
12 --> 24
12 --> 42
12 --> 17
12 --> 18
12 --> 43
12 --> 44
12 --> 31

45[golang.org/x/tools/go/types/objectpath]
click 45 href "https://godoc.org/golang.org/x/tools/go/types/objectpath"
45 --> 9
45 --> 21
45 --> 46
45 --> 38
45 --> 42
45 --> 17

47[golang.org/x/tools/go/types/typeutil]
click 47 href "https://godoc.org/golang.org/x/tools/go/types/typeutil"
47 --> 2
47 --> 9
47 --> 19
47 --> 21
47 --> 48
47 --> 49
47 --> 18
47 --> 50

46[golang.org/x/tools/internal/aliases]
click 46 href "https://godoc.org/golang.org/x/tools/internal/aliases"
46 --> 19
46 --> 27
46 --> 20
46 --> 21

51[golang.org/x/tools/internal/event]
click 51 href "https://godoc.org/golang.org/x/tools/internal/event"
51 --> 3
51 --> 52
51 --> 53
51 --> 54

52[golang.org/x/tools/internal/event/core]
click 52 href "https://godoc.org/golang.org/x/tools/internal/event/core"
52 --> 3
52 --> 9
52 --> 53
52 --> 54
52 --> 43
52 --> 44
52 --> 50

53[golang.org/x/tools/internal/event/keys]
click 53 href "https://godoc.org/golang.org/x/tools/internal/event/keys"
53 --> 9
53 --> 54
53 --> 34
53 --> 55
53 --> 24
53 --> 42
53 --> 17

54[golang.org/x/tools/internal/event/label]
click 54 href "https://godoc.org/golang.org/x/tools/internal/event/label"
54 --> 9
54 --> 34
54 --> 23
54 --> 50

33[golang.org/x/tools/internal/gcimporter]
click 33 href "https://godoc.org/golang.org/x/tools/internal/gcimporter"
33 --> 1
33 --> 2
33 --> 5
33 --> 7
33 --> 9
33 --> 25
33 --> 26
33 --> 20
33 --> 21
33 --> 45
33 --> 46
33 --> 56
33 --> 38
33 --> 34
33 --> 57
33 --> 13
33 --> 35
33 --> 14
33 --> 23
33 --> 24
33 --> 42
33 --> 17
33 --> 18
33 --> 50

36[golang.org/x/tools/internal/gocommand]
click 36 href "https://godoc.org/golang.org/x/tools/internal/gocommand"
36 --> 2
36 --> 3
36 --> 6
36 --> 7
36 --> 9
36 --> 29
36 --> 51
36 --> 53
36 --> 54
36 --> 34
36 --> 39
36 --> 13
36 --> 35
36 --> 14
36 --> 15
36 --> 16
36 --> 42
36 --> 17
36 --> 18
36 --> 58
36 --> 44

37[golang.org/x/tools/internal/packagesinternal]
click 37 href "https://godoc.org/golang.org/x/tools/internal/packagesinternal"

56[golang.org/x/tools/internal/pkgbits]
click 56 href "https://godoc.org/golang.org/x/tools/internal/pkgbits"
56 --> 2
56 --> 4
56 --> 5
56 --> 7
56 --> 9
56 --> 26
56 --> 20
56 --> 34
56 --> 57
56 --> 13
56 --> 16
56 --> 42
56 --> 17

59[golang.org/x/tools/internal/stdlib]
click 59 href "https://godoc.org/golang.org/x/tools/internal/stdlib"
59 --> 9
59 --> 17

48[golang.org/x/tools/internal/typeparams]
click 48 href "https://godoc.org/golang.org/x/tools/internal/typeparams"
48 --> 2
48 --> 7
48 --> 9
48 --> 19
48 --> 20
48 --> 21
48 --> 46
48 --> 13
48 --> 17

38[golang.org/x/tools/internal/typesinternal]
click 38 href "https://godoc.org/golang.org/x/tools/internal/typesinternal"
38 --> 9
38 --> 19
38 --> 20
38 --> 21
38 --> 47
38 --> 46
38 --> 59
38 --> 60
38 --> 23
38 --> 42
38 --> 17
38 --> 50

60[golang.org/x/tools/internal/versions]
click 60 href "https://godoc.org/golang.org/x/tools/internal/versions"
60 --> 19
60 --> 21
60 --> 17

49[hash/maphash]
click 49 href "https://godoc.org/hash/maphash"
class 49 goroot

34[io]
click 34 href "https://godoc.org/io"
class 34 goroot

39[log]
click 39 href "https://godoc.org/log"
class 39 goroot

55[math]
click 55 href "https://godoc.org/math"
class 55 goroot

57[math/big]
click 57 href "https://godoc.org/math/big"
class 57 goroot

13[os]
click 13 href "https://godoc.org/os"
class 13 goroot

35[os/exec]
click 35 href "https://godoc.org/os/exec"
class 35 goroot

40[path]
click 40 href "https://godoc.org/path"
class 40 goroot

14[path/filepath]
click 14 href "https://godoc.org/path/filepath"
class 14 goroot

23[reflect]
click 23 href "https://godoc.org/reflect"
class 23 goroot

15[regexp]
click 15 href "https://godoc.org/regexp"
class 15 goroot

16[runtime]
click 16 href "https://godoc.org/runtime"
class 16 goroot

41[slices]
click 41 href "https://godoc.org/slices"
class 41 goroot

24[sort]
click 24 href "https://godoc.org/sort"
class 24 goroot

42[strconv]
click 42 href "https://godoc.org/strconv"
class 42 goroot

17[strings]
click 17 href "https://godoc.org/strings"
class 17 goroot

18[sync]
click 18 href "https://godoc.org/sync"
class 18 goroot

43[sync/atomic]
click 43 href "https://godoc.org/sync/atomic"
class 43 goroot

58[syscall]
click 58 href "https://godoc.org/syscall"
class 58 goroot

44[time]
click 44 href "https://godoc.org/time"
class 44 goroot

31[unicode]
click 31 href "https://godoc.org/unicode"
class 31 goroot

50[unsafe]
click 50 href "https://godoc.org/unsafe"
class 50 goroot
```
