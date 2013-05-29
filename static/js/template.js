var app = {
	_tmplCache : {},
	
	parseTemplate : function(str, data) {
		// From http://www.west-wind.com/weblog/posts/2008/Oct/13/Client-Templating-with-jQuery
		// on 2013-04-23
		// <summary>
		// Client side template parser that uses &lt;#= #&gt; and &lt;# code #&gt;
		// expressions and # # code blocks for template expansion.
		// NOTE: chokes on single quotes in the document in some situations
		// use &amp;rsquo; for literals in text and avoid any single quote
		// attribute delimiters.
		// </summary>
		// <param name="str" type="string">The text of the template to
		// expand</param>
		// <param name="data" type="var">
		// Any data that is to be merged. Pass an object and
		// that object's properties are visible as variables.
		// </param>
		// <returns type="string" />
		var err = "";
		try {
			var func = this._tmplCache[str];
			if (!func) {
	            var strFunc =
	            "var p=[],print=function(){p.push.apply(p,arguments);};" +
	                        "with(obj){p.push('" +
	            str.replace(/[\r\t\n]/g, " ")
	               .replace(/'(?=[^#]*#>)/g, "\t")
	               .split("'").join("\\'")
	               .split("\t").join("'")
	               .replace(/<#=(.+?)#>/g, "',$1,'")
	               .split("<#").join("');")
	               .split("#>").join("p.push('")
	               + "');}return p.join('');";
	            // alert(strFunc);
				func = new Function("obj", strFunc);
				this._tmplCache[str] = func;
			}
			return func(data);
		} catch (e) {
			err = e.message;
		}
		return "< # ERROR: " + err.htmlEncode() + " # >";
	},
	
	// templateId is the ID of a template;
	// data is an object containing named variables;
	// containerId is the ID of the element to which to append the new element. 
	appendTemplate : function(templateId, data, containerId) {
		var newDiv = document.createElement("div");
		
		if(typeof data.id != "undefined"){
			newDiv.id = data.id
		}
		
		newDiv.innerHTML = this.parseTemplate(document.getElementById(templateId).innerHTML, data);
		newDiv = newDiv.children[0];
		
		document.getElementById(containerId).appendChild(newDiv);
		
		return;
	},
	
	// templateId is the ID of a template;
	// data is an object containing named variables;
	// containerId is the ID of the element to which to prepend the new element. 
	appendTemplate : function(templateId, data, containerId) {
		var newDiv = document.createElement("div");
		
		if(typeof data.id != "undefined"){
			newDiv.id = data.id
		}
		
		newDiv.innerHTML = this.parseTemplate(document.getElementById(templateId).innerHTML, data);
		newDiv = newDiv.children[0];
		
		//Prepend
		var par = document.getElementById(containerId) 
		par.insertBefore(newDiv, par.firstChild);
		
		return;
	},
	
	// templateId is the ID of a template;
	// data is an object containing named variables;
	// containerId is the ID of the element to which to append the new element. 
	replaceTemplate : function(templateId, data, containerId) {
		var newDiv = document.createElement("div");
		
		if(typeof data.id != "undefined"){
			newDiv.id = data.id
		}
		
		var replacement = this.parseTemplate(document.getElementById(templateId).innerHTML, data);
		
		newDiv.innerHTML = replacement;
		newDiv = newDiv.children[0];
		
		cont = document.getElementById(containerId);
		
		cont.innerHTML = newDiv.innerHTML;
		
		return;
	},
}