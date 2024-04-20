# C04 Templates

Plan: 
1. Pass dynamic data on templates
1. Use function on templates
1. Create template cache 
1. Gracefully handle template rendering errors 
1. Implement a pattern for passing dynamic data to web pages 
1. Create custom functions to displ,ay data in html template

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



