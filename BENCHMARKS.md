# Generating kubernetest models

with minimal patch:
rm -rf temp && mkdir temp && /usr/bin/time swagger generate server -t temp -f fixtures/canary/kubernetes/swagger.json --quiet
2024/01/15 09:55:47 profile: memory profiling enabled (rate 4096), /home/fred/src/github.com/go-swagger/go-swagger/prof/prof-3694800708/mem.pprof
45.54user 2.29system 0:33.80elapsed 141%CPU (0avgtext+0avgdata 916192maxresident)k
3736inputs+31440outputs (0major+432247minor)pagefaults 0swaps

max RAM: 916,192 kB
alloc objects: 330,519,558 (ToGoName: 57%) (from templates Value.call: 211,711,764)

same run, with swag patch
32.51user 2.20system 0:23.98elapsed 144%CPU (0avgtext+0avgdata 1037848maxresident)k

max RAM: 1,037,848
alloc objects: 95,751,001 (ToGoName disappears from sample)
  * imports.Process: 35%
  * Validate: 19% (18M allocs)
  * analysis.New / initialize ~ 6M allocs


* real: 34.2 s, user CPU: 48s, total CPU: 51.21s
  * validateAndFlatten: 5.86
    * validate: 5.56
  * appGenerator.Generate: 19.50s 
    * ast.Parse: 2.20s
    * makeCodeGenApp: 18.46s
      * template render: 0.41s
      * imports.Process: 2.92s
        * fixImportDefault: 2.85s
  * runtime mallocgc: 2.09s (3.99%)
  * runtime.gcBMarkWorker: 15.39s (30%)
    * gcDrain: 0.45s
    * scanobject: 5.09s

* with new swag ~ same (~ -1%)
* with new spec: 27.7s, user CPU: 41s  (~ -5%)
* with both: 26.8s (40+3 CPU)
     
# Running the codegen CI suite

Running a patch swagger CLI that collects mem profiles.

## With spec using std lib

time go test -v -timeout 30m -parallel 3 hack/codegen_nonreg_test.go -args -fixture-file codegen-fixtures.yaml -skip-models -skip-full-flatten

real	4m38.270s
user	16m1.063s
sys	4m17.300s

time go test -v -timeout 30m -parallel 3 hack/codegen_nonreg_test.go -args -fixture-file canary-fixtures.yaml -skip-models -skip-expand -skip-full-flatten

real	5m6.957s
user	15m38.306s
sys	1m0.152s

pprof -tree -functions -sample_index=alloc_objects  merged.data
File: swagger
Build ID: e4c891238f5069a88d59a94ecf2f07f5d6371c02
Type: alloc_objects
Time: Jan 9, 2024 at 4:16pm (CET)
Showing nodes accounting for 2523606974, 84.43% of 2988991312 total
Dropped 1779 nodes (cum <= 14944956)
Showing top 80 nodes out of 353
----------------------------------------------------------+-------------
      flat  flat%   sum%        cum   cum%   calls calls% + context 	 	 
----------------------------------------------------------+-------------
                                         710103241 81.35% |   github.com/go-openapi/swag.ToGoName
                                         102256992 11.72% |   reflect.Value.call
 872861133 29.20% 29.20%  872861133 29.20%                | github.com/go-openapi/swag.(*splitter).gatherInitialismMatches
----------------------------------------------------------+-------------
                                          59964487 51.54% |   github.com/go-openapi/swag.(*splitter).breakCasualString
 116346464  3.89% 33.10%  116346464  3.89%                | strings.(*Builder).grow
----------------------------------------------------------+-------------
                                          59147141 50.60% |   go/parser.(*parser).parseOperand
                                          22890850 19.58% |   go/parser.(*parser).parseTypeName
                                          12856709 11.00% |   go/parser.(*parser).parseSelector
                                          10872725  9.30% |   go/parser.(*parser).parseParameterList
                                           4212810  3.60% |   go/parser.(*parser).parseFuncDecl
                                           3433226  2.94% |   go/parser.(*parser).parseGenDecl
 112232956  3.75% 36.85%  116881015  3.91%                | go/parser.(*parser).parseIdent
----------------------------------------------------------+-------------
                                          97571206 71.67% |   path/filepath.readDir
                                          22296125 16.38% |   golang.org/x/tools/internal/imports.(*ModuleResolver).canonicalize
  84449153  2.83% 39.68%  136133781  4.55%                | os.(*File).readdir
                                          51667618 37.95% |   os.newUnixDirent
----------------------------------------------------------+-------------
                                          32717882 54.88% |   encoding/json.(*decodeState).objectInterface (inline)
                                          26896461 45.12% |   encoding/json.(*decodeState).literalInterface (inline)
  59361690  1.99% 41.66%   59614343  1.99%                | encoding/json.unquote
----------------------------------------------------------+-------------
                                          10899466 19.13% |   go/parser.(*parser).parseBinaryExpr
                                           5756834 10.10% |   go/parser.(*parser).parseParameters
                                           5545559  9.73% |   go/parser.(*parser).parseIfStmt
                                           4435397  7.78% |   go/parser.(*parser).parsePointerType
                                           4047525  7.10% |   go/parser.(*parser).parseCallOrConversion
                                           3157126  5.54% |   go/parser.(*parser).parseBlockStmt
                                           3091939  5.43% |   go/parser.(*parser).parseTypeName
  56977598  1.91% 43.57%   56977598  1.91%                | go/scanner.(*Scanner).scanIdentifier
----------------------------------------------------------+-------------
                                         149164766 89.07% |   encoding/json.(*decodeState).valueInterface
                                         131756327 78.68% |   encoding/json.(*decodeState).object
  56509765  1.89% 45.46%  167461623  5.60%                | encoding/json.(*decodeState).objectInterface
                                         153532210 91.68% |   encoding/json.(*decodeState).valueInterface
                                          32717882 19.54% |   encoding/json.unquote (inline)
----------------------------------------------------------+-------------
                                          34717361 65.03% |   encoding/gob.(*Encoder).sendActualType
                                          11360733 21.28% |   encoding/gob.(*Decoder).compileDec
  53386542  1.79% 47.24%   53386542  1.79%                | reflect.(*structType).Field
----------------------------------------------------------+-------------
                                          52787400 99.93% |   go/parser.(*parser).parseFile (inline)
  52825378  1.77% 49.01%   52825378  1.77%                | go/ast.NewScope
----------------------------------------------------------+-------------
                                          51667618   100% |   os.(*File).readdir
  51667618  1.73% 50.74%   51667618  1.73%                | os.newUnixDirent
----------------------------------------------------------+-------------
                                          23192225 47.77% |   reflect.cvtBytesString
                                          10843393 22.34% |   encoding/json.(*decodeState).object
                                          10618332 21.87% |   encoding/gob.(*Decoder).recvType
  48544335  1.62% 52.36%   48546884  1.62%                | reflect.New
----------------------------------------------------------+-------------
                                         155046793 77.97% |   go/parser.(*parser).parseSimpleStmt
                                          42601359 21.42% |   go/parser.(*parser).parseReturnStmt
  40715769  1.36% 53.73%  198850390  6.65%                | go/parser.(*parser).parseExprList
                                         156779155 78.84% |   go/parser.(*parser).parseBinaryExpr
----------------------------------------------------------+-------------
                                          67229201   100% |   encoding/json.(*decodeState).valueInterface
  39596778  1.32% 55.05%   67229201  2.25%                | encoding/json.(*decodeState).literalInterface
                                          26896461 40.01% |   encoding/json.unquote (inline)
----------------------------------------------------------+-------------
                                         120650642   100% |   encoding/gob.(*Decoder).compileDec
  38135691  1.28% 56.33%  120650642  4.04%                | encoding/gob.(*Decoder).decOpFor
                                         115296834 95.56% |   encoding/gob.(*Decoder).getDecEnginePtr
                                           6538347  5.42% |   encoding/gob.overflow (inline)
----------------------------------------------------------+-------------
                                         398513845 98.65% |   encoding/json.Unmarshal
  35875547  1.20% 57.53%  403951853 13.51%                | encoding/json.(*decodeState).object
                                         298451409 73.88% |   github.com/go-openapi/spec.(*Schema).UnmarshalJSON
                                         150511066 37.26% |   encoding/json.Unmarshal
                                         131756327 32.62% |   encoding/json.(*decodeState).objectInterface
                                          43681468 10.81% |   encoding/json.(*decodeState).literalStore
                                          40626093 10.06% |   encoding/json.(*decodeState).arrayInterface
                                          36797541  9.11% |   reflect.cvtBytesString
                                          10843393  2.68% |   reflect.New
----------------------------------------------------------+-------------
                                          26125461 80.88% |   github.com/go-openapi/spec.expandSchemaRef (inline)
  32301168  1.08% 58.61%   32301168  1.08%                | strings.(*Builder).WriteString
----------------------------------------------------------+-------------
                                          23373422 72.91% |   encoding/gob.(*Decoder).decodeArrayHelper
                                           7115213 22.19% |   encoding/gob.(*Decoder).recvType
  32057769  1.07% 59.68%   32057769  1.07%                | encoding/gob.decString
----------------------------------------------------------+-------------
                                          43681468   100% |   encoding/json.(*decodeState).object
                                           6100240 13.97% |   encoding/json.Unmarshal
  28943379  0.97% 60.65%   43681524  1.46%                | encoding/json.(*decodeState).literalStore
                                          18971868 43.43% |   github.com/go-openapi/spec.(*StringOrArray).UnmarshalJSON
----------------------------------------------------------+-------------
                                          62585551   100% |   go/parser.(*parser).parseBinaryExpr
  27498183  0.92% 61.57%   62585935  2.09%                | go/parser.(*parser).parseCallOrConversion
                                          30018813 47.96% |   go/parser.(*parser).parseBinaryExpr
                                           4047525  6.47% |   go/scanner.(*Scanner).scanIdentifier
----------------------------------------------------------+-------------
                                          25752313 97.44% |   encoding/gob.overflow (inline)
  26429188  0.88% 62.45%   26429188  0.88%                | errors.New
----------------------------------------------------------+-------------
                                          44969457 87.31% |   encoding/gob.(*Decoder).compileDec (inline)
                                           6538347 12.69% |   encoding/gob.(*Decoder).decOpFor (inline)
  25755491  0.86% 63.31%   51507804  1.72%                | encoding/gob.overflow
                                          25752313 50.00% |   errors.New (inline)
----------------------------------------------------------+-------------
                                          66988427 83.10% |   github.com/go-openapi/spec.Schema.MarshalJSON
                                          10452208 12.97% |   github.com/go-openapi/spec.VendorExtensible.MarshalJSON
                                           9679357 12.01% |   encoding/json.mapEncoder.encode
  25295308  0.85% 64.16%   80612829  2.70%                | encoding/json.Marshal
                                          61141435 75.85% |   github.com/go-openapi/spec.Schema.MarshalJSON
                                          26789197 33.23% |   encoding/json.mapEncoder.encode
                                          17215695 21.36% |   github.com/go-openapi/spec.VendorExtensible.MarshalJSON
----------------------------------------------------------+-------------
                                         295480528 73.85% |   github.com/go-openapi/spec.(*Schema).UnmarshalJSON
                                         150511066 37.62% |   encoding/json.(*decodeState).object
                                          97998212 24.49% |   github.com/go-openapi/spec.expandSchemaRef
                                          10405666  2.60% |   github.com/go-openapi/spec.(*StringOrArray).UnmarshalJSON
  24982886  0.84% 65.00%  400107897 13.39%                | encoding/json.Unmarshal
                                         398513845 99.60% |   encoding/json.(*decodeState).object
                                          23746904  5.94% |   encoding/json.(*scanner).pushParseState
                                           6100240  1.52% |   encoding/json.(*decodeState).literalStore
----------------------------------------------------------+-------------
                                          23746904 99.93% |   encoding/json.Unmarshal
  23764181   0.8% 65.79%   23764181   0.8%                | encoding/json.(*scanner).pushParseState
----------------------------------------------------------+-------------
                                          12819975 54.87% |   path/filepath.readDir
                                           3724727 15.94% |   golang.org/x/tools/internal/imports.(*ModuleResolver).canonicalize
  23365592  0.78% 66.57%   23365592  0.78%                | os.newFile
