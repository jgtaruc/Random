var name_entry = document.querySelector("#entry");
var msg = document.querySelector("#msg_id");

entry.addEventListener('input', function(){
	var xhr = new XMLHttpRequest();
	xhr.open('POST', '/api/check');
	xhr.send(entry.value);
	xhr.addEventListener('readystatechange', function(){
		if (xhr.readyState === 4 && xhr.status === 200){
			var taken = xhr.responseText;
			if (taken == 'true'){
				msg.style.display = 'block';
				document.querySelector("#btn").disabled = true;
				document.querySelector("#btn").style.background = 'red';
			}
			else{
				msg.style.display = 'none';
				document.querySelector("#btn").disabled = false;
				document.querySelector("#btn").style.background = 'rgb(244, 194, 13)';
			}
		}
	});
});