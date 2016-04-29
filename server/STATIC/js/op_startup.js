(function(muleObj) {

var html = muleObj.html;
var geometry = muleObj.geometry;

var data = muleObj.overpower.data;
var render = muleObj.overpower.render;
var parse = muleObj.overpower.parse;


(function() {
    data.map = {};
    var map  = data.map;
    var canvas = document.getElementById('mainscreen');
    map.screen = new geometry.ScreenTransform(canvas, 0,0, 35, 0.25, 0.55);
    map.screen.hexesIn = function() {
        return geometry.hexesIn(0,0, map.screen.canvas.width, map.screen.canvas.height, map.screen.transform);
    };
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

    function mapClick(point, button, shift) {
        var hex = data.map.hexGrid.hexAt(point);
        console.log("MAPCLICK", point, hex, button, shift);
    }

    function mapWheel(up, shift, ctrl) {
        console.log("MAPWHEEL", up, shift, ctrl);
        if (shift) {
            parse.rotateStars(up > 0, 0.5);
            muleObj.overpower.render.stars();
            parse.rotateMap(up > 0, 1.0);
            muleObj.overpower.render.map();
            return;
        }
    }

    render.map();
})();

(function() {

})();


})(muleObj);
