(function(muleObj) {

if (!muleObj.geometry || !muleObj.geometry.Hex || !muleObj.geometry.Point) {
    console.log("GRID CREATION FAILED: requires geometry.Hex and geometry.Point");
    return;
}

function HexGrid(canvas, transform) {
    this.canvas = canvas;
    this.transform = transform;
}
HERE

function visibleHexes(minX, minY, maxX, maxY) {
    var minObj = {};
    var maxObj = {};
    function xChecker(x, y) {
        var pt = new Point(x, y);
        var hex = pt.hexAt();
        var key = ""+hex.x;
        if (!minObj[key] || minObj[key].y > hex.y) {
            minObj[key] = hex;
        }
        if (!maxObj[key] || maxObj[key].y < hex.y) {
            maxObj[key] = hex;
        }
    }
    var y = 0;
    var x = minX;
    for (y = minY; y <=maxY; y+=1) {
        xChecker(x,y);
    }
    x = maxX;
    for (y = minY; y <=maxY; y+=1) {
        xChecker(x,y);
    }
    y = minY;
    for (x = minX; x <=maxX; x+=1) {
        xChecker(x,y);
    }
    y = maxY;
    for (x = minX; x <=maxX; x+=1) {
        xChecker(x,y);
    }
    var hexList = [];
    var hexMap = new Map();
    hexMap.hasHex = function(hex) {
        var hexSet = this.get(hex.x);
        if (!hexSet) {
            return false;
        }
        return hexSet.has(hex.y);
    };
    Object.keys(minObj).forEach(function(key) {
        var minHex = minObj[key];
        var maxHex = maxObj[key];
        var hexSet = hexMap.get(minHex.x);
        if (!hexSet) {
            hexSet = new Set();
            hexMap.set(minHex.x, hexSet);
        }
        for (i = minHex[1]; i<=maxHex[1];i++) {
            var hex = new Hex(minHex.x, i);
            hexList.push(hex);
            hexSet.add(i);
        }
    });
    hexMap.list = hexList;
    return hexMap;
}

})(muleObj);
