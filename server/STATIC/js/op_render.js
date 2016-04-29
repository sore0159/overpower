(function(muleObj) {

if (!muleObj.geometry) {
    console.log("OP START FAILED: REQUIRES GEOMETRY");
    return;
}

if (!muleObj.overpower) {
    muleObj.overpower = {};
}
if (!muleObj.overpower.render) {
    muleObj.overpower.render = {};
}
if (!muleObj.overpower.data) {
    muleObj.overpower.data = {};
}

var data = muleObj.overpower.data;
var render = muleObj.overpower.render;

render.map = function() {
    var map = data.map;
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

})(muleObj);
