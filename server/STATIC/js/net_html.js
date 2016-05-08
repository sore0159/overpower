(function(muleObj) {

if (!document) {
    console.log("HTML LIB FAILED: NO DOCUMENT ENVIRONMENT FOUND");
    return;
}
if (!muleObj.html) {
    muleObj.html = {};
}
var html = muleObj.html;

html.clear = function(elem) {
    while (elem.firstChild) {
        elem.removeChild(elem.firstChild);
    }
};
html.spur = function(elem, kind, text) {
    var child = document.createElement(kind);
    if (text === 0 || text) {
        child.textContent = text;
    }
    elem.appendChild(child);
    return child;
};
html.spurText = function(elem, text) {
    var child = document.createTextNode(text);
    elem.appendChild(child);
};
html.setText = function(name, text) {
    var elem = document.getElementById(name);
    elem.textContent = text;
};

html.clickWrap = function(f) {
    var g = function(event) {
        var clickx = event.pageX - this.offsetLeft;
        var clicky = event.pageY - this.offsetTop;
        var pt  = (muleObj.geometry)? new muleObj.geometry.Point(clickx, clicky) : {x: clickx, y: clicky};
        f.call(this, pt, event.button, event.shiftKey, event.ctrlKey);
    };
    return g;
};

html.setPointClick = function(elem, f) {
    var g = html.clickWrap(f);
    elem.addEventListener("mouseup", g);
    elem.oncontextmenu= function() { return false; };
};

html.wheelWrap = function(f) {
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
        f.call(this, up, event.shiftKey, event.ctrlKey);
        return false;
    };
    return g;
};

html.setWheel = function(elem, f) {
    var g = html.wheelWrap(f);
    elem.onmousewheel = g;
    elem.onDOMMouseScroll = g;
    elem.addEventListener("DOMMouseScroll", g);
};

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
    html.setPointClick(this.canvas, f);
};
ScreenTransform.prototype.setInClick = function(f) {
    var transform = this.transform;
    var g = function(pt, button, shift, ctrl) {
        var inPt = transform.out2in(pt);
        return f.call(this, inPt, button, shift, ctrl);
    };
    html.setPointClick(this.canvas, g);
};
ScreenTransform.prototype.setWheel = function(f) {
    html.setWheel(this.canvas, f);
};

html.ScreenTransform = ScreenTransform;

function Tree(name) {
    if (name) {
        this.elem = document.getElementById(name);
        this.style = this.elem.style;
    }
}
Tree.prototype.spur = function(kind, text) {
    var newT = new Tree();
    newT.elem = html.spur(this.elem, kind, text);
    newT.style = newT.elem.style;
    newT.parent = this;
    return newT;
};
Tree.prototype.spurElement = function(elem) {
    this.elem.appendChild(elem);
    var newT = new Tree();
    newT.elem = elem;
    newT.style = elem.style;
    newT.parent = this;
    return newT;
};
Tree.prototype.addText = function(text) {
    var child = document.createTextNode(text);
    this.elem.appendChild(child);
    return this;
};

Tree.prototype.spurClass = function(kind, className, text) {
    var newT = new Tree();
    newT.elem = html.spur(this.elem, kind, text);
    newT.style = newT.elem.style;
    newT.elem.className = className;
    newT.parent = this;
    return newT;
};
Tree.prototype.and = function(kind, text) {
    return this.parent.spur(kind, text);
};
Tree.prototype.andClass = function(kind, className, text) {
    return this.parent.spurClass(kind, className, text);
};
Tree.prototype.clear = function() {
    while (this.elem.firstChild) {
        this.elem.removeChild(this.elem.firstChild);
    }
};
Tree.prototype.setText = function(text) {
    this.elem.textContent = text;
};
Tree.prototype.setClass = function(text) {
    this.elem.className = text;
};
Tree.prototype.setClick = function(f, noMenu) {
    if (noMenu) {
        this.elem.oncontextmenu= function() { return false; };
    }
    this.elem.addEventListener("mouseup", f, false);
};
Tree.prototype.setPointClick = function(f, noMenu) {
    if (noMenu) {
        this.elem.oncontextmenu= function() { return false; };
    }
    this.elem.addEventListener("mouseup", html.clickWrap(f), false);
};
Tree.prototype.setWheel = function(f) {
    html.setWheel(this.elem, f);
};

html.Tree = Tree;

})(muleObj);