----------------------------------------------------------+-------------
                                          11292632 38.94% |   github.com/go-openapi/analysis.(*Spec).analyzeSchema
  22995931  0.77% 67.34%   28998573  0.97%                | path.Join
----------------------------------------------------------+-------------
                                          20828396 99.91% |   go/parser.(*parser).parseFile (inline)
  20846972   0.7% 68.04%   20846972   0.7%                | go/ast.NewObj
----------------------------------------------------------+-------------
                                         290266363 98.72% |   go/parser.(*parser).parseBody
                                         132810075 45.17% |   go/parser.(*parser).parseBlockStmt
  20662977  0.69% 68.73%  294036589  9.84%                | go/parser.(*parser).parseStmtList
                                         216305120 73.56% |   go/parser.(*parser).parseIfStmt
                                          67483641 22.95% |   go/parser.(*parser).parseSimpleStmt
                                          53704185 18.26% |   go/parser.(*parser).parseReturnStmt
                                          25134614  8.55% |   go/parser.(*parser).parseBlockStmt
                                          13813027  4.70% |   go/parser.(*parser).parseGenDecl
----------------------------------------------------------+-------------
                                          75619736   100% |   go/parser.(*parser).parseParameters
  20634367  0.69% 69.42%   75619736  2.53%                | go/parser.(*parser).parseParameterList
                                          16391583 21.68% |   go/parser.(*parser).parseParameterList.func1 (inline)
                                          14514896 19.19% |   go/parser.(*parser).parsePointerType
                                          10872725 14.38% |   go/parser.(*parser).parseIdent
                                          10755134 14.22% |   go/parser.(*parser).parseTypeName
----------------------------------------------------------+-------------
                                           8619995 42.15% |   os.statNolog
                                           6378501 31.19% |   path/filepath.readDir
  20452902  0.68% 70.11%   20452902  0.68%                | syscall.ByteSliceFromString
----------------------------------------------------------+-------------
                                         115296834 87.17% |   encoding/gob.(*Decoder).decOpFor
                                          34033516 25.73% |   encoding/gob.(*Decoder).recvType
  18379827  0.61% 70.72%  132262562  4.42%                | encoding/gob.(*Decoder).getDecEnginePtr
                                         129038648 97.56% |   encoding/gob.(*Decoder).compileDec
----------------------------------------------------------+-------------
                                         108876172 61.75% |   go/parser.(*parser).parseIfStmt
                                          67483641 38.27% |   go/parser.(*parser).parseStmtList
  17964845   0.6% 71.32%  176326723  5.90%                | go/parser.(*parser).parseSimpleStmt
                                         155046793 87.93% |   go/parser.(*parser).parseExprList
----------------------------------------------------------+-------------
                                          16391583   100% |   go/parser.(*parser).parseParameterList (inline)
  16391583  0.55% 71.87%   16391583  0.55%                | go/parser.(*parser).parseParameterList.func1
----------------------------------------------------------+-------------
  15330924  0.51% 72.38%   23950919   0.8%                | os.statNolog
                                           8619995 35.99% |   syscall.ByteSliceFromString
----------------------------------------------------------+-------------
                                         298451409 99.77% |   encoding/json.(*decodeState).object
  14332776  0.48% 72.86%  299147541 10.01%                | github.com/go-openapi/spec.(*Schema).UnmarshalJSON
                                         295480528 98.77% |   encoding/json.Unmarshal
                                           3026731  1.01% |   github.com/go-openapi/jsonpointer.(*Pointer).parse
----------------------------------------------------------+-------------
  14303079  0.48% 73.34%   14996431   0.5%                | reflect.(*structType).FieldByNameFunc
----------------------------------------------------------+-------------
                                          36797541   100% |   encoding/json.(*decodeState).object
  13605316  0.46% 73.80%   36797541  1.23%                | reflect.cvtBytesString
                                          23192225 63.03% |   reflect.New
----------------------------------------------------------+-------------
                                          15917192 70.36% |   io.WriteString
  13518367  0.45% 74.25%   22623156  0.76%                | github.com/kr/text.(*indentWriter).Write
----------------------------------------------------------+-------------
                                          72577306 89.62% |   github.com/go-openapi/swag.ToGoName
                                           6840831  8.45% |   reflect.Value.call
  13512738  0.45% 74.70%   80985201  2.71%                | github.com/go-openapi/swag.(*splitter).breakCasualString
                                          59964487 74.04% |   strings.(*Builder).grow
----------------------------------------------------------+-------------
                                          25772591   100% |   go/parser.(*parser).parseBinaryExpr
  12915882  0.43% 75.13%   25772591  0.86%                | go/parser.(*parser).parseSelector
                                          12856709 49.89% |   go/parser.(*parser).parseIdent
----------------------------------------------------------+-------------
                                          73837918 76.92% |   encoding/json.(*decodeState).valueInterface
                                          40626093 42.32% |   encoding/json.(*decodeState).object
  12735317  0.43% 75.56%   95998978  3.21%                | encoding/json.(*decodeState).arrayInterface
                                          88047597 91.72% |   encoding/json.(*decodeState).valueInterface
----------------------------------------------------------+-------------
  12451576  0.42% 75.98%   51886640  1.74%                | encoding/gob.(*Encoder).sendActualType
                                          34717361 66.91% |   reflect.(*structType).Field
----------------------------------------------------------+-------------
  12018884   0.4% 76.38%  106015441  3.55%                | encoding/gob.(*Decoder).recvType
                                          34033516 32.10% |   encoding/gob.(*Decoder).getDecEnginePtr
                                          28792125 27.16% |   encoding/gob.(*Decoder).decodeArrayHelper
                                          12657474 11.94% |   reflect.MakeSlice
                                          10618332 10.02% |   reflect.New
                                           7115213  6.71% |   encoding/gob.decString
----------------------------------------------------------+-------------
                                          61141435 86.20% |   encoding/json.Marshal
                                          24235602 34.17% |   github.com/go-openapi/spec.expandSchemaRef
                                           7049257  9.94% |   encoding/json.mapEncoder.encode
  11669683  0.39% 76.77%   70931644  2.37%                | github.com/go-openapi/spec.Schema.MarshalJSON
                                          66988427 94.44% |   encoding/json.Marshal
----------------------------------------------------------+-------------
                                         129038648   100% |   encoding/gob.(*Decoder).getDecEnginePtr
  10882604  0.36% 77.13%  129038648  4.32%                | encoding/gob.(*Decoder).compileDec
                                         120650642 93.50% |   encoding/gob.(*Decoder).decOpFor
                                          44969457 34.85% |   encoding/gob.overflow (inline)
                                          11360733  8.80% |   reflect.(*structType).Field
----------------------------------------------------------+-------------
                                          26789197 91.50% |   encoding/json.Marshal
   9677401  0.32% 77.46%   29277083  0.98%                | encoding/json.mapEncoder.encode
                                           9679357 33.06% |   encoding/json.Marshal
                                           7049257 24.08% |   github.com/go-openapi/spec.Schema.MarshalJSON
----------------------------------------------------------+-------------
                                          21107515 98.91% |   go/parser.(*parser).parseFuncDecl
   9629441  0.32% 77.78%   21339715  0.71%                | go/parser.(*parser).parseResult
                                           7139164 33.45% |   go/parser.(*parser).parseParameters
----------------------------------------------------------+-------------
   9323801  0.31% 78.09%   36286056  1.21%                | github.com/go-openapi/analysis.(*Spec).analyzeSchema
                                          11292632 31.12% |   path.Join
                                           7208740 19.87% |   github.com/go-openapi/jsonpointer.(*Pointer).parse
----------------------------------------------------------+-------------
                                          53704185   100% |   go/parser.(*parser).parseStmtList
   9289812  0.31% 78.40%   53720019  1.80%                | go/parser.(*parser).parseReturnStmt
                                          42601359 79.30% |   go/parser.(*parser).parseExprList
----------------------------------------------------------+-------------
   8911559   0.3% 78.70%   90734475  3.04%                | golang.org/x/tools/internal/imports.(*dirInfoCache).Store
                                          74066792 81.63% |   golang.org/x/tools/internal/imports.(*ModuleResolver).canonicalize
----------------------------------------------------------+-------------
   8798392  0.29% 78.99%  111755938  3.74%                | github.com/go-openapi/validate.(*objectValidator).Validate
                                          61603333 55.12% |   github.com/go-openapi/validate.NewSchemaValidator
----------------------------------------------------------+-------------
                                          18971868   100% |   encoding/json.(*decodeState).literalStore
   8566202  0.29% 79.28%   18971868  0.63%                | github.com/go-openapi/spec.(*StringOrArray).UnmarshalJSON
                                          10405666 54.85% |   encoding/json.Unmarshal
----------------------------------------------------------+-------------
                                         411905908 99.22% |   go/parser.(*parser).parseFile
   8453319  0.28% 79.56%  415144370 13.89%                | go/parser.(*parser).parseFuncDecl
                                         294002592 70.82% |   go/parser.(*parser).parseBody
                                          81941226 19.74% |   go/parser.(*parser).parseParameters
                                          21107515  5.08% |   go/parser.(*parser).parseResult
                                           4212810  1.01% |   go/parser.(*parser).parseIdent
----------------------------------------------------------+-------------
                                         132754659 94.42% |   go/parser.(*parser).parseIfStmt
                                          25134614 17.88% |   go/parser.(*parser).parseStmtList
   8395788  0.28% 79.84%  140604634  4.70%                | go/parser.(*parser).parseBlockStmt
                                         132810075 94.46% |   go/parser.(*parser).parseStmtList
                                           3157126  2.25% |   go/scanner.(*Scanner).scanIdentifier
----------------------------------------------------------+-------------
                                           7208740 43.20% |   github.com/go-openapi/analysis.(*Spec).analyzeSchema
                                           5803825 34.78% |   github.com/go-openapi/spec.expandSchemaRef
                                           3026731 18.14% |   github.com/go-openapi/spec.(*Schema).UnmarshalJSON
   8385189  0.28% 80.12%   16687621  0.56%                | github.com/go-openapi/jsonpointer.(*Pointer).parse
----------------------------------------------------------+-------------
                                          81941226 91.55% |   go/parser.(*parser).parseFuncDecl
                                           7139164  7.98% |   go/parser.(*parser).parseResult
   8124899  0.27% 80.40%   89501483  2.99%                | go/parser.(*parser).parseParameters
                                          75619736 84.49% |   go/parser.(*parser).parseParameterList
                                           5756834  6.43% |   go/scanner.(*Scanner).scanIdentifier
----------------------------------------------------------+-------------
                                          12657474 79.54% |   encoding/gob.(*Decoder).recvType
   7965322  0.27% 80.66%   15914286  0.53%                | reflect.MakeSlice
----------------------------------------------------------+-------------
                                         216305120   100% |   go/parser.(*parser).parseStmtList
   7783922  0.26% 80.92%  216327103  7.24%                | go/parser.(*parser).parseIfStmt
                                         132754659 61.37% |   go/parser.(*parser).parseBlockStmt
                                         108876172 50.33% |   go/parser.(*parser).parseSimpleStmt
                                           5545559  2.56% |   go/scanner.(*Scanner).scanIdentifier
----------------------------------------------------------+-------------
                                         156779155 99.81% |   go/parser.(*parser).parseExprList
                                          30018813 19.11% |   go/parser.(*parser).parseCallOrConversion
   7703308  0.26% 81.18%  157071755  5.26%                | go/parser.(*parser).parseBinaryExpr
                                          67788960 43.16% |   go/parser.(*parser).parseOperand
                                          62585551 39.85% |   go/parser.(*parser).parseCallOrConversion
                                          25772591 16.41% |   go/parser.(*parser).parseSelector
                                          10899466  6.94% |   go/scanner.(*Scanner).scanIdentifier
                                           5522468  3.52% |   go/parser.(*parser).parsePointerType
----------------------------------------------------------+-------------
   7539951  0.25% 81.43%   23584715  0.79%                | io.WriteString
                                          15917192 67.49% |   github.com/kr/text.(*indentWriter).Write
----------------------------------------------------------+-------------
                                          67788960   100% |   go/parser.(*parser).parseBinaryExpr
   7355094  0.25% 81.68%   67789088  2.27%                | go/parser.(*parser).parseOperand
                                          59147141 87.25% |   go/parser.(*parser).parseIdent
