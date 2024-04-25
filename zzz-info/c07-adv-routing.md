# CH07 Advanced routing 

## Target endpoints

Before | After  | Handler 
---|---|---
`ANY`  `/` | `GET`  `/` | `HandleHome` | 
`ANY`  `/snippet/view?id=123` | `GET`  `/snippet/view/:id`  | `HandleViewSnippet` | 
=none= | `GET` `/snippet/create` | `HandleCreateSnippetForm` | 
`POST` | `/snippet/create` | `HandleCreateSnippet` 
`ANY`  `/static/` | `GET` `/static/` | `http.Fileserver(...)` | Using the `http.Fileserver()` handler + `http.StripPrefix()`. 