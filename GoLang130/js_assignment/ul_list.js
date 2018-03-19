function build_ul(){
	var div1 = document.getElementById("div1");
	var ul = document.createElement("ul");
	ul.innerHTML = "<li id=li1 onclick=display_alert('li1')>item1</li>" +
				   "<li id=li2 onclick=display_alert('li2')>item2</li>" +
				   "<li id=li3 onclick=display_alert('li3')>item3</li>" +
				   "<li id=li4 onclick=display_alert('li4')>item4</li>";
	div1.appendChild(ul);
}