----------------------------------------------------------+-------------
                                          28792125 83.01% |   encoding/gob.(*Decoder).recvType
   7186971  0.24% 81.92%   34684181  1.16%                | encoding/gob.(*Decoder).decodeArrayHelper
                                          23373422 67.39% |   encoding/gob.decString
----------------------------------------------------------+-------------
                                          17215695   100% |   encoding/json.Marshal
   6761436  0.23% 82.15%   17215695  0.58%                | github.com/go-openapi/spec.VendorExtensible.MarshalJSON
                                          10452208 60.71% |   encoding/json.Marshal
----------------------------------------------------------+-------------
                                          14514896 60.10% |   go/parser.(*parser).parseParameterList
                                           5522468 22.87% |   go/parser.(*parser).parseBinaryExpr
   6532981  0.22% 82.36%   24151541  0.81%                | go/parser.(*parser).parsePointerType
                                          13181284 54.58% |   go/parser.(*parser).parseTypeName
                                           4435397 18.36% |   go/scanner.(*Scanner).scanIdentifier
----------------------------------------------------------+-------------
   6462260  0.22% 82.58%  131804374  4.41%                | path/filepath.readDir
                                          97571206 74.03% |   os.(*File).readdir
                                          12819975  9.73% |   os.newFile
                                           6378501  4.84% |   syscall.ByteSliceFromString
----------------------------------------------------------+-------------
                                          36454243 72.52% |   go/parser.(*parser).parseFile
                                          13813027 27.48% |   go/parser.(*parser).parseStmtList
   6369712  0.21% 82.79%   50269047  1.68%                | go/parser.(*parser).parseGenDecl
                                          17728149 35.27% |   go/parser.(*parser).parseStructType
                                           3433226  6.83% |   go/parser.(*parser).parseIdent
----------------------------------------------------------+-------------
                                         128119488 69.01% |   github.com/go-openapi/validate.NewSchemaValidator
                                          91956735 49.53% |   github.com/go-openapi/spec.expandSchemaRef
                                          44372694 23.90% |   github.com/go-openapi/spec.expandItems
   5351878  0.18% 82.97%  185660778  6.21%                | github.com/go-openapi/spec.expandSchema
                                         184441918 99.34% |   github.com/go-openapi/spec.expandSchemaRef
                                          46162139 24.86% |   github.com/go-openapi/spec.expandItems
----------------------------------------------------------+-------------
                                          13181284 42.38% |   go/parser.(*parser).parsePointerType
                                          10755134 34.58% |   go/parser.(*parser).parseParameterList
   5119212  0.17% 83.14%   31102001  1.04%                | go/parser.(*parser).parseTypeName
                                          22890850 73.60% |   go/parser.(*parser).parseIdent
                                           3091939  9.94% |   go/scanner.(*Scanner).scanIdentifier
----------------------------------------------------------+-------------
                                         153532210 95.51% |   encoding/json.(*decodeState).objectInterface
                                          88047597 54.78% |   encoding/json.(*decodeState).arrayInterface
   5000611  0.17% 83.31%  160741477  5.38%                | encoding/json.(*decodeState).valueInterface
                                         149164766 92.80% |   encoding/json.(*decodeState).objectInterface
                                          73837918 45.94% |   encoding/json.(*decodeState).arrayInterface
                                          67229201 41.82% |   encoding/json.(*decodeState).literalInterface
----------------------------------------------------------+-------------
                                          14957349   100% |   go/parser.(*parser).parseStructType
   4709451  0.16% 83.47%   14957349   0.5%                | go/parser.(*parser).parseFieldDecl
----------------------------------------------------------+-------------
                                          46162139 97.40% |   github.com/go-openapi/spec.expandSchema
   4652214  0.16% 83.62%   47392150  1.59%                | github.com/go-openapi/spec.expandItems
                                          44372694 93.63% |   github.com/go-openapi/spec.expandSchema
----------------------------------------------------------+-------------
                                         294002592 99.46% |   go/parser.(*parser).parseFuncDecl
   4283604  0.14% 83.77%  295596438  9.89%                | go/parser.(*parser).parseBody
                                         290266363 98.20% |   go/parser.(*parser).parseStmtList
----------------------------------------------------------+-------------
                                          74066792   100% |   golang.org/x/tools/internal/imports.(*dirInfoCache).Store
   3515157  0.12% 83.89%   74066792  2.48%                | golang.org/x/tools/internal/imports.(*ModuleResolver).canonicalize
                                          22296125 30.10% |   os.(*File).readdir
                                          16312657 22.02% |   go/build.(*Context).matchFile
                                           4240051  5.72% |   go/parser.(*parser).parseFile
                                           3724727  5.03% |   os.newFile
----------------------------------------------------------+-------------
                                           4240051  0.77% |   golang.org/x/tools/internal/imports.(*ModuleResolver).canonicalize
   3053019   0.1% 83.99%  547486252 18.32%                | go/parser.(*parser).parseFile
                                         411905908 75.24% |   go/parser.(*parser).parseFuncDecl
                                          52787400  9.64% |   go/ast.NewScope (inline)
                                          36454243  6.66% |   go/parser.(*parser).parseGenDecl
                                          20828396  3.80% |   go/ast.NewObj (inline)
----------------------------------------------------------+-------------
   2847997 0.095% 84.08%  789356594 26.41%                | reflect.Value.call
                                         675301756 85.55% |   github.com/go-openapi/swag.ToGoName
                                         102256992 12.95% |   github.com/go-openapi/swag.(*splitter).gatherInitialismMatches
                                           6840831  0.87% |   github.com/go-openapi/swag.(*splitter).breakCasualString
----------------------------------------------------------+-------------
                                          16312657 83.64% |   golang.org/x/tools/internal/imports.(*ModuleResolver).canonicalize
   2262189 0.076% 84.16%   19503556  0.65%                | go/build.(*Context).matchFile
----------------------------------------------------------+-------------
                                          61603333 43.08% |   github.com/go-openapi/validate.(*objectValidator).Validate
   2143016 0.072% 84.23%  142996894  4.78%                | github.com/go-openapi/validate.NewSchemaValidator
                                         128119488 89.60% |   github.com/go-openapi/spec.expandSchema
----------------------------------------------------------+-------------
                                         675301756 85.76% |   reflect.Value.call
   2056455 0.069% 84.30%  787466685 26.35%                | github.com/go-openapi/swag.ToGoName
                                         710103241 90.18% |   github.com/go-openapi/swag.(*splitter).gatherInitialismMatches
                                          72577306  9.22% |   github.com/go-openapi/swag.(*splitter).breakCasualString
----------------------------------------------------------+-------------
                                          17728149 99.85% |   go/parser.(*parser).parseGenDecl
   1986257 0.066% 84.37%   17754355  0.59%                | go/parser.(*parser).parseStructType
                                          14957349 84.25% |   go/parser.(*parser).parseFieldDecl
----------------------------------------------------------+-------------
                                         184441918 94.91% |   github.com/go-openapi/spec.expandSchema
   1923042 0.064% 84.43%  194324372  6.50%                | github.com/go-openapi/spec.expandSchemaRef
                                          97998212 50.43% |   encoding/json.Unmarshal
                                          91956735 47.32% |   github.com/go-openapi/spec.expandSchema
                                          26125461 13.44% |   strings.(*Builder).WriteString (inline)
                                          24235602 12.47% |   github.com/go-openapi/spec.Schema.MarshalJSON
                                           5803825  2.99% |   github.com/go-openapi/jsonpointer.(*Pointer).parse
----------------------------------------------------------+-------------

## With spec using go-json


time go test -v -timeout 30m -parallel 3 hack/codegen_nonreg_test.go -args -fixture-file codegen-fixtures.yaml -skip-models -skip-full-flatten

real	4m3.087s
user	14m55.012s
sys	4m21.041s

