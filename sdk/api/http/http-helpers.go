/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package sdkhttp

import (
	"html/template"
)

// HttpHelpers are methods available in html templates as .Helpers.
// For example, to use the Translate() method in html templates, use <% .Helpers.Translate "label" "network_settings" %>.
type HttpHelpers interface {

	// Returns the uri path of a static file in resources/assets directory from your plugin
	AssetPath(path string) (uri string)

	// Returns the uri path of a file in resources/assets directory from your plugin.
	// The file is parsed using text/template go module with access to <% .Helpers %> object.
	AssetWithHelpersPath(path string) (uri string)

	// Returns the html for the ads view.
	AdsView() (html template.HTML)

	// Translates a message into the current language settings from application config.
	// msgtype is the message type, e.g. "error", "success", "info", "warning".
	// For example, if the current language is "en", then the following code in your template:
	//  <% .Helpers.Translate "error" "some-key" %>
	// will look for the file "/resources/translations/en/error/some-key.txt" under the plugin root directory
	// and displays the text inside that file.
	Translate(msgtype string, msgk string) string

	// Returns the url for the route.
	UrlForRoute(name string, pairs ...string) (uri string)

	// Returns the url from other plugins.
	UrlForPkgRoute(pkg string, name string, pairs ...string) (uri string)

	// Returns the vue route name for a named route which can be used in vue router, e.g.
	//   $this.push({name: '<% .Helpers.VueRouteName "login" %>'})
	VueRouteName(name string) string

	// Returns the vue route path for a named route
	VueRoutePath(name string, pairs ...string) string

	// Returns the vue route path for a named route from a third-party plugin
	VuePkgRoutePath(pkg string, name string, pairs ...string) string

	// Returns the uri path of a file in resources/components directory from your plugin.
	// The file is parsed using text/template go module with access to <% .Helpers %> object.
	VueComponentPath(path string) (uri string)
}
