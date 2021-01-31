// This function will update the LEDs on the image to reflect the selected color.
function changeColor(AreaID) {
    var obj = document.getElementById("parent_svg").contentDocument.getElementById(AreaID);
    var newColor = document.getElementById("colorpicker").value;
    var newBright = document.getElementById("bright").value;
    obj.style.fill = newColor;
    obj.style.stroke = newColor;
    document.getElementById("input_" + AreaID).value = newColor;
    document.getElementById("input_" + AreaID + "b").value = newBright;
}

// You have to convert the form, otherwise, the data is sent as multipart/form-data
// and it will not be parsable by golang.
function urlencodeFormData(fd) {
    var params = new URLSearchParams();
    for(var pair of fd.entries()) {
        typeof pair[1] == 'string' && params.append(pair[0], pair[1]);
    }
    return params.toString();
}

// Async call to light the LEDs on Blinkt.
function submitColors() {
    var f = document.getElementById("form");
    var formData = urlencodeFormData(new FormData(f));
    const xhr = new XMLHttpRequest();
    xhr.open("POST", "/update_led");
    xhr.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
    xhr.send(formData);
}