pprof -tree -functions -sample_index=alloc_objects  prof/*/mem.pprof
File: swagger
Build ID: 0bb5d01437bcbb4ac7f5771c8e507402c3b13112
Type: alloc_objects
Time: Jan 9, 2024 at 4:56pm (CET)
Showing nodes accounting for 2544992097, 84.85% of 2999383944 total
Dropped 1922 nodes (cum <= 14996919)
Showing top 80 nodes out of 345
----------------------------------------------------------+-------------
      flat  flat%   sum%        cum   cum%   calls calls% + context 	 	 
----------------------------------------------------------+-------------
                                         702512477 81.41% |   github.com/go-openapi/swag.ToGoName
                                         100470698 11.64% |   reflect.Value.call
 862921570 28.77% 28.77%  862921570 28.77%                | github.com/go-openapi/swag.(*splitter).gatherInitialismMatches
----------------------------------------------------------+-------------
                                         291447197 70.87% |   github.com/goccy/go-json.unmarshal
                                         249681572 60.71% |   github.com/goccy/go-json/internal/decoder.(*interfaceDecoder).decodeEmptyInterface
                                         191379918 46.54% |   github.com/goccy/go-json/internal/decoder.(*structDecoder).Decode
 182058509  6.07% 34.84%  411243745 13.71%                | github.com/goccy/go-json/internal/decoder.(*mapDecoder).Decode
                                         270746751 65.84% |   github.com/goccy/go-json/internal/decoder.(*interfaceDecoder).decodeEmptyInterface
                                         179097077 43.55% |   github.com/go-openapi/spec.(*Schema).UnmarshalJSON
                                          43148616 10.49% |   github.com/goccy/go-json/internal/decoder.(*mapDecoder).mapassign
                                          37503530  9.12% |   reflect.makemap
----------------------------------------------------------+-------------
                                          59334797 50.33% |   github.com/go-openapi/swag.(*splitter).breakCasualString
 117882301  3.93% 38.77%  117882301  3.93%                | strings.(*Builder).grow
----------------------------------------------------------+-------------
                                          58990739 50.39% |   go/parser.(*parser).parseOperand
                                          22961420 19.61% |   go/parser.(*parser).parseTypeName
                                          12900671 11.02% |   go/parser.(*parser).parseSelector
                                          10846712  9.27% |   go/parser.(*parser).parseParameterList
                                           4215630  3.60% |   go/parser.(*parser).parseFuncDecl
                                           3458434  2.95% |   go/parser.(*parser).parseGenDecl
 112416887  3.75% 42.52%  117071279  3.90%                | go/parser.(*parser).parseIdent
----------------------------------------------------------+-------------
                                         101441974 71.45% |   path/filepath.readDir
                                          24219758 17.06% |   golang.org/x/tools/internal/imports.(*ModuleResolver).canonicalize
  88742709  2.96% 45.48%  141980144  4.73%                | os.(*File).readdir
                                          53222338 37.49% |   os.newUnixDirent
----------------------------------------------------------+-------------
                                         270746751 97.70% |   github.com/goccy/go-json/internal/decoder.(*mapDecoder).Decode
                                         127594484 46.04% |   github.com/goccy/go-json/internal/decoder.(*sliceDecoder).Decode
                                           9567790  3.45% |   github.com/goccy/go-json.unmarshal
  68276315  2.28% 47.75%  277130020  9.24%                | github.com/goccy/go-json/internal/decoder.(*interfaceDecoder).decodeEmptyInterface
                                         249681572 90.10% |   github.com/goccy/go-json/internal/decoder.(*mapDecoder).Decode
                                         139651019 50.39% |   github.com/goccy/go-json/internal/decoder.(*sliceDecoder).Decode
----------------------------------------------------------+-------------
                                           9205709 16.08% |   go/parser.(*parser).parseUnaryExpr
                                           5679735  9.92% |   go/parser.(*parser).parseParameters
                                           5520545  9.64% |   go/parser.(*parser).parseIfStmt
                                           4476059  7.82% |   go/parser.(*parser).parsePointerType
                                           4066952  7.10% |   go/parser.(*parser).parseCallOrConversion
                                           3178287  5.55% |   go/parser.(*parser).parseBlockStmt
                                           3082484  5.38% |   go/parser.(*parser).parseTypeName
  57262174  1.91% 49.66%   57262174  1.91%                | go/scanner.(*Scanner).scanIdentifier
----------------------------------------------------------+-------------
                                          34718615 64.93% |   encoding/gob.(*Encoder).sendActualType
                                          11455663 21.43% |   encoding/gob.(*Decoder).compileDec
  53466936  1.78% 51.44%   53466936  1.78%                | reflect.(*structType).Field
----------------------------------------------------------+-------------
                                          53286893   100% |   go/parser.(*parser).parseFile (inline)
  53313136  1.78% 53.22%   53313136  1.78%                | go/ast.NewScope
----------------------------------------------------------+-------------
                                          53222338   100% |   os.(*File).readdir
  53222338  1.77% 55.00%   53222338  1.77%                | os.newUnixDirent
----------------------------------------------------------+-------------
                                         155128457 78.11% |   go/parser.(*parser).parseSimpleStmt
                                          42257301 21.28% |   go/parser.(*parser).parseReturnStmt
  40694821  1.36% 56.35%  198599453  6.62%                | go/parser.(*parser).parseExprList
                                         156562033 78.83% |   go/parser.(*parser).parseBinaryExpr
----------------------------------------------------------+-------------
                                         120695800   100% |   encoding/gob.(*Decoder).compileDec
  37972301  1.27% 57.62%  120695800  4.02%                | encoding/gob.(*Decoder).decOpFor
                                         115491513 95.69% |   encoding/gob.(*Decoder).getDecEnginePtr
                                           6532837  5.41% |   encoding/gob.overflow (inline)
----------------------------------------------------------+-------------
                                          37503530 98.82% |   github.com/goccy/go-json/internal/decoder.(*mapDecoder).Decode
  37950234  1.27% 58.88%   37950234  1.27%                | reflect.makemap
----------------------------------------------------------+-------------
                                          43148616   100% |   github.com/goccy/go-json/internal/decoder.(*mapDecoder).Decode
  36938160  1.23% 60.12%   43148616  1.44%                | github.com/goccy/go-json/internal/decoder.(*mapDecoder).mapassign
----------------------------------------------------------+-------------
                                          27053421 81.04% |   github.com/go-openapi/spec.expandSchemaRef (inline)
  33383458  1.11% 61.23%   33383458  1.11%                | strings.(*Builder).WriteString
----------------------------------------------------------+-------------
                                          23311902 72.72% |   encoding/gob.(*Decoder).decodeArrayHelper
                                           7107139 22.17% |   encoding/gob.(*Decoder).recvType
  32054995  1.07% 62.30%   32054995  1.07%                | encoding/gob.decString
----------------------------------------------------------+-------------
                                          62772362   100% |   go/parser.(*parser).parseUnaryExpr
  27593708  0.92% 63.22%   62772362  2.09%                | go/parser.(*parser).parseCallOrConversion
                                          30046788 47.87% |   go/parser.(*parser).parseBinaryExpr
                                           4066952  6.48% |   go/scanner.(*Scanner).scanIdentifier
----------------------------------------------------------+-------------
                                         296456175 73.29% |   github.com/go-openapi/spec.(*Schema).UnmarshalJSON
                                         112878330 27.91% |   github.com/goccy/go-json/internal/decoder.(*ptrDecoder).Decode
                                          12497328  3.09% |   github.com/goccy/go-json/internal/decoder.(*sliceDecoder).Decode
                                           8691570  2.15% |   github.com/go-openapi/spec.(*StringOrArray).UnmarshalJSON
                                           5645254  1.40% |   github.com/go-openapi/validate.NewSchemaValidator
                                           5065144  1.25% |   github.com/go-openapi/spec.expandSchemaRef
  27361164  0.91% 64.13%  404472107 13.49%                | github.com/goccy/go-json.unmarshal
                                         303193647 74.96% |   github.com/goccy/go-json/internal/decoder.(*structDecoder).Decode
                                         291447197 72.06% |   github.com/goccy/go-json/internal/decoder.(*mapDecoder).Decode
                                         111560971 27.58% |   github.com/go-openapi/spec.(*Schema).UnmarshalJSON
                                           9567790  2.37% |   github.com/goccy/go-json/internal/decoder.(*interfaceDecoder).decodeEmptyInterface
----------------------------------------------------------+-------------
                                          25756981 97.32% |   encoding/gob.overflow (inline)
  26467449  0.88% 65.01%   26467449  0.88%                | errors.New
----------------------------------------------------------+-------------
                                          45019046 87.33% |   encoding/gob.(*Decoder).compileDec (inline)
                                           6532837 12.67% |   encoding/gob.(*Decoder).decOpFor (inline)
  25794902  0.86% 65.87%   51551883  1.72%                | encoding/gob.overflow
                                          25756981 49.96% |   errors.New (inline)
----------------------------------------------------------+-------------
                                          13502906 54.11% |   path/filepath.readDir
                                           4452783 17.84% |   golang.org/x/tools/internal/imports.(*ModuleResolver).canonicalize
  24954365  0.83% 66.70%   24954365  0.83%                | os.newFile
----------------------------------------------------------+-------------
                                          11315260 38.40% |   github.com/go-openapi/analysis.(*Spec).analyzeSchema
  23568675  0.79% 67.49%   29469059  0.98%                | path.Join
----------------------------------------------------------+-------------
                                           8971419 41.59% |   os.statNolog
                                           6685958 31.00% |   path/filepath.readDir
  21569175  0.72% 68.21%   21569175  0.72%                | syscall.ByteSliceFromString
----------------------------------------------------------+-------------
                                          20925982 99.94% |   go/parser.(*parser).parseFile (inline)
  20938362   0.7% 68.91%   20938362   0.7%                | go/ast.NewObj
----------------------------------------------------------+-------------
                                          76245667   100% |   go/parser.(*parser).parseParameters
  20878259   0.7% 69.60%   76245667  2.54%                | go/parser.(*parser).parseParameterList
                                          16603072 21.78% |   go/parser.(*parser).parseParameterList.func1 (inline)
                                          14710563 19.29% |   go/parser.(*parser).parsePointerType
                                          10846712 14.23% |   go/parser.(*parser).parseIdent
                                          10794075 14.16% |   go/parser.(*parser).parseTypeName
----------------------------------------------------------+-------------
                                         289896482 98.80% |   go/parser.(*parser).parseBody
                                         131268902 44.74% |   go/parser.(*parser).parseBlockStmt
                                           3191576  1.09% |   go/parser.(*parser).parseStmt
  20610283  0.69% 70.29%  293418930  9.78%                | go/parser.(*parser).parseStmtList
                                         283609642 96.66% |   go/parser.(*parser).parseStmt
----------------------------------------------------------+-------------
                                         139651019 62.08% |   github.com/goccy/go-json/internal/decoder.(*interfaceDecoder).decodeEmptyInterface
                                         117755836 52.35% |   github.com/goccy/go-json/internal/decoder.(*structDecoder).Decode
  19628719  0.65% 70.95%  224937763  7.50%                | github.com/goccy/go-json/internal/decoder.(*sliceDecoder).Decode
                                         127594484 56.72% |   github.com/goccy/go-json/internal/decoder.(*interfaceDecoder).decodeEmptyInterface
                                         101830467 45.27% |   github.com/go-openapi/spec.(*Schema).UnmarshalJSON
                                          12497328  5.56% |   github.com/goccy/go-json.unmarshal
                                           8174217  3.63% |   reflect.unsafe_NewArray
----------------------------------------------------------+-------------
                                         115491513 87.28% |   encoding/gob.(*Decoder).decOpFor
                                          34097983 25.77% |   encoding/gob.(*Decoder).recvType
  18383050  0.61% 71.56%  132319047  4.41%                | encoding/gob.(*Decoder).getDecEnginePtr
                                         129086181 97.56% |   encoding/gob.(*Decoder).compileDec
----------------------------------------------------------+-------------
                                          44881377 85.24% |   github.com/go-openapi/spec.Schema.MarshalJSON (inline)
  18140642   0.6% 72.16%   52650447  1.76%                | github.com/goccy/go-json.MarshalWithOption
                                          42395132 80.52% |   github.com/go-openapi/spec.Schema.MarshalJSON
----------------------------------------------------------+-------------
                                         109285655 61.98% |   go/parser.(*parser).parseIfStmt
                                          67072852 38.04% |   go/parser.(*parser).parseStmt
  17885131   0.6% 72.76%  176323302  5.88%                | go/parser.(*parser).parseSimpleStmt
                                         155128457 87.98% |   go/parser.(*parser).parseExprList
----------------------------------------------------------+-------------
                                          16603072   100% |   go/parser.(*parser).parseParameterList (inline)
  16603072  0.55% 73.31%   16603072  0.55%                | go/parser.(*parser).parseParameterList.func1
----------------------------------------------------------+-------------
                                           8174217 50.59% |   github.com/goccy/go-json/internal/decoder.(*sliceDecoder).Decode
                                           7982858 49.41% |   reflect.MakeSlice
  16157075  0.54% 73.85%   16157075  0.54%                | reflect.unsafe_NewArray
----------------------------------------------------------+-------------
  16062423  0.54% 74.39%   25033842  0.83%                | os.statNolog
                                           8971419 35.84% |   syscall.ByteSliceFromString
----------------------------------------------------------+-------------
                                         179097077 59.65% |   github.com/goccy/go-json/internal/decoder.(*mapDecoder).Decode
                                         111560971 37.16% |   github.com/goccy/go-json.unmarshal
                                         108557696 36.16% |   github.com/go-openapi/spec.expandSchemaRef
                                         101830467 33.92% |   github.com/goccy/go-json/internal/decoder.(*sliceDecoder).Decode
                                           6658319  2.22% |   github.com/goccy/go-json/internal/decoder.(*ptrDecoder).Decode
  14396603  0.48% 74.87%  300233253 10.01%                | github.com/go-openapi/spec.(*Schema).UnmarshalJSON
                                         296456175 98.74% |   github.com/goccy/go-json.unmarshal
                                           3012425  1.00% |   github.com/go-openapi/jsonpointer.(*Pointer).parse
----------------------------------------------------------+-------------
                                          15909065 70.30% |   io.WriteString
  13507901  0.45% 75.32%   22629915  0.75%                | github.com/kr/text.(*indentWriter).Write
----------------------------------------------------------+-------------
                                          71734712 89.58% |   github.com/go-openapi/swag.ToGoName
                                           6743574  8.42% |   reflect.Value.call
  13378328  0.45% 75.76%   80078624  2.67%                | github.com/go-openapi/swag.(*splitter).breakCasualString
                                          59334797 74.10% |   strings.(*Builder).grow
----------------------------------------------------------+-------------
                                          25772419   100% |   go/parser.(*parser).parseUnaryExpr
  12871748  0.43% 76.19%   25772419  0.86%                | go/parser.(*parser).parseSelector
                                          12900671 50.06% |   go/parser.(*parser).parseIdent
----------------------------------------------------------+-------------
  12408371  0.41% 76.61%   51854637  1.73%                | encoding/gob.(*Encoder).sendActualType
                                          34718615 66.95% |   reflect.(*structType).Field
----------------------------------------------------------+-------------
  12031242   0.4% 77.01%  106133622  3.54%                | encoding/gob.(*Decoder).recvType
                                          34097983 32.13% |   encoding/gob.(*Decoder).getDecEnginePtr
                                          28772595 27.11% |   encoding/gob.(*Decoder).decodeArrayHelper
                                          12681112 11.95% |   reflect.MakeSlice
                                           7107139  6.70% |   encoding/gob.decString
----------------------------------------------------------+-------------
                                          42395132 86.81% |   github.com/goccy/go-json.MarshalWithOption
                                          10028446 20.54% |   github.com/go-openapi/spec.expandSchemaRef
  11821422  0.39% 77.40%   48834517  1.63%                | github.com/go-openapi/spec.Schema.MarshalJSON
                                          44881377 91.91% |   github.com/goccy/go-json.MarshalWithOption (inline)
----------------------------------------------------------+-------------
                                         129086181   100% |   encoding/gob.(*Decoder).getDecEnginePtr
  10958029  0.37% 77.77%  129086181  4.30%                | encoding/gob.(*Decoder).compileDec
                                         120695800 93.50% |   encoding/gob.(*Decoder).decOpFor
                                          45019046 34.88% |   encoding/gob.overflow (inline)
                                          11455663  8.87% |   reflect.(*structType).Field
----------------------------------------------------------+-------------
                                          21139786 98.89% |   go/parser.(*parser).parseFuncDecl
   9734060  0.32% 78.09%   21377160  0.71%                | go/parser.(*parser).parseResult
                                           7096380 33.20% |   go/parser.(*parser).parseParameters
----------------------------------------------------------+-------------
   9552396  0.32% 78.41%   99611557  3.32%                | golang.org/x/tools/internal/imports.(*dirInfoCache).Store
                                          81876193 82.20% |   golang.org/x/tools/internal/imports.(*ModuleResolver).canonicalize
----------------------------------------------------------+-------------
   9299154  0.31% 78.72%   36311615  1.21%                | github.com/go-openapi/analysis.(*Spec).analyzeSchema
                                          11315260 31.16% |   path.Join
                                           7198741 19.82% |   github.com/go-openapi/jsonpointer.(*Pointer).parse
----------------------------------------------------------+-------------
                                          53362164   100% |   go/parser.(*parser).parseStmt
   9276741  0.31% 79.03%   53372659  1.78%                | go/parser.(*parser).parseReturnStmt
                                          42257301 79.17% |   go/parser.(*parser).parseExprList
----------------------------------------------------------+-------------
   8747997  0.29% 79.32%  121364090  4.05%                | github.com/go-openapi/validate.(*objectValidator).Validate
                                          71597586 58.99% |   github.com/go-openapi/validate.NewSchemaValidator
                                          22996074 18.95% |   regexp.compile
----------------------------------------------------------+-------------
                                          17422635   100% |   github.com/goccy/go-json/internal/decoder.(*structDecoder).Decode
   8731065  0.29% 79.61%   17422635  0.58%                | github.com/go-openapi/spec.(*StringOrArray).UnmarshalJSON
                                           8691570 49.89% |   github.com/goccy/go-json.unmarshal
----------------------------------------------------------+-------------
                                         412356454 99.25% |   go/parser.(*parser).parseFile
   8521563  0.28% 79.90%  415471756 13.85%                | go/parser.(*parser).parseFuncDecl
                                         293646391 70.68% |   go/parser.(*parser).parseBody
                                          82578163 19.88% |   go/parser.(*parser).parseParameters
                                          21139786  5.09% |   go/parser.(*parser).parseResult
                                           4215630  1.01% |   go/parser.(*parser).parseIdent
----------------------------------------------------------+-------------
                                           7198741 42.74% |   github.com/go-openapi/analysis.(*Spec).analyzeSchema
                                           5968068 35.43% |   github.com/go-openapi/spec.expandSchemaRef
                                           3012425 17.88% |   github.com/go-openapi/spec.(*Schema).UnmarshalJSON
   8427205  0.28% 80.18%   16844706  0.56%                | github.com/go-openapi/jsonpointer.(*Pointer).parse
----------------------------------------------------------+-------------
                                         131961123 94.86% |   go/parser.(*parser).parseIfStmt
                                          24194644 17.39% |   go/parser.(*parser).parseStmt
   8327747  0.28% 80.45%  139107492  4.64%                | go/parser.(*parser).parseBlockStmt
                                         131268902 94.37% |   go/parser.(*parser).parseStmtList
                                           3178287  2.28% |   go/scanner.(*Scanner).scanIdentifier
----------------------------------------------------------+-------------
                                          82578163 91.67% |   go/parser.(*parser).parseFuncDecl
                                           7096380  7.88% |   go/parser.(*parser).parseResult
   8153473  0.27% 80.73%   90078802  3.00%                | go/parser.(*parser).parseParameters
                                          76245667 84.64% |   go/parser.(*parser).parseParameterList
                                           5679735  6.31% |   go/scanner.(*Scanner).scanIdentifier
----------------------------------------------------------+-------------
                                          12681112 79.51% |   encoding/gob.(*Decoder).recvType
   7965854  0.27% 80.99%   15948712  0.53%                | reflect.MakeSlice
                                           7982858 50.05% |   reflect.unsafe_NewArray
----------------------------------------------------------+-------------
                                         216393381   100% |   go/parser.(*parser).parseStmt
   7768641  0.26% 81.25%  216393765  7.21%                | go/parser.(*parser).parseIfStmt
                                         131961123 60.98% |   go/parser.(*parser).parseBlockStmt
                                         109285655 50.50% |   go/parser.(*parser).parseSimpleStmt
                                           5520545  2.55% |   go/scanner.(*Scanner).scanIdentifier
----------------------------------------------------------+-------------
                                         156562033 99.84% |   go/parser.(*parser).parseExprList
                                          30046788 19.16% |   go/parser.(*parser).parseCallOrConversion
   7597976  0.25% 81.50%  156815871  5.23%                | go/parser.(*parser).parseBinaryExpr
                                         148646644 94.79% |   go/parser.(*parser).parseUnaryExpr
----------------------------------------------------------+-------------
   7523269  0.25% 81.76%   23559901  0.79%                | io.WriteString
                                          15909065 67.53% |   github.com/kr/text.(*indentWriter).Write
----------------------------------------------------------+-------------
                                          67622677   100% |   go/parser.(*parser).parseUnaryExpr
   7326469  0.24% 82.00%   67622677  2.25%                | go/parser.(*parser).parseOperand
                                          58990739 87.24% |   go/parser.(*parser).parseIdent
----------------------------------------------------------+-------------
                                          28772595 83.26% |   encoding/gob.(*Decoder).recvType
   7192646  0.24% 82.24%   34558587  1.15%                | encoding/gob.(*Decoder).decodeArrayHelper
                                          23311902 67.46% |   encoding/gob.decString
----------------------------------------------------------+-------------
   6717471  0.22% 82.46%  137137009  4.57%                | path/filepath.readDir
                                         101441974 73.97% |   os.(*File).readdir
                                          13502906  9.85% |   os.newFile
                                           6685958  4.88% |   syscall.ByteSliceFromString
----------------------------------------------------------+-------------
                                          14710563 60.69% |   go/parser.(*parser).parseParameterList
                                           5436427 22.43% |   go/parser.(*parser).parseUnaryExpr
   6529707  0.22% 82.68%   24236928  0.81%                | go/parser.(*parser).parsePointerType
                                          13231077 54.59% |   go/parser.(*parser).parseTypeName
                                           4476059 18.47% |   go/scanner.(*Scanner).scanIdentifier
----------------------------------------------------------+-------------
                                          37115714 73.14% |   go/parser.(*parser).parseFile
                                          13627969 26.86% |   go/parser.(*parser).parseStmt
   6359311  0.21% 82.89%   50745094  1.69%                | go/parser.(*parser).parseGenDecl
                                          17970727 35.41% |   go/parser.(*parser).parseStructType
                                           3458434  6.82% |   go/parser.(*parser).parseIdent
----------------------------------------------------------+-------------
                                         141001603 73.63% |   github.com/go-openapi/validate.NewSchemaValidator
                                          89514266 46.74% |   github.com/go-openapi/spec.expandSchemaRef
                                          47441874 24.77% |   github.com/go-openapi/spec.expandItems
   5380657  0.18% 83.07%  191512478  6.39%                | github.com/go-openapi/spec.expandSchema
                                         190065654 99.24% |   github.com/go-openapi/spec.expandSchemaRef
                                          49987062 26.10% |   github.com/go-openapi/spec.expandItems
----------------------------------------------------------+-------------
                                          13231077 42.34% |   go/parser.(*parser).parsePointerType
                                          10794075 34.54% |   go/parser.(*parser).parseParameterList
   5204969  0.17% 83.25%   31248873  1.04%                | go/parser.(*parser).parseTypeName
                                          22961420 73.48% |   go/parser.(*parser).parseIdent
                                           3082484  9.86% |   go/scanner.(*Scanner).scanIdentifier
----------------------------------------------------------+-------------
                                          15242923   100% |   go/parser.(*parser).parseStructType
   4802458  0.16% 83.41%   15242923  0.51%                | go/parser.(*parser).parseFieldDecl
----------------------------------------------------------+-------------
                                          49987062 98.97% |   github.com/go-openapi/spec.expandSchema
   4661458  0.16% 83.56%   50506315  1.68%                | github.com/go-openapi/spec.expandItems
                                          47441874 93.93% |   github.com/go-openapi/spec.expandSchema
----------------------------------------------------------+-------------
                                         293646391 99.48% |   go/parser.(*parser).parseFuncDecl
   4215809  0.14% 83.70%  295176936  9.84%                | go/parser.(*parser).parseBody
                                         289896482 98.21% |   go/parser.(*parser).parseStmtList
----------------------------------------------------------+-------------
                                          81876193   100% |   golang.org/x/tools/internal/imports.(*dirInfoCache).Store
   3671695  0.12% 83.82%   81876193  2.73%                | golang.org/x/tools/internal/imports.(*ModuleResolver).canonicalize
                                          24219758 29.58% |   os.(*File).readdir
                                          16750152 20.46% |   go/build.(*Context).matchFile
                                           5523528  6.75% |   go/parser.(*parser).parseFile
                                           4452783  5.44% |   os.newFile
----------------------------------------------------------+-------------
                                         115232986 98.41% |   github.com/goccy/go-json/internal/decoder.(*structDecoder).Decode
   3512667  0.12% 83.94%  117099151  3.90%                | github.com/goccy/go-json/internal/decoder.(*ptrDecoder).Decode
                                         112878330 96.40% |   github.com/goccy/go-json.unmarshal
                                           6658319  5.69% |   github.com/go-openapi/spec.(*Schema).UnmarshalJSON
----------------------------------------------------------+-------------
                                           5523528  1.00% |   golang.org/x/tools/internal/imports.(*ModuleResolver).canonicalize
   3313715  0.11% 84.05%  550563921 18.36%                | go/parser.(*parser).parseFile
                                         412356454 74.90% |   go/parser.(*parser).parseFuncDecl
                                          53286893  9.68% |   go/ast.NewScope (inline)
                                          37115714  6.74% |   go/parser.(*parser).parseGenDecl
                                          20925982  3.80% |   go/ast.NewObj (inline)
----------------------------------------------------------+-------------
   2816587 0.094% 84.15%  779785041 26.00%                | reflect.Value.call
                                         667665230 85.62% |   github.com/go-openapi/swag.ToGoName
                                         100470698 12.88% |   github.com/go-openapi/swag.(*splitter).gatherInitialismMatches
                                           6743574  0.86% |   github.com/go-openapi/swag.(*splitter).breakCasualString
----------------------------------------------------------+-------------
                                         303193647 99.35% |   github.com/goccy/go-json.unmarshal
   2584132 0.086% 84.23%  305180004 10.17%                | github.com/goccy/go-json/internal/decoder.(*structDecoder).Decode
                                         191379918 62.71% |   github.com/goccy/go-json/internal/decoder.(*mapDecoder).Decode
                                         117755836 38.59% |   github.com/goccy/go-json/internal/decoder.(*sliceDecoder).Decode
                                         115232986 37.76% |   github.com/goccy/go-json/internal/decoder.(*ptrDecoder).Decode
                                          17422635  5.71% |   github.com/go-openapi/spec.(*StringOrArray).UnmarshalJSON
----------------------------------------------------------+-------------
                                          16750152 84.11% |   golang.org/x/tools/internal/imports.(*ModuleResolver).canonicalize
   2306296 0.077% 84.31%   19914483  0.66%                | go/build.(*Context).matchFile
----------------------------------------------------------+-------------
                                          71597586 45.13% |   github.com/go-openapi/validate.(*objectValidator).Validate
                                          55631729 35.07% |   github.com/go-openapi/validate.newSchemaPropsValidator
   2133870 0.071% 84.38%  158637340  5.29%                | github.com/go-openapi/validate.NewSchemaValidator
                                         141001603 88.88% |   github.com/go-openapi/spec.expandSchema
                                          50280560 31.70% |   github.com/go-openapi/validate.newSchemaPropsValidator
                                           5645254  3.56% |   github.com/goccy/go-json.unmarshal
----------------------------------------------------------+-------------
                                         667665230 85.71% |   reflect.Value.call
                                          62555650  8.03% |   github.com/go-swagger/go-swagger/generator.typeResolver.knownDefGoType
   2026283 0.068% 84.45%  778980887 25.97%                | github.com/go-openapi/swag.ToGoName
                                         702512477 90.18% |   github.com/go-openapi/swag.(*splitter).gatherInitialismMatches
                                          71734712  9.21% |   github.com/go-openapi/swag.(*splitter).breakCasualString
----------------------------------------------------------+-------------
                                         190065654 95.76% |   github.com/go-openapi/spec.expandSchema
   2005913 0.067% 84.51%  198478386  6.62%                | github.com/go-openapi/spec.expandSchemaRef
                                         108557696 54.69% |   github.com/go-openapi/spec.(*Schema).UnmarshalJSON
                                          89514266 45.10% |   github.com/go-openapi/spec.expandSchema
                                          27053421 13.63% |   strings.(*Builder).WriteString (inline)
                                          10028446  5.05% |   github.com/go-openapi/spec.Schema.MarshalJSON
                                           5968068  3.01% |   github.com/go-openapi/jsonpointer.(*Pointer).parse
                                           5065144  2.55% |   github.com/goccy/go-json.unmarshal
----------------------------------------------------------+-------------
                                          17970727 99.76% |   go/parser.(*parser).parseGenDecl
   1991208 0.066% 84.58%   18014500   0.6%                | go/parser.(*parser).parseStructType
                                          15242923 84.61% |   go/parser.(*parser).parseFieldDecl
----------------------------------------------------------+-------------
                                         148646644   100% |   go/parser.(*parser).parseBinaryExpr
   1818645 0.061% 84.64%  148646644  4.96%                | go/parser.(*parser).parseUnaryExpr
                                          67622677 45.49% |   go/parser.(*parser).parseOperand
                                          62772362 42.23% |   go/parser.(*parser).parseCallOrConversion
                                          25772419 17.34% |   go/parser.(*parser).parseSelector
                                           9205709  6.19% |   go/scanner.(*Scanner).scanIdentifier
                                           5436427  3.66% |   go/parser.(*parser).parsePointerType
----------------------------------------------------------+-------------
                                          50280560 88.48% |   github.com/go-openapi/validate.NewSchemaValidator
   1622667 0.054% 84.70%   56828421  1.89%                | github.com/go-openapi/validate.newSchemaPropsValidator
                                          55631729 97.89% |   github.com/go-openapi/validate.NewSchemaValidator
----------------------------------------------------------+-------------
   1608727 0.054% 84.75%   64229063  2.14%                | github.com/go-swagger/go-swagger/generator.typeResolver.knownDefGoType
                                          62555650 97.39% |   github.com/go-openapi/swag.ToGoName
----------------------------------------------------------+-------------
                                         283609642   100% |   go/parser.(*parser).parseStmtList
   1551736 0.052% 84.80%  283618574  9.46%                | go/parser.(*parser).parseStmt
                                         216393381 76.30% |   go/parser.(*parser).parseIfStmt
                                          67072852 23.65% |   go/parser.(*parser).parseSimpleStmt
                                          53362164 18.81% |   go/parser.(*parser).parseReturnStmt
                                          24194644  8.53% |   go/parser.(*parser).parseBlockStmt
                                          13627969  4.81% |   go/parser.(*parser).parseGenDecl
                                           3191576  1.13% |   go/parser.(*parser).parseStmtList
----------------------------------------------------------+-------------
                                          22996074 99.68% |   github.com/go-openapi/validate.(*objectValidator).Validate
   1482918 0.049% 84.85%   23069777  0.77%                | regexp.compile
----------------------------------------------------------+-------------

pprof -tree -functions -sample_index=alloc_space  prof/*/mem.pprof
File: swagger
Build ID: 0bb5d01437bcbb4ac7f5771c8e507402c3b13112
Type: alloc_space
Time: Jan 9, 2024 at 4:56pm (CET)
Showing nodes accounting for 169.38GB, 82.96% of 204.17GB total
Dropped 1911 nodes (cum <= 1.02GB)
Showing top 80 nodes out of 356
----------------------------------------------------------+-------------
      flat  flat%   sum%        cum   cum%   calls calls% + context 	 	 
