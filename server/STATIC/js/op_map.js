(function(muleObj) {

var html = muleObj.html;
var geometry = muleObj.geometry;

var map = {};
muleObj.overpower.map = map;


var canvas = document.getElementById('mainscreen');
map.screen = new html.ScreenTransform(canvas, 0,0, 0.25, 35, 0.55);
map.hexGrid = new muleObj.geometry.HexGrid(map.screen.transform);

map.hexPath = function(hex, path) {
    if (!path) {
        path = new Path2D();
    }
    var corners = map.hexGrid.cornerPts(hex);
    path.moveTo(corners[5].x, corners[5].y);
    corners.forEach(function(pt) {
        path.lineTo(pt.x, pt.y);
    });
    return path;
};
map.visibleHexes = function() {
    return map.hexGrid.hexesIn(0,0, map.screen.canvas.width, map.screen.canvas.height);
};
map.setFrame = function(width, height) {
    var spacer = document.getElementById('canvasspacer');
    spacer.style.width = width+'px';
    spacer.style.height = height+'px';
    var targetframe = document.getElementById('targetframe');
    targetframe.style.marginLeft = ''+(width + 10)+'px';
    map.screen.resize(width, height);
    var stars = muleObj.overpower.stars;
    stars.screen.resize(width, height);
    stars.moreStars();
};
map.scaleStep = function(zoomIn) {
    var curScale = map.screen.transform.scale;
    if ((curScale >= 40 && zoomIn) || (curScale <= 0.25 && !zoomIn)) {
        return false;
    }
    var dir = (zoomIn) ? 1 : -1;
    var delta;
    if (curScale > 20 || (zoomIn && curScale === 20)) {
        delta = 2*dir;
    } else if (curScale > 15 || (zoomIn && curScale == 15)) {
        delta = dir;
    } else if (curScale > 5 || (zoomIn && curScale == 5)) {
        delta = 0.5*dir;
    } else {
        delta = 0.25*dir;
    }
    var dScale;
    if (curScale + delta <= 0.25) {
        dScale = 0.25 - curScale;
    } else if (curScale+delta >= 40) {
        dScale = 40 - curScale;
    } else {
        dScale = delta;
    }
    map.screen.setScale(dScale);
    return true;
};

function mapClick(inPoint, button, shift) {
    var hex = inPoint.hexAt();
    console.log("MAPCLICK", inPoint, hex, button, shift);
}
map.screen.setInClick(mapClick);

function mapWheel(up, shift, ctrl) {
    console.log("MAPWHEEL", up, shift, ctrl);
    if (shift) {
        var theta = (up > 0) ? 1: -1;
        theta *= 0.015;
        map.screen.rotate(theta);
        map.render();
        muleObj.overpower.stars.screen.rotate(theta*0.5);
        return;
    }
    if (ctrl) {
        var delta = (up > 0) ? 1: -1;
        delta *= 10;
        map.setFrame(map.screen.canvas.width + delta, map.screen.canvas.height + delta);
        map.render();
        return;
    }
    if (map.scaleStep(up > 0)) {
        map.render();
    }
}

map.screen.setWheel(mapWheel);

map.render = function() {
    var canvas = map.screen.canvas;
    var ctx = canvas.getContext('2d');
    // CLEAR //
    ctx.clearRect(0,0,canvas.width, canvas.height);
    //
    var scale = map.screen.transform.scale;
    var showGrid = scale >= 15;
    var visHexes;
    if (scale > 5) {
        visHexes = map.visibleHexes();
    }
    var path;
    // GRID //
    if (showGrid) {
        if (scale < 25) {
            ctx.lineWidth = 0.5;
        } else if (scale < 35) {
            ctx.lineWidth = 1;
        } else {
            ctx.lineWidth = 1.5;
        }
        path = new Path2D();
        var drawHex = function(hex) {
            map.hexPath(hex, path);
        };
        visHexes.list.forEach(drawHex);
        //hexList.forEach(drawHex);
        ctx.strokeStyle = "#2f2f8f";
        ctx.stroke(path);
    }
};

map.render();

})(muleObj);
