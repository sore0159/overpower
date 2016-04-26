(function(muleObj) {

if (!document) {
    console.log("HTML LIB FAILED: NO DOCUMENT ENVIRONMENT FOUND");
    return;
}
if (!muleObj.html) {
    muleObj.html = {};
}
muleObj.html.clear = function(elem) {
    while (elem.firstChild) {
        elem.removeChild(elem.firstChild);
    }
};
muleObj.html.spur = function(elem, kind, text) {
    var child = document.createElement(kind);
    if (text) {
        child.textContent = text;
    }
    this.elem.appendChild(child);
    return child;
};

muleObj.html.clickWrap = function(f) {
    var g = function(event) {
        var clickx = event.pageX - this.offsetLeft;
        var clicky = event.pageY - this.offsetTop;
        var pt  = (muleObj.geometry)? new muleObj.geometry.Point(clickx, clicky) : {x: clickx, y: clicky};
        f(pt, event.button, event.ShiftKey);
    };
    return g;
};

muleObj.html.wheelWrap = function(f) {
    var g = function(event) {
        event.preventDefault();
        var up = (event.detail)? -1*event.detail/3: (event.wheelDelta)/120;
        if (up > 0 && up < 1) {
            up = 1;
        } else if (up < 0 && up > -1) {
            up = -1;
        } else if (up === 0) {
            return true;
        }
        f(up, event.shiftKey, event.ctrlKey);
        return false;
    };
    return g;
};

muleObj.html.setWheel = function(elem, f) {
    var g = muleObj.html.wheelWrap(f);
    elem.onmousewheel = g;
    elem.onDOMMouseScroll = g;
    elem.addEventListener("DOMMouseScroll", g);
};

muleObj.html.setPointClick = function(elem, f) {
    var g = muleObj.html.clickWrap(f);
    elem.addEventListener("mouseup", g);
    elem.oncontextmenu= function() { return false; };
};


})(muleObj);