----------------------------------------------------------+-------------
   13.38GB  6.55%  6.55%    13.38GB  6.55%                | github.com/go-openapi/swag.(*splitter).gatherInitialismMatches
----------------------------------------------------------+-------------
                                           13.87GB   100% |   github.com/goccy/go-json/internal/decoder.(*mapDecoder).Decode
   10.75GB  5.26% 11.82%    13.87GB  6.79%                | github.com/goccy/go-json/internal/decoder.(*mapDecoder).mapassign
                                            3.12GB 22.52% |   reflect.mapassign0
----------------------------------------------------------+-------------
                                            9.23GB 82.24% |   go/parser.ParseFile
                                            1.78GB 15.89% |   golang.org/x/tools/internal/gopathwalk.(*walker).walk
   10.22GB  5.01% 16.82%    11.22GB  5.50%                | os.ReadFile
                                            0.26GB  2.30% |   os.newFile
                                            0.23GB  2.01% |   syscall.ByteSliceFromString
----------------------------------------------------------+-------------
                                            8.97GB   100% |   bytes.(*Buffer).grow
    8.97GB  4.39% 21.22%     8.97GB  4.39%                | bytes.growSlice
----------------------------------------------------------+-------------
                                           30.10GB 77.54% |   github.com/go-openapi/spec.(*Schema).UnmarshalJSON
                                            7.68GB 19.79% |   github.com/goccy/go-json/internal/decoder.(*ptrDecoder).Decode
                                            4.60GB 11.85% |   github.com/go-openapi/spec.(*SchemaOrArray).UnmarshalJSON
                                            0.78GB  2.01% |   github.com/goccy/go-json/internal/decoder.(*sliceDecoder).Decode
                                            0.72GB  1.85% |   github.com/go-openapi/validate.NewSchemaValidator
                                            0.27GB  0.69% |   github.com/go-openapi/spec.expandSchemaRef
    6.20GB  3.04% 24.26%    38.82GB 19.01%                | github.com/goccy/go-json.unmarshal
                                           32.88GB 84.70% |   github.com/goccy/go-json/internal/decoder.(*mapDecoder).Decode
                                           12.17GB 31.34% |   github.com/go-openapi/spec.(*Schema).UnmarshalJSON
                                           11.62GB 29.92% |   github.com/goccy/go-json/internal/decoder.(*sliceDecoder).Decode
                                           10.66GB 27.47% |   github.com/goccy/go-json/internal/decoder.(*ptrDecoder).Decode
                                            0.33GB  0.84% |   github.com/goccy/go-json/internal/decoder.(*interfaceDecoder).decodeEmptyInterface
