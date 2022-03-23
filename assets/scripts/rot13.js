// rot13 function
var rot13 = function() {
	// grab input text from form and read the value
    var the_text = document.getElementById("main_text").value;
	// replace the text via rot13, i.e. each letter value + 13
    the_text = the_text.replace(/[a-z]/gi, letter => String.fromCharCode(letter.charCodeAt(0) + (letter.toLowerCase() <= 'm' ? 13 : -13)));
	// replace the output text's inner html with the enciphered text
    document.getElementById("output_text").innerHTML=the_text;
};
