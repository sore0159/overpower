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

map.onmousedown = mapClick;
map.oncontextmenu= function() { return false; };