----------------------------------------------------------+-------------
                                           32.88GB 94.69% |   github.com/goccy/go-json.unmarshal
                                           12.66GB 36.47% |   github.com/goccy/go-json/internal/decoder.(*interfaceDecoder).decodeEmptyInterface
    6.08GB  2.98% 27.23%    34.72GB 17.01%                | github.com/goccy/go-json/internal/decoder.(*mapDecoder).Decode
                                           19.83GB 57.11% |   github.com/go-openapi/spec.(*Schema).UnmarshalJSON
                                           13.87GB 39.95% |   github.com/goccy/go-json/internal/decoder.(*mapDecoder).mapassign
                                           13.04GB 37.54% |   github.com/goccy/go-json/internal/decoder.(*interfaceDecoder).decodeEmptyInterface
                                            1.68GB  4.84% |   reflect.makemap
----------------------------------------------------------+-------------
    6.03GB  2.95% 30.18%     6.03GB  2.95%                | encoding/gob.(*encBuffer).Write
----------------------------------------------------------+-------------
                                            0.88GB 15.17% |   github.com/go-openapi/swag.(*splitter).breakCasualString
                                            0.60GB 10.31% |   golang.org/x/tools/internal/gopathwalk.(*walker).walk
                                            0.23GB  3.88% |   golang.org/x/tools/internal/imports.(*ModuleResolver).canonicalize
    5.83GB  2.86% 33.04%     5.83GB  2.86%                | strings.(*Builder).grow
----------------------------------------------------------+-------------
                                            4.51GB 99.69% |   go/build.(*Context).matchFile (inline)
    4.52GB  2.21% 35.26%     4.52GB  2.21%                | bufio.NewReaderSize
----------------------------------------------------------+-------------
                                            4.22GB   100% |   github.com/go-openapi/spec.(*Schema).UnmarshalJSON
    4.22GB  2.07% 37.32%     4.22GB  2.07%                | github.com/go-openapi/swag.(*NameProvider).GetJSONNames
----------------------------------------------------------+-------------
                                           19.83GB 63.87% |   github.com/goccy/go-json/internal/decoder.(*mapDecoder).Decode
                                           12.97GB 41.78% |   github.com/go-openapi/spec.expandSchemaRef
                                           12.17GB 39.18% |   github.com/goccy/go-json.unmarshal
                                            9.91GB 31.93% |   github.com/goccy/go-json/internal/decoder.(*sliceDecoder).Decode
                                            0.76GB  2.45% |   github.com/goccy/go-json/internal/decoder.(*ptrDecoder).Decode
    3.90GB  1.91% 39.23%    31.05GB 15.21%                | github.com/go-openapi/spec.(*Schema).UnmarshalJSON
                                           30.10GB 96.94% |   github.com/goccy/go-json.unmarshal
                                            4.22GB 13.60% |   github.com/go-openapi/swag.(*NameProvider).GetJSONNames
----------------------------------------------------------+-------------
                                            1.15GB 16.38% |   golang.org/x/tools/internal/imports.(*ModuleResolver).canonicalize
    3.82GB  1.87% 41.10%     7.03GB  3.44%                | os.(*File).readdir
                                            3.17GB 45.11% |   os.newUnixDirent
----------------------------------------------------------+-------------
    3.76GB  1.84% 42.95%     3.76GB  1.84%                | internal/saferio.ReadData
----------------------------------------------------------+-------------
                                            3.12GB 88.37% |   github.com/goccy/go-json/internal/decoder.(*mapDecoder).mapassign
    3.54GB  1.73% 44.68%     3.54GB  1.73%                | reflect.mapassign0
----------------------------------------------------------+-------------
                                            1.98GB 53.28% |   go/parser.(*parser).parseOperand
                                            0.67GB 17.97% |   go/parser.(*parser).parseParameterList
                                            0.46GB 12.34% |   go/parser.(*parser).parseBinaryExpr
                                            0.34GB  9.19% |   go/parser.(*parser).parseGenDecl
                                            0.22GB  5.97% |   go/parser.(*parser).parseFuncDecl
    3.35GB  1.64% 46.32%     3.72GB  1.82%                | go/parser.(*parser).parseIdent
                                            0.28GB  7.62% |   go/token.(*File).AddLine
----------------------------------------------------------+-------------
                                            3.17GB   100% |   os.(*File).readdir
    3.17GB  1.55% 47.87%     3.17GB  1.55%                | os.newUnixDirent
----------------------------------------------------------+-------------
                                            8.63GB 74.61% |   github.com/goccy/go-json/internal/encoder.AppendMarshalJSON
                                            3.81GB 32.95% |   github.com/go-openapi/spec.expandSchemaRef
       3GB  1.47% 49.34%    11.56GB  5.66%                | github.com/go-openapi/spec.Schema.MarshalJSON
                                           10.71GB 92.63% |   github.com/goccy/go-json.MarshalWithOption (inline)
                                            1.34GB 11.58% |   bytes.(*Buffer).grow
----------------------------------------------------------+-------------
                                           10.71GB 75.05% |   github.com/go-openapi/spec.Schema.MarshalJSON (inline)
                                           10.12GB 70.88% |   github.com/goccy/go-json/internal/encoder.AppendMarshalJSON (inline)
    2.96GB  1.45% 50.79%    14.27GB  6.99%                | github.com/goccy/go-json.MarshalWithOption
                                           13.51GB 94.68% |   github.com/goccy/go-json/internal/encoder.AppendMarshalJSON
----------------------------------------------------------+-------------
                                           20.09GB 79.75% |   github.com/go-openapi/validate.NewSchemaValidator
                                           11.85GB 47.05% |   github.com/go-openapi/spec.expandSchemaRef
                                            5.61GB 22.27% |   github.com/go-openapi/spec.expandItems
    2.92GB  1.43% 52.22%    25.19GB 12.34%                | github.com/go-openapi/spec.expandSchema
                                           24.74GB 98.22% |   github.com/go-openapi/spec.expandSchemaRef
                                            7.19GB 28.52% |   github.com/go-openapi/spec.expandItems
