var map = document.getElementById("map");

function mapClick(event) {
    var form = document.getElementById("mapform");
    var button = document.createElement("input");
    var clickx = document.createElement("input");
    var clicky = document.createElement("input");
    var shift = document.createElement("input");
    button.setAttribute("type", "hidden");
    shift.setAttribute("type", "hidden");
    clickx.setAttribute("type", "hidden");
    clicky.setAttribute("type", "hidden");
    button.setAttribute("name", "button");
    clickx.setAttribute("name", "clickx");
    clicky.setAttribute("name", "clicky");
    shift.setAttribute("name", "shift");
    button.setAttribute("value", event.button);
    clickx.setAttribute("value", event.pageX - map.x);
    clicky.setAttribute("value", event.pageY - map.y); 
    if (event.shiftKey) {
        shift.setAttribute("value", "true"); 
    } else {
        shift.setAttribute("value", "false"); 
    }
    form.appendChild(button);
    form.appendChild(shift);
    form.appendChild(clickx);
    form.appendChild(clicky);
    form.submit();
    return false;
}

function mapWheel(event) {
    event.preventDefault();
    var up;
    if (event.detail) {
        up = -1*event.detail/3;
    } else {
        up = (event.wheelDelta)/120;
    }
    if (up > 0 && up < 1) {
        up = 1;
    } else if (up < 0 && up > -1) {
        up = -1;
    } else if (up === 0) {
        return true;
    }
    var form = document.getElementById("mapform");
    var button = document.createElement("input");
    var clickx = document.createElement("input");
    var clicky = document.createElement("input");
    var shift = document.createElement("input");
    button.setAttribute("type", "hidden");
    shift.setAttribute("type", "hidden");
    clickx.setAttribute("type", "hidden");
    clicky.setAttribute("type", "hidden");
    button.setAttribute("name", "button");
    clickx.setAttribute("name", "clickx");
    clicky.setAttribute("name", "clicky");
    shift.setAttribute("name", "shift");
    button.setAttribute("value", "3");
    clickx.setAttribute("value", 0);
    clicky.setAttribute("value", up); 
    if (event.shiftKey) {
        shift.setAttribute("value", "true"); 
    } else {
        shift.setAttribute("value", "false"); 
    }
    form.appendChild(button);
    form.appendChild(shift);
    form.appendChild(clickx);
    form.appendChild(clicky);

    form.submit();
    return false;
}
map.onmousedown = mapClick;
map.oncontextmenu= function() { return false; };
map.onmousewheel = mapWheel;
map.onDOMMouseScroll = mapWheel;
map.addEventListener("DOMMouseScroll", mapWheel);
