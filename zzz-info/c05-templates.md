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