----------------------------------------------------------+-------------
                                           13.47GB 76.46% |   github.com/go-openapi/validate.(*schemaPropsValidator).Validate
                                            3.55GB 20.17% |   github.com/go-openapi/validate.(*objectValidator).validatePatternProperty
    2.56GB  1.26% 53.48%    17.62GB  8.63%                | github.com/go-openapi/validate.(*objectValidator).Validate
                                           10.95GB 62.15% |   github.com/go-openapi/validate.NewSchemaValidator
                                            7.80GB 44.29% |   github.com/go-openapi/validate.(*schemaPropsValidator).Validate
                                               6GB 34.06% |   github.com/go-openapi/validate.(*objectValidator).validatePatternProperty
                                            0.99GB  5.59% |   regexp.compile
----------------------------------------------------------+-------------
    2.55GB  1.25% 54.73%     2.55GB  1.25%                | github.com/go-openapi/analysis.(*Spec).reset
----------------------------------------------------------+-------------
    2.54GB  1.24% 55.97%     4.85GB  2.38%                | github.com/go-openapi/analysis.(*Spec).analyzeSchema
                                            0.65GB 13.34% |   path.Join
                                            0.48GB  9.96% |   net/url.parse
----------------------------------------------------------+-------------
                                            7.19GB 98.79% |   github.com/go-openapi/spec.expandSchema
    2.53GB  1.24% 57.21%     7.27GB  3.56%                | github.com/go-openapi/spec.expandItems
                                            5.61GB 77.13% |   github.com/go-openapi/spec.expandSchema
----------------------------------------------------------+-------------
                                            2.31GB   100% |   go/build.(*Context).matchFile
    2.31GB  1.13% 58.34%     2.31GB  1.13%                | go/build.(*importReader).readByte
----------------------------------------------------------+-------------
                                            2.14GB   100% |   github.com/goccy/go-json/internal/encoder.AppendMarshalJSON
    2.14GB  1.05% 59.39%     2.14GB  1.05%                | github.com/go-openapi/spec.SchemaProperties.ToOrderedSchemaItems
----------------------------------------------------------+-------------
                                            2.40GB 85.93% |   golang.org/x/tools/internal/gopathwalk.(*walker).walk
                                            0.30GB 10.87% |   golang.org/x/tools/internal/imports.(*dirInfoCache).Store
    2.13GB  1.05% 60.44%     2.80GB  1.37%                | os.statNolog
                                            0.66GB 23.67% |   syscall.ByteSliceFromString
----------------------------------------------------------+-------------
    2.11GB  1.03% 61.47%     2.11GB  1.03%                | text/template.addValueFuncs
----------------------------------------------------------+-------------
                                            1.95GB 99.90% |   go/parser.(*parser).parseFile (inline)
    1.96GB  0.96% 62.43%     1.96GB  0.96%                | go/ast.(*Scope).Insert
----------------------------------------------------------+-------------
                                            1.95GB   100% |   github.com/goccy/go-json/internal/encoder.compactObject
    1.95GB  0.96% 63.39%     1.95GB  0.96%                | github.com/goccy/go-json/internal/encoder.compactString
----------------------------------------------------------+-------------
                                            0.36GB 20.27% |   go/parser.(*parser).parseFuncDecl
                                            0.35GB 19.54% |   go/parser.(*parser).parseGenDecl
                                            0.34GB 19.23% |   go/parser.(*parser).parseBlockStmt
                                            0.28GB 15.82% |   go/parser.(*parser).parseIdent
    1.79GB  0.88% 64.26%     1.79GB  0.88%                | go/token.(*File).AddLine
----------------------------------------------------------+-------------
                                            1.68GB 96.02% |   github.com/goccy/go-json/internal/decoder.(*mapDecoder).Decode
    1.75GB  0.86% 65.12%     1.75GB  0.86%                | reflect.makemap
----------------------------------------------------------+-------------
    1.73GB  0.85% 65.97%     5.13GB  2.51%                | encoding/gob.(*Encoder).sendActualType
                                            2.27GB 44.14% |   bytes.(*Buffer).grow
----------------------------------------------------------+-------------
                                            0.86GB 50.22% |   github.com/goccy/go-json/internal/decoder.(*sliceDecoder).Decode
                                            0.75GB 43.72% |   encoding/gob.(*Decoder).recvType
    1.72GB  0.84% 66.81%     1.72GB  0.84%                | reflect.unsafe_NewArray
----------------------------------------------------------+-------------
    1.69GB  0.83% 67.64%     1.69GB  0.83%                | text/template.addFuncs
----------------------------------------------------------+-------------
                                            0.66GB 39.80% |   os.statNolog
                                            0.23GB 13.56% |   os.ReadFile
    1.66GB  0.81% 68.45%     1.66GB  0.81%                | syscall.ByteSliceFromString
----------------------------------------------------------+-------------
                                            0.89GB 53.63% |   github.com/go-openapi/spec.expandSchemaRef
                                            0.48GB 29.18% |   github.com/go-openapi/analysis.(*Spec).analyzeSchema
    1.65GB  0.81% 69.26%     1.65GB  0.81%                | net/url.parse
----------------------------------------------------------+-------------
                                            1.59GB 99.94% |   go/parser.(*parser).parseFile (inline)
    1.59GB  0.78% 70.04%     1.59GB  0.78%                | go/ast.NewScope
----------------------------------------------------------+-------------
                                            1.56GB 99.94% |   go/parser.(*parser).parseFile (inline)
    1.56GB  0.76% 70.81%     1.56GB  0.76%                | go/ast.NewObj
----------------------------------------------------------+-------------
                                            1.28GB 83.01% |   github.com/go-openapi/spec.expandSchemaRef (inline)
    1.54GB  0.76% 71.56%     1.54GB  0.76%                | strings.(*Builder).WriteString
----------------------------------------------------------+-------------
                                            4.40GB 85.25% |   encoding/gob.(*Decoder).decOpFor
                                            1.21GB 23.52% |   encoding/gob.(*Decoder).recvType
    1.45GB  0.71% 72.27%     5.17GB  2.53%                | encoding/gob.(*Decoder).getDecEnginePtr
                                            4.84GB 93.62% |   encoding/gob.(*Decoder).compileDec
----------------------------------------------------------+-------------
                                            4.84GB   100% |   encoding/gob.(*Decoder).getDecEnginePtr
    1.38GB  0.68% 72.95%     4.84GB  2.37%                | encoding/gob.(*Decoder).compileDec
                                            4.50GB 92.96% |   encoding/gob.(*Decoder).decOpFor
                                            1.18GB 24.42% |   encoding/gob.overflow (inline)
----------------------------------------------------------+-------------
                                            0.65GB 40.71% |   github.com/go-openapi/analysis.(*Spec).analyzeSchema
                                            0.63GB 39.89% |   golang.org/x/tools/internal/gopathwalk.(*walker).walk
    1.31GB  0.64% 73.59%     1.59GB  0.78%                | path.Join
----------------------------------------------------------+-------------
                                            0.26GB 21.46% |   os.ReadFile
    1.20GB  0.59% 74.18%     1.20GB  0.59%                | os.newFile
----------------------------------------------------------+-------------
                                           13.04GB 99.19% |   github.com/goccy/go-json/internal/decoder.(*mapDecoder).Decode
                                            6.05GB 46.06% |   github.com/goccy/go-json/internal/decoder.(*sliceDecoder).Decode
                                            0.33GB  2.47% |   github.com/goccy/go-json.unmarshal
    1.12GB  0.55% 74.73%    13.14GB  6.44%                | github.com/goccy/go-json/internal/decoder.(*interfaceDecoder).decodeEmptyInterface
                                           12.66GB 96.37% |   github.com/goccy/go-json/internal/decoder.(*mapDecoder).Decode
                                            6.33GB 48.15% |   github.com/goccy/go-json/internal/decoder.(*sliceDecoder).Decode
----------------------------------------------------------+-------------
                                           11.62GB 70.24% |   github.com/goccy/go-json.unmarshal
                                            6.33GB 38.26% |   github.com/goccy/go-json/internal/decoder.(*interfaceDecoder).decodeEmptyInterface
    1.11GB  0.54% 75.27%    16.54GB  8.10%                | github.com/goccy/go-json/internal/decoder.(*sliceDecoder).Decode
                                            9.91GB 59.94% |   github.com/go-openapi/spec.(*Schema).UnmarshalJSON
                                            6.05GB 36.60% |   github.com/goccy/go-json/internal/decoder.(*interfaceDecoder).decodeEmptyInterface
                                            0.86GB  5.23% |   reflect.unsafe_NewArray
                                            0.78GB  4.73% |   github.com/goccy/go-json.unmarshal
----------------------------------------------------------+-------------
    1.09GB  0.53% 75.80%     4.48GB  2.19%                | encoding/gob.(*Decoder).recvType
                                            1.21GB 27.13% |   encoding/gob.(*Decoder).getDecEnginePtr
                                            0.75GB 16.82% |   reflect.unsafe_NewArray
----------------------------------------------------------+-------------
                                            2.12GB   100% |   go/parser.(*parser).parseBinaryExpr
    1.07GB  0.53% 76.33%     2.12GB  1.04%                | go/parser.(*parser).parseCallOrConversion
                                            0.90GB 42.58% |   go/parser.(*parser).parseBinaryExpr
----------------------------------------------------------+-------------
                                            1.18GB 85.97% |   encoding/gob.(*Decoder).compileDec (inline)
    0.99GB  0.48% 76.81%     1.37GB  0.67%                | encoding/gob.overflow
----------------------------------------------------------+-------------
                                           18.61GB   100% |   golang.org/x/tools/internal/gopathwalk.(*walker).walk
    0.99GB  0.48% 77.30%    18.61GB  9.11%                | golang.org/x/tools/internal/imports.(*dirInfoCache).Store
                                           16.10GB 86.50% |   golang.org/x/tools/internal/imports.(*ModuleResolver).canonicalize
                                            0.98GB  5.28% |   go/parser.ParseFile
                                            0.30GB  1.63% |   os.statNolog
----------------------------------------------------------+-------------
                                            4.50GB   100% |   encoding/gob.(*Decoder).compileDec
    0.73GB  0.36% 77.65%     4.50GB  2.20%                | encoding/gob.(*Decoder).decOpFor
                                            4.40GB 97.96% |   encoding/gob.(*Decoder).getDecEnginePtr
----------------------------------------------------------+-------------
                                            3.08GB 59.32% |   go/parser.(*parser).parseIfStmt
                                            2.11GB 40.69% |   go/parser.(*parser).parseStmtList
    0.67GB  0.33% 77.98%     5.19GB  2.54%                | go/parser.(*parser).parseSimpleStmt
                                            4.43GB 85.43% |   go/parser.(*parser).parseExprList
----------------------------------------------------------+-------------
                                            4.43GB 76.71% |   go/parser.(*parser).parseSimpleStmt
                                            1.31GB 22.61% |   go/parser.(*parser).parseReturnStmt
    0.67GB  0.33% 78.31%     5.77GB  2.83%                | go/parser.(*parser).parseExprList
                                            5.09GB 88.09% |   go/parser.(*parser).parseBinaryExpr
----------------------------------------------------------+-------------
                                            9.32GB 98.49% |   go/parser.(*parser).parseFuncDecl
                                            3.98GB 42.12% |   go/parser.(*parser).parseBlockStmt
    0.65GB  0.32% 78.63%     9.46GB  4.63%                | go/parser.(*parser).parseStmtList
                                            6.84GB 72.26% |   go/parser.(*parser).parseIfStmt
                                            2.11GB 22.30% |   go/parser.(*parser).parseSimpleStmt
                                            1.61GB 17.01% |   go/parser.(*parser).parseReturnStmt
                                            0.74GB  7.84% |   go/parser.(*parser).parseBlockStmt
                                            0.49GB  5.21% |   go/parser.(*parser).parseGenDecl
