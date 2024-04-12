# Notes 

## Disable directory list

Add a blank index.html to each dir 
```sh
$ find ./ui/static -type d -exec touch {}/index.html \;
```
A more complicated (but arguably better) solution is to create a custom implementation of `http.FileSystem`, and have it return an `os.ErrNotExist` error for any directories. A full explanation and sample code can be found here: [How to Disable FileServer Directory Listings](https://www.alexedwards.net/blog/disable-http-fileserver-directory-listings).