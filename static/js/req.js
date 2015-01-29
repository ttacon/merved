/**
 *
 * req.js - Lightly wrapping native js functionality.
 *
 * Copyright (c) 2014 Trey Tacon (ttacon@gmail.com)
 * Licensed under the MIT License (see LICENSE file).
 *
 * Created 17/01/2014
 *
 */

(function(window, document) {

	var $req = {
		get: function(url, options) {
			return new Promise(function(resolve, reject) {
				var req = new XMLHttpRequest();
				req.open('GET', url);

				req.onload = function() {
					// This is called even on 404 etc
					// so check the status
					if (req.status >=  200 && req.status < 300) {
						resolve(req);
					} else {
						reject(req);
					}
				};
				// Handle network errors
				req.onerror = function() {
					reject(Error("Network Error"));
				};

				req.send();
			});
		},
		post: function(url, data, options) {
			return new Promise(function(resolve, reject) {
				var req = new XMLHttpRequest();
				req.open('POST', url, true);

				// for now this is just JSON, just cuz...
				req.setRequestHeader('Content-Type', 'application/json');

				req.onload = function() {
					// This is called even on 404 etc
					// so check the status
					if (req.status >=  200 && req.status < 300) {
						resolve(req);
					} else {
						reject(req);
					}
				};
				// Handle network errors
				req.onerror = function() {
					reject(Error("Network Error"));
				};

				req.send(data);
			});
		},
		put: function(url, data, options) {

		},
		delete: function(url, options) {

		}
	};
	

	// cause I'm lazy...
	window.$req = $req;

})(window, document);
