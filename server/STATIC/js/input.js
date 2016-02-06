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
    this.muleWheeled(up, event.shiftKey, event.ctrlKey);
    return false;
}
canvas.onmousedown = mapClick;
canvas.oncontextmenu= function() { return false; };
canvas.onmousewheel = mapWheel;
canvas.onDOMMouseScroll = mapWheel;
canvas.addEventListener("DOMMouseScroll", mapWheel);

// ----------- GAME COMMANDS ------------- //

canvas.muleClicked = function(pt, button, shift) {
    var grid = this.muleGrid;
    var hex = grid.p2Hex(pt);
    var help = !this.gridShowing();
    if (button === 0 && !shift) {
        this.overpowerData.setTargetOne(hex, help);
    } else if (button === 2) {
        this.overpowerData.setTargetTwo(hex, help, shift);
    } else if (button === 1 || button === 0) {
        this.setCenterDest(hex);
    }
    this.drawMap();
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
            this.refreshTargetOrderInfo();
            this.drawMap();
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


/*
canvas.muleWheeled = function(up, shiftKey) {
    console.log("wheel", up, shiftKey);
};
*/



})();