----------------------------------------------------------+-------------
                                           13.51GB 96.68% |   github.com/goccy/go-json.MarshalWithOption
    0.59GB  0.29% 78.92%    13.98GB  6.85%                | github.com/goccy/go-json/internal/encoder.AppendMarshalJSON
                                           10.12GB 72.37% |   github.com/goccy/go-json.MarshalWithOption (inline)
                                            8.63GB 61.73% |   github.com/go-openapi/spec.Schema.MarshalJSON
                                            2.19GB 15.66% |   github.com/goccy/go-json/internal/encoder.compactObject
                                            2.14GB 15.34% |   github.com/go-openapi/spec.SchemaProperties.ToOrderedSchemaItems
                                            1.43GB 10.21% |   bytes.(*Buffer).grow
----------------------------------------------------------+-------------
                                           10.66GB 91.38% |   github.com/goccy/go-json.unmarshal
    0.51GB  0.25% 79.17%    11.67GB  5.72%                | github.com/goccy/go-json/internal/decoder.(*ptrDecoder).Decode
                                            7.68GB 65.82% |   github.com/goccy/go-json.unmarshal
                                            3.65GB 31.25% |   github.com/go-openapi/spec.(*SchemaOrArray).UnmarshalJSON
                                            0.76GB  6.51% |   github.com/go-openapi/spec.(*Schema).UnmarshalJSON
----------------------------------------------------------+-------------
                                            6.84GB   100% |   go/parser.(*parser).parseStmtList
    0.46GB  0.23% 79.40%     6.84GB  3.35%                | go/parser.(*parser).parseIfStmt
                                            4.23GB 61.93% |   go/parser.(*parser).parseBlockStmt
                                            3.08GB 45.00% |   go/parser.(*parser).parseSimpleStmt
----------------------------------------------------------+-------------
                                            2.04GB   100% |   go/parser.(*parser).parseParameters
    0.46GB  0.22% 79.62%     2.04GB  1.00%                | go/parser.(*parser).parseParameterList
                                            0.67GB 32.71% |   go/parser.(*parser).parseIdent
----------------------------------------------------------+-------------
                                            6.88GB 21.67% |   golang.org/x/tools/internal/imports.(*ModuleResolver).canonicalize
                                            0.98GB  3.09% |   golang.org/x/tools/internal/imports.(*dirInfoCache).Store
    0.43GB  0.21% 79.83%    31.77GB 15.56%                | go/parser.ParseFile
                                           21.57GB 67.91% |   go/parser.(*parser).parseFile
                                            9.23GB 29.04% |   os.ReadFile
----------------------------------------------------------+-------------
                                            7.53GB 90.01% |   github.com/go-openapi/validate.NewSchemaValidator
    0.43GB  0.21% 80.04%     8.36GB  4.10%                | github.com/go-openapi/validate.newSchemaPropsValidator
                                            8.05GB 96.21% |   github.com/go-openapi/validate.NewSchemaValidator
----------------------------------------------------------+-------------
                                           24.74GB 94.35% |   github.com/go-openapi/spec.expandSchema
    0.41GB   0.2% 80.24%    26.22GB 12.84%                | github.com/go-openapi/spec.expandSchemaRef
                                           12.97GB 49.47% |   github.com/go-openapi/spec.(*Schema).UnmarshalJSON
                                           11.85GB 45.19% |   github.com/go-openapi/spec.expandSchema
                                            3.81GB 14.53% |   github.com/go-openapi/spec.Schema.MarshalJSON
                                            1.28GB  4.89% |   strings.(*Builder).WriteString (inline)
                                            0.89GB  3.38% |   net/url.parse
                                            0.27GB  1.03% |   github.com/goccy/go-json.unmarshal
----------------------------------------------------------+-------------
                                           21.57GB 99.70% |   go/parser.ParseFile
    0.41GB   0.2% 80.44%    21.64GB 10.60%                | go/parser.(*parser).parseFile
                                           13.46GB 62.21% |   go/parser.(*parser).parseFuncDecl
                                            1.95GB  9.03% |   go/ast.(*Scope).Insert (inline)
                                            1.59GB  7.34% |   go/ast.NewScope (inline)
                                            1.56GB  7.21% |   go/ast.NewObj (inline)
                                            1.56GB  7.20% |   go/parser.(*parser).parseGenDecl
----------------------------------------------------------+-------------
                                               6GB 95.99% |   github.com/go-openapi/validate.(*objectValidator).Validate
    0.40GB   0.2% 80.63%     6.25GB  3.06%                | github.com/go-openapi/validate.(*objectValidator).validatePatternProperty
                                            3.55GB 56.83% |   github.com/go-openapi/validate.(*objectValidator).Validate
                                            2.42GB 38.68% |   github.com/go-openapi/validate.(*schemaPropsValidator).Validate
                                               1GB 16.04% |   regexp.compile
----------------------------------------------------------+-------------
                                            4.23GB 94.64% |   go/parser.(*parser).parseIfStmt
                                            0.74GB 16.59% |   go/parser.(*parser).parseStmtList
    0.37GB  0.18% 80.81%     4.47GB  2.19%                | go/parser.(*parser).parseBlockStmt
                                            3.98GB 89.08% |   go/parser.(*parser).parseStmtList
                                            0.34GB  7.71% |   go/token.(*File).AddLine
----------------------------------------------------------+-------------
                                            2.48GB 99.31% |   go/parser.(*parser).parseFuncDecl
    0.36GB  0.18% 80.99%     2.50GB  1.22%                | go/parser.(*parser).parseParameters
                                            2.04GB 81.79% |   go/parser.(*parser).parseParameterList
----------------------------------------------------------+-------------
                                            5.09GB 99.85% |   go/parser.(*parser).parseExprList
                                            0.90GB 17.71% |   go/parser.(*parser).parseCallOrConversion
    0.34GB  0.17% 81.16%     5.09GB  2.49%                | go/parser.(*parser).parseBinaryExpr
                                            2.25GB 44.21% |   go/parser.(*parser).parseOperand
                                            2.12GB 41.58% |   go/parser.(*parser).parseCallOrConversion
                                            0.46GB  9.01% |   go/parser.(*parser).parseIdent
----------------------------------------------------------+-------------
                                            2.19GB 91.17% |   github.com/goccy/go-json/internal/encoder.AppendMarshalJSON
    0.33GB  0.16% 81.32%     2.40GB  1.18%                | github.com/goccy/go-json/internal/encoder.compactObject
                                            1.95GB 81.30% |   github.com/goccy/go-json/internal/encoder.compactString
----------------------------------------------------------+-------------
                                            2.27GB 24.38% |   encoding/gob.(*Encoder).sendActualType
                                            1.43GB 15.35% |   github.com/goccy/go-json/internal/encoder.AppendMarshalJSON
                                            1.34GB 14.41% |   github.com/go-openapi/spec.Schema.MarshalJSON
    0.32GB  0.16% 81.48%     9.30GB  4.55%                | bytes.(*Buffer).grow
                                            8.97GB 96.52% |   bytes.growSlice
----------------------------------------------------------+-------------
                                           13.46GB 99.22% |   go/parser.(*parser).parseFile
    0.32GB  0.16% 81.63%    13.57GB  6.65%                | go/parser.(*parser).parseFuncDecl
                                            9.32GB 68.67% |   go/parser.(*parser).parseStmtList
                                            2.48GB 18.30% |   go/parser.(*parser).parseParameters
                                            0.36GB  2.68% |   go/token.(*File).AddLine
                                            0.22GB  1.64% |   go/parser.(*parser).parseIdent
----------------------------------------------------------+-------------
                                            1.56GB 75.96% |   go/parser.(*parser).parseFile
                                            0.49GB 24.04% |   go/parser.(*parser).parseStmtList
    0.28GB  0.14% 81.77%     2.05GB  1.00%                | go/parser.(*parser).parseGenDecl
                                            0.35GB 17.07% |   go/token.(*File).AddLine
                                            0.34GB 16.67% |   go/parser.(*parser).parseIdent
----------------------------------------------------------+-------------
                                            1.61GB   100% |   go/parser.(*parser).parseStmtList
    0.28GB  0.14% 81.91%     1.61GB  0.79%                | go/parser.(*parser).parseReturnStmt
                                            1.31GB 81.12% |   go/parser.(*parser).parseExprList
----------------------------------------------------------+-------------
    0.26GB  0.13% 82.03%    25.40GB 12.44%                | golang.org/x/tools/internal/gopathwalk.(*walker).walk
                                           18.61GB 73.25% |   golang.org/x/tools/internal/imports.(*dirInfoCache).Store
                                            2.40GB  9.46% |   os.statNolog
                                            1.78GB  7.02% |   os.ReadFile
                                            0.63GB  2.50% |   path.Join
                                            0.60GB  2.37% |   strings.(*Builder).grow
----------------------------------------------------------+-------------
                                           10.95GB 49.22% |   github.com/go-openapi/validate.(*objectValidator).Validate
                                            8.05GB 36.16% |   github.com/go-openapi/validate.newSchemaPropsValidator
    0.24GB  0.12% 82.15%    22.25GB 10.90%                | github.com/go-openapi/validate.NewSchemaValidator
                                           20.09GB 90.29% |   github.com/go-openapi/spec.expandSchema
                                            7.53GB 33.83% |   github.com/go-openapi/validate.newSchemaPropsValidator
                                            0.72GB  3.23% |   github.com/goccy/go-json.unmarshal
----------------------------------------------------------+-------------
                                               1GB 49.87% |   github.com/go-openapi/validate.(*objectValidator).validatePatternProperty
                                            0.99GB 49.00% |   github.com/go-openapi/validate.(*objectValidator).Validate
    0.22GB  0.11% 82.26%     2.01GB  0.98%                | regexp.compile
----------------------------------------------------------+-------------
                                            2.25GB   100% |   go/parser.(*parser).parseBinaryExpr
    0.22GB  0.11% 82.37%     2.25GB  1.10%                | go/parser.(*parser).parseOperand
                                            1.98GB 88.04% |   go/parser.(*parser).parseIdent
----------------------------------------------------------+-------------
                                           16.10GB   100% |   golang.org/x/tools/internal/imports.(*dirInfoCache).Store
    0.22GB  0.11% 82.47%    16.10GB  7.88%                | golang.org/x/tools/internal/imports.(*ModuleResolver).canonicalize
                                            6.88GB 42.77% |   go/parser.ParseFile
                                            6.35GB 39.45% |   go/build.(*Context).matchFile
                                            1.15GB  7.16% |   os.(*File).readdir
                                            0.23GB  1.41% |   strings.(*Builder).grow
----------------------------------------------------------+-------------
    0.20GB 0.099% 82.57%     1.27GB  0.62%                | github.com/go-openapi/swag.(*splitter).breakCasualString
                                            0.88GB 69.72% |   strings.(*Builder).grow
----------------------------------------------------------+-------------
    0.20GB 0.099% 82.67%     1.40GB  0.69%                | github.com/kr/text.(*indentWriter).Write
----------------------------------------------------------+-------------
                                            7.80GB 56.93% |   github.com/go-openapi/validate.(*objectValidator).Validate
                                            2.42GB 17.65% |   github.com/go-openapi/validate.(*objectValidator).validatePatternProperty
    0.20GB 0.098% 82.77%    13.71GB  6.71%                | github.com/go-openapi/validate.(*schemaPropsValidator).Validate
                                           13.47GB 98.29% |   github.com/go-openapi/validate.(*objectValidator).Validate
----------------------------------------------------------+-------------
                                            3.65GB 76.16% |   github.com/goccy/go-json/internal/decoder.(*ptrDecoder).Decode
    0.20GB 0.097% 82.87%     4.79GB  2.34%                | github.com/go-openapi/spec.(*SchemaOrArray).UnmarshalJSON
                                            4.60GB 96.06% |   github.com/goccy/go-json.unmarshal
----------------------------------------------------------+-------------
                                            6.35GB 85.09% |   golang.org/x/tools/internal/imports.(*ModuleResolver).canonicalize
    0.19GB 0.093% 82.96%     7.46GB  3.65%                | go/build.(*Context).matchFile
                                            4.51GB 60.41% |   bufio.NewReaderSize (inline)
                                            2.31GB 30.95% |   go/build.(*importReader).readByte
----------------------------------------------------------+-------------

## With spec using jsoniterator

