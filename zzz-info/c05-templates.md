# C04 Templates

Plan: 
1. Pass dynamic data on templates
2. Use function on templates
3. Create template cache 
4. Gracefully handle template rendering errors 
5. Implement a pattern for passing dynamic data to web pages 
6. Create custom functions to displ,ay data in html template

## Content escaping 

The `html/template` package dynamically escapes all {{}} tags. 
Strips all HTML comments. 

## Calling methods from tags

`time.Time.Weekday()`

```html
<time>Created: {{.Created.Weekday}}</time><br>
<!-- Add 6 months -->
<time>Expires: {{.Expires.AddDate 0 6 0}}</time>
```

## If, else, end, with, range 

Action | Descr 
---|---
`{{if .Foo}} C1 {{else}} C2 {{end}}` | If `.foo` not empty, render C1, else render c2. 
`{{with .Foo}} C1 {{else}} C2 {{end}}` | If `.foo` not empty, set dot to `.Foo`, and render C1, otherwise render C2. 
`{{range .Foo}} C1 {{else}} C2 {{end}}` | If len(.Foo) > 0, loop over each element, setting dot to the value of each element and rendering the content C1. If en(.Foo) == 0, render C2. Foo: array, slice, map, or channel only.

## Eq, ne, not, or, index, printf, len, $bar 
Action | Descr 
---|---
`{{eq .Foo .Bar}}` | Yields true if `.Foo == .Bar`
`{{ne .Foo .Bar}}` | Yields true if `.Foo != .Bar`
`{{not .Foo }}` | Yields the boolean negation of `.Foo`
`{{or .Foo .Bar}}` | Yields .Foo if .Foo not empty; otherwise yields .Bar
`{{index .Foo i}}` | Yields the value of `.Foo[i]` (Map, slice or array only, i int)
`{{printf "%s-%s" .Foo .Bar}}` | Yields form string containing .Foo and .Bar (sim to `Sprintf()`)
`{{len .Foo }}` | Yields len(.Foo)
`{{$bar := len .Foo }}` | Assign `len(.Foo)` to var `$bar`


## Templates 

Add error to a template:
```html
{{define "title"}}Snippet #{{.Snippet.ID}}{{end}}
{{define "main"}}
{{ with .Snippet }}
<div class="snippet">
    <div class="metadata">
        <strong>{{.Title}}</strong>
        <span>#{{.ID}}</span>
    </div>
    {{ len nil }} <!-- Deliberate error (nil doesn't have length) -->
    <!-- ... -->
{{end}}
{{end}}
```
Query:

```
curl -i "http://localhost:1111/snippet/view?id=2"
HTTP/1.1 200 OK
Date: Sun, 21 Apr 2024 09:27:16 GMT
Content-Length: 725
Content-Type: text/html; charset=utf-8


<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link rel="stylesheet" href="./static/css/main.css">
    <link rel="shortcut icon" href="./static/img/favicon.ico" type="image/x-icon">
    <link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Ubuntu+Mono:400,700">
    <title>Snippet #2 - Snippetbox</title>
</head>
<body>
    <header>
        <h1><a href="/">S-BOX</a></h1>
    </header>
    
<nav>
    <a href="/">Home</a>
</nav>

    <main>
        

<div class="snippet">
    <div class="metadata">
        <strong>Over the wintry forest</strong>
        <span>#2</span>
    </div>
    Internal Server Error
```
Sol: 
1. Make a trial render by writing templ into a buff
2. If fails: err; if works: write to response writer
