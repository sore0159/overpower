(function(muleObj) {

var html = muleObj.html;
var geometry = muleObj.geometry;
var overpower = muleObj.overpower;

if (!overpower.map) {
    overpower.map = {};
}

var map = overpower.map;

var canvas = document.getElementById('mainscreen');
map.screen = new html.ScreenTransform(canvas, 0,0, 0.25, 15, 0.55);
map.hexGrid = new geometry.HexGrid(map.screen.transform);
map.getScale = function() {
    return this.screen.transform.scale;
};

map.setFrame = function(width, height) {
    var spacer = document.getElementById('canvasspacer');
    spacer.style.width = width+'px';
    spacer.style.height = height+'px';
    var targetframe = document.getElementById('targetframe');
    targetframe.style.marginLeft = ''+(width + 10)+'px';
    map.screen.resize(width, height);
    var stars = overpower.stars;
    stars.screen.resize(width, height);
    stars.moreStars();
};
map.scaleStep = function(zoomIn) {
    var curScale = map.getScale();
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
map.snapTo = function(hex) {
    map.center = hex;
    map.screen.setCenter(hex.centerPt());
};
map.distToCenter = function() {
    if (!map.center) {
        return 0;
    }
    var centerAt = map.hexGrid.centerPt(map.center);
    return centerAt.dist(map.screen.center());
};
map.moveTowardCenter = function(distLeft) {
    var endpt;
    if (distLeft <= 10) {
        endpt = map.center.centerPt();
    } else {
        var speed;
        if (distLeft <= 20) {
            speed = 10;
        } else {
            //speed = 10 + Math.pow((distLeft-20), 3);
            speed = 10 + Math.pow((distLeft - 20)/25, 1.25);
        }
        var startPt = map.screen.center();
        var towardPt = map.hexGrid.centerPt(map.center);
        var polarDat = startPt.polarTo(towardPt);
        endpt = startPt.addPolar(speed, polarDat[1]);
        endpt = map.screen.out2in(endpt);
    }
    map.screen.setCenter(endpt);
};

function mapClick(inPoint, button, shift) {
    var hex = inPoint.hexAt();
    console.log("MAPCLICK", inPoint, hex, button, shift);
    if (button === 0) { 
        if (shift) {
            overpower.map.center = hex;
        } else {
            overpower.data.targets.setT1(hex);
            map.redraw = true;
        }
    } else if (button === 2) {
        overpower.data.targets.setT2(hex, shift);
        map.redraw = true;
    }
}
map.screen.setInClick(mapClick);

function mapWheel(up, shift, ctrl) {
    console.log("MAPWHEEL", up, shift, ctrl);
    if (shift) {
        var theta = (up > 0) ? 1: -1;
        theta *= 0.015;
        map.screen.rotate(theta);
        map.redraw = true;
        overpower.stars.screen.rotate(theta*0.5);
        return;
    }
    if (ctrl) {
        var delta = (up > 0) ? 1: -1;
        delta *= 10;
        map.setFrame(map.screen.canvas.width + delta, map.screen.canvas.height + delta);
        //map.render();
        map.redraw = true;
        return;
    }
    if (map.scaleStep(up > 0)) {
        map.redraw = true;
    }
}

map.screen.setWheel(mapWheel);

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
    //return map.hexGrid.hexesIn(map.screen.canvas.width * 0.25, map.screen.canvas.height * 0.25 , map.screen.canvas.width * 0.75, map.screen.canvas.height * 0.75);
    return map.hexGrid.hexesIn(0,0, map.screen.canvas.width, map.screen.canvas.height);
};

map.renderGrid = function(ctx, list) {
    map.setLineWidth(ctx);
    path = new Path2D();
    var drawHex = function(hex) {
        var corners = map.hexGrid.cornerPts(hex);
        path.moveTo(corners[0].x, corners[0].y);
        path.lineTo(corners[1].x, corners[1].y);
        path.lineTo(corners[2].x, corners[2].y);
        path.lineTo(corners[3].x, corners[3].y);
    };
    list.forEach(drawHex);
    ctx.strokeStyle = "#2f2f8f";
    ctx.stroke(path);
};

map.renderTargets = function(ctx, visHexes, targets) {
    ctx.lineWidth += 1;
    if (targets.secondary && (!targets.primary || !targets.primary.hex.eq(targets.secondary.hex)) && visHexes.hasHex(targets.secondary.hex)) {
        ctx.strokeStyle = "#999900";
        ctx.stroke(map.hexPath(targets.secondary.hex));
    }
    if (targets.primary && visHexes.hasHex(targets.primary.hex)) {
        ctx.strokeStyle = "#ffff00";
        ctx.stroke(map.hexPath(targets.primary.hex));
    }
    ctx.lineWidth -= 1;
};

map.renderPlanet = function(ctx, pv, planetRad, showGrid) {
    var scale = map.getScale();
    if (!showGrid && overpower.data.targets.isT1(pv.hex)) {
        ctx.fillStyle = "#ffff00";
    } else if (!showGrid && overpower.data.targets.isT2(pv.hex)) {
        ctx.fillStyle = "#999900";
    } else if (pv.primaryfaction === overpower.FID) {
        if (pv.secondaryfaction) {
            ctx.fillStyle = "#af00af";  // contested
        } else {
            ctx.fillStyle = "#0fff0f";
        }
    } else if (pv.secondaryfaction === overpower.FID) {
        ctx.fillStyle = "#af00af";  // contested
    } else if (pv.primaryfaction) {
        ctx.fillStyle = "#ff0f0f";
    //} else if (pv.hex.x === 0) {
    } else {
        ctx.fillStyle = "#ffffff";
    }
    var path = new Path2D();
    var center = map.hexGrid.centerPt(pv.hex);
    if (showGrid) {
        center.y -= scale * 0.25;
    }
    path.moveTo(center.x+planetRad, center.y);
    path.arc(center.x, center.y, planetRad, 0, 2*Math.PI);
    ctx.fill(path);
    ctx.stroke(path);
    if (scale < 1.5) {
        return;
    }
    var nameStr;
    // TOP //
    if (scale < 3.75) {
        if (pv.avail) {
            nameStr = "("+pv.avail+")";
         // ctx.strokeText(nameStr, center.x+planetRad*0.25, center.y-1.15*planetRad);
            ctx.fillText(nameStr, center.x+planetRad*0.25, center.y-1.15*planetRad);
        }
    } else {
        if (pv.avail) {
            nameStr = "("+pv.avail+")"+pv.name;
        } else {
            nameStr = pv.name;
        }
        //ctx.strokeText(nameStr, center.x+planetRad*0.25, center.y-1.15*planetRad);
        ctx.fillText(nameStr, center.x+planetRad*0.25, center.y-1.15*planetRad);
    }
    // BOTTOM //
};

map.setLineWidth = function(ctx) {
   var scale = map.getScale();
    if (scale < 25) {
        ctx.lineWidth = 0.5;
    } else if (scale < 35) {
        ctx.lineWidth = 1;
    } else {
        ctx.lineWidth = 1.5;
    }
};
map.isShowGrid = function() {
    return map.getScale() >= 15;
};

//map.render = function() {
function renderMap(timestamp) {
    var data = overpower.data;
    var canvas = map.screen.canvas;
    var ctx = canvas.getContext('2d');
    // CLEAR //
    ctx.clearRect(0,0,canvas.width, canvas.height);
    //
    var scale = map.getScale();
    var showGrid = map.isShowGrid();
    var visHexes;
    if (scale > 5) {
        visHexes = map.visibleHexes();
    }
    var path;
    // GRID //
    if (showGrid) {
        map.renderGrid(ctx, visHexes.list);
        if (data.targets) {
            map.renderTargets(ctx, visHexes, data.targets);
        }
    }
    var planetRad;
    if (scale > 10) {
        planetRad = scale * 0.45;
    } else if (scale > 4) {
        planetRad = 5;
    }  else if (scale > 2) {
        planetRad = 4;
    } else {
        planetRad = 3;
    }
    var fontHeight;
    if (scale >= 20) {
        ctx.font = "12pt Courier New";
        fontHeight = 14;
    } else if (scale >= 15) {
        ctx.font = "11pt Courier New";
        fontHeight = 13;
    } else if (scale >= 10) {
        ctx.font = "10pt Courier New";
        fontHeight = 12;
    } else if (scale >= 7.5) {
    //} else if (scale >= 5) {
        ctx.font = "9pt Courier New";
        fontHeight = 11;
    } else { 
        ctx.font = "8pt Courier New";
        fontHeight = 10;
    }
    // ORDERS/TRAILS/DESTINATIONS  //
    // SHIPS //
    // PLANETS //
    if (data.planetGrid) {
        map.setLineWidth(ctx);
        ctx.strokeStyle = "#000000";
        var later = {};
        data.planetGrid.forEach(function(pv, hex) {
            if (!visHexes || visHexes.hasHex(hex)) {
                if (data.targets.isT1(hex)) {
                    later.t1 = pv;
                } else if (data.targets.isT2(hex)) {
                    later.t2 = pv;
                } else {
                    map.renderPlanet(ctx, pv, planetRad, showGrid);
                }
            }
        });
        if (later.t2) {
            map.renderPlanet(ctx, later.t2, planetRad, showGrid);
        }
        if (later.t1) {
            map.renderPlanet(ctx, later.t1, planetRad, showGrid);
        }
    }
}

map.animate = function(timestamp) {
    var redraw = map.redraw;
    if (redraw !== false) {
        redraw = true;
    }
    map.redraw = false;
    var dist = map.distToCenter();
    if (dist > 1) {
        map.moveTowardCenter(dist);
        redraw = true;
    }
    if (redraw) {
        renderMap(timestamp);
    }
    window.requestAnimationFrame(map.animate);
};

map.animate();

})(muleObj);
