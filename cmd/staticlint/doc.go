/*
Package main provides complex multichecker, which consists of:
  - shift.Analyzer:						checks for shifts that exceed the width of an integer.
  - printf.Analyzer:					checks consistency of Printf format strings and arguments.
  - shadow.Analyzer:					checks for possible unintended shadowing of variables.
  - asmdecl.Analyzer:				    defines an Analyzer that reports mismatches between assembly files and Go declarations.
  - assign.Analyzer:   					check for useless assignments.
  - atomic.Analyzer:					check for common mistakes using the sync/atomic package.
  - atomicalign.Analyzer:				defines an Analyzer that checks for non-64-bit-aligned arguments to sync/atomic functions.
  - bools.Analyzer:						defines an Analyzer that detects common mistakes involving boolean operators.
  - buildssa.Analyzer:					defines an Analyzer that constructs the SSA representation of an error-free package and returns the set of all functions within it.
  - buildtag.Analyzer:					defines an Analyzer that checks build tags.
  - cgocall.Analyzer:					defines an Analyzer that detects some violations of the cgo pointer passing rules.
  - composite.Analyzer:					defines an Analyzer that checks for unkeyed composite literals.
  - copylock.Analyzer:					defines an Analyzer that checks for locks erroneously passed by value.
  - ctrlflow.Analyzer:					provides a syntactic control-flow graph (CFG) for the body of a function.
  - deepequalerrors.Analyzer:			defines an Analyzer that checks for the use of reflect.DeepEqual with error values.
  - errorsas.Analyzer:					defines an Analyzer that checks that the second argument to errors.As is a pointer to a type implementing error.
  - fieldalignment.Analyzer:			defines an Analyzer that detects structs that would use less memory if their fields were sorted.
  - findcall.Analyzer:					defines an Analyzer that serves as a trivial example and test of the Analysis API.
  - framepointer.Analyzer:				defines an Analyzer that reports assembly code that clobbers the frame pointer before saving it.
  - httpresponse.Analyzer:				defines an Analyzer that checks for mistakes using HTTP responses.
  - ifaceassert.Analyzer:				defines an Analyzer that flags impossible interface-interface type assertions.
  - inspect.Analyzer:					inspect defines an Analyzer that provides an AST inspector (golang.org/x/tools/go/ast/inspector.Inspector) for the syntax trees of a package.
  - loopclosure.Analyzer:				defines an Analyzer that checks for references to enclosing loop variables from within nested functions.
  - lostcancel.Analyzer:				defines an Analyzer that checks for failure to call a context cancellation function.
  - nilfunc.Analyzer:					defines an Analyzer that checks for useless comparisons against nil.
  - nilness.Analyzer:				    inspects the control-flow graph of an SSA function and reports errors such as nil pointer dereferences and degenerate nil pointer comparisons.
  - pkgfact.Analyzer:					demonstration and test of the package fact mechanism.
  - reflectvaluecompare.Analyzer:		defines an Analyzer that checks for accidentally using == or reflect.DeepEqual to compare reflect.Value values.
  - sigchanyzer.Analyzer:				defines an Analyzer that detects misuse of unbuffered signal as argument to signal.Notify.
  - sortslice.Analyzer:					defines an Analyzer that checks for calls to sort.Slice that do not use a slice type as first argument.
  - stdmethods.Analyzer:				defines an Analyzer that checks for misspellings in the signatures of methods similar to well-known interfaces.
  - stringintconv.Analyzer:				defines an Analyzer that flags type conversions from integers to strings.
  - tests.Analyzer:						defines an Analyzer that checks for common mistaken usages of tests and examples.
  - unmarshal.Analyzer:					defines an Analyzer that checks for passing non-pointer or non-interface types to unmarshal and decode functions.
  - unreachable.Analyzer:				defines an Analyzer that checks for unreachable code.
  - unsafeptr.Analyzer:					defines an Analyzer that checks for invalid conversions of uintptr to unsafe.Pointer.
  - unusedresult.Analyzer:				defines an analyzer that checks for unused results of calls to certain pure functions.
  - unusedwrite.Analyzer:				checks for unused writes to the elements of a struct or array object.
  - structtag.Analyzer:					checks struct field tags are well-formed.
  - nakedret.NakedReturnAnalyzer: 		checks naked returns in functions greater than a specified function length (github.com/alexkohler/nakedret).
  - bidichk.NewAnalyzer:				checks dangerous unicode character sequences in Go source files (https://github.com/breml/bidichk).
  - exitcheck.Analyzer: 				checks calling os.Exit().
  - staticcheck:						all SA*, S1*, ST1* checks from static checks (staticcheck.io).

How to use:
 1. build package:
    $> go build -o staticlint
 2. exec multichecker
    $> ./staticlint ./cmd/main.go
 3. for information type
    $> ./staticlint help
*/
package main
