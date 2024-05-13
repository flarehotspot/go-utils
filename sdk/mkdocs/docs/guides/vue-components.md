# Vue Components

Flare Hotspot uses [Vue.js](https://v2.vuejs.org) to build the user interface. But we are not using the standard build tools for Vue.js project since we need to support dynamic components from the plugins. Hence, the syntax for declaring vue components are slightly different. This guide will help you understand how to build and use Vue components in the Flare Hotspot project. Vue components are placed in the `resources/components` directory in your plugin.

Take a look at the following example:

```html title="resources/components/portal/Welcome.vue"
<template>
    <div>
        <h1>Welcome {{ name }}</h1>
    </div>
</template>

<script>
define(function () {
    return {
        template: template,
        data: function(){
            return {
                name: "John"
            }
        }
    }
})
</script>
```

The equivalent traditional [single-file component](https://v2.vuejs.org/v2/guide/single-file-components) would look like:

```html
<template>
    <div>
        <h1>Welcome {{ name }}</h1>
    </div>
</template>

<script>
export default {
    data() {
        return {
            name: "John"
        }
    }
}
</script>
```

## 1. The `template` variable {#template-variable}

The `template` variable is a string containing the HTML code automatically extracted from the `<template>` tag.

!!! warning "Important"
    Note that there must be **only one root html tag** of the template. The following template will not work:
    ```html
    <template>
        <h1>Welcome {{ flareView.data.name }}</h1>
        <p>Some other text</p> <!-- the <p> tag will not be displayed -->
    </template>
    ```

    Below is the correct way:
    ```html
    <template>
        <div>
            <h1>Welcome {{ flareView.data.name }}</h1>
            <p>Some other text</p>
        </div>
    </template>
    ```

## 2. Template helpers {#template-helpers}

Aside from the [HttpHelpers.VueRoutePath](../api/http-helpers.md#vueroutepath) method we used to create a link, there are other useful methods within the [HttpHelpers](../api/http-helpers.md) API. The [HttpHelpers](../api/http-helpers.md) can be accessed anywhere inside the component as `.Helpers` (notice the dot prefix) enclosed by `<%` and `%>` delimiters. Visit the [HttpHelpers](../api/http-helpers.md) API documentation to learn more.

For example, to build a link to another route, you can use the `HttpHelpers.VueRoutePath` method as shown below:

```html
<router-link :to='<% .Helpers.VueRoutePath "portal.welcome" %>'>Welcome</router-link>
```

## 3. Loading child components {#loading-child-components}

Loading child components can be done using the [HttpHelpers.VueComponentPath](../api/http-helpers.md#vuecomponentpath) method in combination with the [$flare.vueLazyLoad](../api/flare-variable.md#flare-vuelazyload) method:

The parent component:

```html title="resources/components/SampleParent.vue"
<template>
    <sample-child></sample-child>
</template>

<script>
define(function(){
    var child = $flare.vueLazyLoad('<% .Helpers.VueComponentPath "SampleChild.vue" %>');

    return {
        template: template,
        components: {
            'sample-child': child
        }
    };
});
</script>
```

The child component:

```html title="resources/components/SampleChild.vue"
<template>
    <div>
        <h1>Sample Child</h1>
    </div>
</template>

<script>
define(function(){
    return {
        template: template
    };
});
</script>
```

## 5. Browser Compatibility {#browser-compatibility}

Since we are not using standard build tools like webpack or vite, it is recommended to use basic form of javascript and css to ensure compatibility with older browsers.
For example, use `var` instead of `let` or `const` and use `function` instead of arrow functions.
Array `map`, `filter`, `reduce`, etc. should also be avoided.

It's important to note that when working with legacy Javascript codes, the context of `this` keyword may not be the same as in the standard Vue.js components. Hence, it is recommended to use `var self = this` to store the context of `this` keyword.

```html title="SampleComponent.vue"
<template>
  <!-- Rest of the template... -->
</template>
<script>
describe(function(){
    return {
        template: template,
        data: function(){
            return {
                name: "John"
            }
        },
        mounted: function(){
            // store the context of `this` keyword to `self`
            var self = this;

            setTimeout(function(){

                // use `self` instead of `this`
                self.name = "Doe";

            }, 2000);
        }
    })
</script>
```