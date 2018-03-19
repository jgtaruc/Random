var quotes = [
	"Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
	"Nam tristique condimentum est at euismod.",
	"Suspendisse feugiat lacus quis arcu dignissim, nec molestie odio efficitur.",
	"Vestibulum euismod, libero a pharetra sollicitudin, dui ante luctus magna, nec dictum metus arcu eget arcu."
	]

var counter = 1

function transition() {
  if(counter == 1) {
    document.querySelector("#bg1").className = "bg-active";
    document.querySelector("#bg2").className = "";
    document.querySelector("#bg3").className = "";
    document.querySelector("#bg4").className = "";
    counter++;
  }
  else if(counter == 2) {
    document.querySelector("#bg1").className = "";
    document.querySelector("#bg2").className = "bg-active";
    document.querySelector("#bg3").className = "";
    document.querySelector("#bg4").className = "";
    counter++;
  }
  else if(counter == 3) {
    document.querySelector("#bg1").className = "";
    document.querySelector("#bg2").className = "";
    document.querySelector("#bg3").className = "bg-active";
    document.querySelector("#bg4").className = "";
    counter++;
  }
  else {
    document.querySelector("#bg1").className = "";
    document.querySelector("#bg2").className = "";
    document.querySelector("#bg3").className = "";
    document.querySelector("#bg4").className = "bg-active";
    counter = 1;
  }
  var j = Math.floor((Math.random() * 4));
  document.querySelector("#quote").innerHTML = quotes[j];
}

setInterval(transition, 4000);