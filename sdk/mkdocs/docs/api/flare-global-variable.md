# The `$flare` global variable
The `$flare` global variable is a global variable in the browser that contains helper functions to work with the Flare API.

## $flare.http {#flare-http}
The `$flare.http` object is used to make AJAX requests. It contains two methods, namely `get` and `post`.

### get {#flare-http-get}
The `$flare.http.get` method accepts two arguments, the first argument is the URL to send the form data to, and the second argument is the query params.
```js
var queryParams = {amount: 100};

$flare.http.get('/path/to/handler', queryParams)
    .then(function(response){
        console.log(response);
    })
    .catch(function(error){
        console.log(error);
    });
```

### post {#flare-http-post}
The `$flare.http.post` method accepts two arguments, the first argument is the URL to send the form data to, and the second argument is the form data.

```js
var formData = {amount: 100};

$flare.http.post('/path/to/handler', formData)
    .then(function(response){
        console.log(response);
    })
    .catch(function(error){
        console.log(error);
    });
```
