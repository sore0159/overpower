(function() {
    var canvas = document.getElementById('mainscreen');
    function mapClick(event) {
        var clickx = event.pageX - canvas.offsetLeft;
        var clicky = event.pageY - canvas.offsetTop;
        this.muleClicked([clickx, clicky], event.button, event.shiftKey);
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
        this.muleWheeled(up, event.shiftKey);
        return false;
    }
    canvas.muleWheeled = function(up, shiftKey) {
        console.log("wheel", up, shiftKey);
    };
    canvas.muleClicked = function(pt, button, shiftKey) {
        console.log("click", pt[0], pt[1], button, shiftKey);
    };
    canvas.onmousedown = mapClick;
    canvas.oncontextmenu= function() { return false; };
    canvas.onmousewheel = mapWheel;
    canvas.onDOMMouseScroll = mapWheel;
    canvas.addEventListener("DOMMouseScroll", mapWheel);
})();
