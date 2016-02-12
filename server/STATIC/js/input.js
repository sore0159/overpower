(function() {


var canvas = document.getElementById('mainscreen');
var boxConfirm = document.getElementById('orderconfirm');
var boxTurn = document.getElementById('turnbox');
var boxTurnButton = document.getElementById('turnchange');
var swapTargetsButton = document.getElementById('swapbutton');
var blockButton = document.getElementById('blockbutton');

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
    this.muleWheeled(up, event.shiftKey, event.ctrlKey);
    return false;
}
canvas.addEventListener("mouseup", mapClick);
canvas.oncontextmenu= function() { return false; };
canvas.onmousewheel = mapWheel;
canvas.onDOMMouseScroll = mapWheel;
canvas.addEventListener("DOMMouseScroll", mapWheel);

boxConfirm.onclick = confirmOrder;
boxTurnButton.onclick = turnChange;
blockButton.onclick = function() {
    canvas.refetchFullView();
};
swapTargetsButton.onclick = function() {
    canvas.overpowerData.swapTargets();
    canvas.redrawPage();
};

// ----------- GAME COMMANDS ------------- //

function putJSON(url, obj, onFail, onPass) {
    var req = new XMLHttpRequest();
    function handlePUT() {
        var err;
        if (req.status === 200) {
            var resp = JSON.parse(req.responseText);
            if (resp.status === "fail") {
                err = "json fail response from server:"+'\n'+resp.data;
            } else if (resp.status === "error") {
                err = "json error response from server:"+'\n'+resp.message;
            } else if (resp.status !== "success") {
                err = "unknown json response from server:"+resp.status;
            }
        } else {
            err = req.status +'\n'+req.statusText+"\nJSON:"+
                req.responseText;
        }
        if (err) {
            onFail(err);
        } else if (onPass) {
            onPass();
        }
    }
    req.addEventListener("load", handlePUT);
    req.addEventListener("error", function() {
        onFail("an error occured with the request");
    });
    req.open("PUT", url, true);
    req.setRequestHeader('Content-Type', 'application/json');
    req.send(JSON.stringify(obj));
}

function putErr(err) {
    console.log("OVERPOWER SERVER ERROR:", err);
    canvas.blockScreen("There was a server error!", "Click here to reload");
}

function turnChange() {
    var faction = canvas.overpowerData.faction;
    if (!faction.donebuffer) {
        faction.donebuffer = 1;
        faction.done = true;
    } else {
        faction.donebuffer = 0;
        faction.done = false;
    }
    canvas.setupText();
    console.log("SENDING:", faction);
    canvas.blockScreen("Checking for new turn...");
    putJSON("/overpower/json/factions", faction, putErr, function() {
        canvas.turnCheck();
    });
}

function confirmOrder() {
    var order = canvas.overpowerData.targetOrder;
    if (order && order.brandnew && order.size > 0) {
        order.brandnew = null;
        order.originSize = order.size;
        canvas.overpowerData.orders.push(order);
        canvas.overpowerData.targetOneInfo.orders.push(order);
        putJSON("/overpower/json/orders", order, putErr);
        canvas.redrawPage();
    } else if (order && !order.brandnew && order.originSize != order.size) {
        order.oldSize = order.originSize;
        order.originSize = order.size;
        if (order.size < 1) {
            canvas.overpowerData.orders = canvas.overpowerData.orders.filter(function(testOrder) {
                return testOrder !== order;
            });
        }
        canvas.overpowerData.targetOrder = null;
        putJSON("/overpower/json/orders", order, putErr);
        canvas.overpowerData.setTargetOne(canvas.overpowerData.targetOne);
    } else {
        console.log("ERROR: CONFIRATION FOR BAD ORDER", order);
    }
    canvas.redrawPage();
}

canvas.muleClicked = function(pt, button, shift) {
    var grid = this.muleGrid;
    var hex = grid.p2Hex(pt);
    var help = !this.gridShowing();
    if (button === 0 && !shift) {
        this.overpowerData.setTargetOne(hex, help);
        this.redrawPage();
    } else if (button === 2) {
        this.overpowerData.setTargetTwo(hex, help, shift);
        this.redrawPage();
    } else if (button === 1 || button === 0) {
        this.setCenterDest(hex);
    }
};

canvas.muleWheeled = function(up, shift, control) {
    if (control) {
        var theta;
        if (up > 0) {
            dTheta = 0.01;
        } else {
            dTheta = -0.01;
        }
        this.muleGrid.rotateAround(dTheta, [this.width/2, this.height/2]);
        this.drawMap();
        return;
    }
    if (shift) {
        var targetOrder = this.overpowerData.targetOrder;
        if (targetOrder) {
            var curSize = targetOrder.size;
            if ((up > 0) && targetOrder.sourcePl.avail) {
                targetOrder.sourcePl.avail -= 1;
                targetOrder.size += 1;
            } else if ((up < 0) && targetOrder.size) {
                targetOrder.sourcePl.avail += 1;
                targetOrder.size -= 1;
            } else {
                return;
            }
            if (!targetOrder.originSize && targetOrder.originSize !== 0) {
                targetOrder.originSize = curSize;
            }
            this.redrawPage();
        }
        return;
    }
    var cur = this.muleGrid.scale;
    var dir;
    if (up > 0) {
        dir = 1;
    } else {
        dir = -1;
    }
    var delta;
    if (cur > 20 || (up > 0 && cur === 20)) {
        delta = 2*dir;
    } else if (cur > 15 || (up > 0 && cur == 15)) {
        delta = dir;
    } else if (cur > 5 || (up > 0 && cur == 5)) {
        delta = 0.5*dir;
    } else {
        delta = 0.25*dir;
    }
    var dScale;
    if (cur + delta <= 0) {
        dScale = 0.25-cur;
    } else if (cur+delta >= 40) {
        dScale = 40-cur;
    } else {
        dScale = delta;
    }
    this.muleGrid.scaleAround(dScale, [this.width/2, this.height/2]);
    this.drawMap();
};


})();
