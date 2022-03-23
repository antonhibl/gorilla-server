function executeScript(input_code) {
    var results = new Function (input_code);
    return(results());
}

document.body.addEventListener("click", function() {
    input = document.getElementById("input").value;
    result = executeScript(input);
    document.getElementById("output").innerText = result; 
});