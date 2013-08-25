$(document).ready(function() {
	// Main entry for the page:
	var id = document.location.href.split(/\//)[4]
	// var id = $(".entry", "#content")[0].id.split(/_/)[1]

	var conn;
	var msg = $("#msg");
	var log = $("#log");

	function appendLog(msg) {
		console.log(msg)
	}

	if (window["ReconnectingWebSocket"]) {
		conn = new ReconnectingWebSocket("ws://" + document.location.host
				+ "/ws/thread/" + id);
		conn.callbacks = {};
		conn.onclose = function(evt) {
			appendLog("Connection closed.");
		};
		conn.onmessage = function(evt) {
			var json = JSON.parse(evt.data)
			conn.dispatch(json.Event, json.Data)
		};
		conn.bind = function(event, callback) {
			conn.callbacks[event] = conn.callbacks[event] || [];
			conn.callbacks[event].push(callback);
			return conn.bind;
		};
		conn.dispatch = function(event, data) {
			var chain = conn.callbacks[event];
			if (typeof chain == "undefined") {
				return;
			}
			for ( var i = 0; i < chain.length; i++) {
				chain[i](data);
			}
		};
		
		conn.bind('num_connections', function(data){
			console.log(data);
		});
		
		conn.bind('thread_post', function(data){
			app.appendTemplate("EntryTemplate", data, "children_"+data.ParentId);
			console.log(data);
		});
	} else {
		appendLog("Your browser does not support WebSockets.")
	}
});