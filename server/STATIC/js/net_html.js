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
        f(pt, event.button, event.shiftKey);
    };
    return g;
};

muleObj.html.setPointClick = function(elem, f) {
    var g = muleObj.html.clickWrap(f);
    elem.addEventListener("mouseup", g);
    elem.oncontextmenu= function() { return false; };
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

function Screen(canvas) {
    this.canvas = canvas;
    var screen = this;
    canvas.addEventListener("mouseup", function(event) {
        if (!screen.handleClick) {
            return;
        }
        var clickx = event.pageX - this.offsetLeft;
        var clicky = event.pageY - this.offsetTop;
        screen.handleClick(clickx, clicky, event.button, event.shiftKey);
    });
    canvas.oncontextmenu= function() { return false; };

    var g = function(event) {
        if (!screen.handleWheel) {
            return true;
        }
        event.preventDefault();
        var up = (event.detail)? -1*event.detail/3: (event.wheelDelta)/120;
        if (up > 0 && up < 1) {
            up = 1;
        } else if (up < 0 && up > -1) {
            up = -1;
        } else if (up === 0) {
            return true;
        }
        screen.handleWheel(up, event.shiftKey, event.ctrlKey);
        return false;
    };
    canvas.onmousewheel = g;
    canvas.onDOMMouseScroll = g;
    canvas.addEventListener("DOMMouseScroll", g);
}

Screen.prototype.resize = function(width, height) {
    this.canvas.width = width;
    this.canvas.height = height;
    if (this.onResize) {
        this.onResize(width, height);
    }
};

muleObj.html.Screen = Screen;

function ScreenTransform(canvas, centerX, centerY, theta, scale, squashY) {
    this.canvas = canvas;
    this.transform = new muleObj.geometry.Transform(scale, 0,0, theta, squashY, false);
    this.transform.setInAtOut(new muleObj.geometry.Point(centerX, centerY), new muleObj.geometry.Point(canvas.width/2, canvas.height/2));
    this.in2out = function(pt) {
        return this.transform.in2out(pt);
    };
    this.out2in = function(pt) {
        return this.transform.out2in(pt);
    };
}
ScreenTransform.prototype.center = function() {
    return new muleObj.geometry.Point(this.canvas.width/2, this.canvas.height/2);
};

ScreenTransform.prototype.setCenter = function(inPt) {
    this.transform.setInAtOut(inPt, this.center());
};
ScreenTransform.prototype.rotate = function(theta) {
    this.transform.rotateAround(theta, this.center());
};
ScreenTransform.prototype.setScale = function(scale) {
    this.transform.scaleAround(scale, this.center());
};
ScreenTransform.prototype.shift = function(dx, dy) {
    this.transform.shift(dx, dy);
};
ScreenTransform.prototype.resize = function(width, height) {
    var curCenter = this.transform.out2in(this.center());
    this.canvas.width = width;
    this.canvas.height = height;
    this.transform.setInAtOut(curCenter, this.center());
};
ScreenTransform.prototype.setClick = function(f) {
    muleObj.html.setPointClick(this.canvas, f);
};
ScreenTransform.prototype.setInClick = function(f) {
    var transform = this.transform;
    var g = function(pt, button, shift) {
        var inPt = transform.out2in(pt);
        return f(inPt, button, shift);
    };
    muleObj.html.setPointClick(this.canvas, g);
};
ScreenTransform.prototype.setWheel = function(f) {
    muleObj.html.setWheel(this.canvas, f);
};

muleObj.html.ScreenTransform = ScreenTransform;


})(muleObj);
