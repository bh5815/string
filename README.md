# String Converter

It is a command line tool to convert clipboard strings and output.

## Usage

Copy texts in an editor, and execute it in the cmd window.

> string option

As below, outputs can be copied to the clipboard again.
So, you can paste the result to the editor.

> string option | clip

## Options

* -ym: Converts YAML to markdown format with headers.
It support YAML blocks.
* -yml: Converts YAML to markdown list format without headers.
It don't support YAML blocks.

## Example

Copied YAML:

```yaml
Header: Bla Bla
List1:
  Foo1: bar
  Foo2: bar
List2:
  - Item1
  - Item2
Block: |
  Line1
  Line2
```

YAML to markdown:

~~~cmd
C:\>string -ym

## Header

Bla Bla

## List1

* Foo1: bar
* Foo2: bar

## List2

* Item1
* Item2

## Block

```txt
Line1
Line2
```
~~~

YAML to markdown list:

~~~cmd
C:\>string -yml
* Header: Bla Bla
* List1:
  * Foo1: bar
  * Foo2: bar
* List2:
  * Item1
  * Item2
* Block: |
  * Line1
  * Line2
~~~
