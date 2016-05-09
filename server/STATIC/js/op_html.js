(function(muleObj) {

if (!muleObj.overpower) {
    muleObj.overpower = {};
}
if (!muleObj.ajax) {
    muleObj.ajax = {};
}
if (!muleObj.html) {
    muleObj.html = {};
}


if (!muleObj.overpower.data) {
    muleObj.overpower.data = {};
}
if (!muleObj.overpower.html) {
    muleObj.overpower.html = {};
}

var geometry = muleObj.geometry;
var html = muleObj.html;
var overpower = muleObj.overpower;

var data = overpower.data;
var ophtml = overpower.html;
var net = overpower.net;

ophtml.icon = {blink: false};
ophtml.icon.elem = document.querySelector('link[rel="shortcut icon"]');
ophtml.icon.animate = function() {
    var nonBlinkSource = "/static/img/yd32.ico";
    if (overpower.data.game.newTurn) {
        ophtml.icon.blink = !ophtml.icon.blink;
        if (ophtml.icon.blink) {
            ophtml.icon.elem.href = "/static/img/yd32blink.ico";
        } else {
            ophtml.icon.elem.href = nonBlinkSource;
        }
        window.setTimeout(ophtml.icon.animate, 1000);
    } else {
        ophtml.icon.blink = false;
        ophtml.icon.elem.href = nonBlinkSource;
    }
};


ophtml.infobox = {};
var infobox = ophtml.infobox;
infobox.render = function() {
    infobox.overview.render();
    infobox.targets.render();
};

function setPlanet(dat) {
    var hex, name;
    if (dat.hex && dat.name) {
        hex = dat.hex;
        name = dat.name;
    } else if ((dat.x === 0 || dat.x) && (dat.y === 0 || dat.y)) {
        hex = dat;
        name = "("+dat.x+","+dat.y+")";
    } else {
        hex = new geometry.Hex();
        name = "ERROR";
    }
    this.textContent = name;
    this.hex = hex;
}

ophtml.jumperButton = function(dat, elem) {
    if (!elem) {
        elem = document.createElement("button");
    }
    elem.className = "target1";
    elem.oncontextmenu= function() { return false; };
    elem.setPlanet = setPlanet;
    if (dat) {
        elem.setPlanet(dat);
    }
    var jump = function(event) {
        overpower.commands.clickHex(elem.hex, event.button, event.shiftKey);
    };
    elem.addEventListener("mouseup", jump);
    return elem;
};

ophtml.target2Button = function(dat, elem) {
    var elem2 = ophtml.jumperButton(dat, elem);
    elem2.className = "target2";
    return elem2;
};

html.Tree.prototype.addDisplay = function(text, num, className) {
    this.addText(text+": [ ");
    if (className) {
        this.spurClass("span", className, num);
    } else {
        this.spur("span", num);
    }
    this.addText(" ] ");
};

})(muleObj);
