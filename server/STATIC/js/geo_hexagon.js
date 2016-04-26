/*
        ____
       /\  /\
  ____/__\/__\
 /\  /\  /\  /
/__\/__\/__\/
\  /\  / 
 \/__\/

*/
(function(muleObj) {

if (!muleObj.geometry) {
    muleObj.geometry = {};
}


function Hex(x, y) {
    this.x = x || 0;
    this.y = y || 0;
}
Hex.prototype.eq = function(hex) {
    return (hex.x === this.x) && (hex.y === this.y);
};
Hex.prototype.stepsTo = function(target) {
    var thisZ = - (this.x + this.y);
    var targetZ = - (target.x + target.y);
    var dx = Math.abs(this.x - target.x);
    var dy = Math.abs(this.y - target.y);
    var dz = Math.abs(thisZ - targetZ);
    return (dx + dy + dz)/2;
};

Hex.prototype.addHex = function(hex) {
    return this.add(hex.x, hex.y);
};
Hex.prototype.add = function(x, y) {
    var dx = x || 0;
    var dy = y || 0;
    var pt = new Hex(this.x + dx, this.y+dy);
    return pt;
};
Hex.prototype.scale = function(s) {
    if (!s) {
        return;
    }
    var pt = new Hex(this.x * s, this.y*s);
    return pt;
};

muleObj.geometry.Hex = Hex;

var Point = muleObj.geometry.Point;
if (!Point) {
    console.log("NO POINT FUNCTION FOUND, NOT INCLUDING HEXAGON/POINT FEATURES");
    return;
}

Hex.prototype.centerPt = function() {
    var px = 1.5 * this.x;
    // 1^2 = (.5)^2 + (dy^2)
    // dy = sqrt( .75 ) =  .86602540378
    var py = (2*this.y + this.x)*0.86602540378;
    var pt = new Point(px, py);
    return pt;
};

Hex.prototype.cornerPts = function() {
    var hexTriH = 0.86602540378;
    var center = this.centerPt();
    var pts = [
        center.add(1, 0),
        center.add(0.5, hexTriH),
        center.add(-0.5, hexTriH),
        center.add(-1, 0),
        center.add(-0.5, -hexTriH),
        center.add(0.5, -hexTriH),
    ];
    return pts;
};

function HexGrid(transform) {
    this.transform = transform;
}

HexGrid.prototype.hexAt = function(pt) {
    var inPt = this.transform.out2in(pt);
	//       __
	//      | \|  width =  1.5
	//box = |_/|  height = sqrt(3) = 2*hexTriH;
    //
    var hexTriH = 0.86602540378;
    var x = inPt.x + 0.5;
    var y = inPt.y + hexTriH;
    var hx = Math.floor(x / 1.5);
    //var hy = Math.floor((y-hx*hexTriH)/ (2*hexTriH));
    var hy = Math.floor(0.5*((y/hexTriH)-hx));
    var boxX = x - 1.5*hx;
    if (boxX < 1) {
        return [hx, hy];
    }
    boxX -= 1;
    var boxY = y - (2*hy+hx)*hexTriH;
    var slope = 2*hexTriH;
    if (boxY > hexTriH) {
        boxY -= 2*hexTriH;
        if (boxX*-slope < boxY) {
            return [hx+1, hy];
        }
        return [hx, hy];
    }
    if (boxX*slope > boxY) {
        return [hx+1, hy-1];
    }
    var hex = new Hex(hx, hy);
    return hex;
};

HexGrid.prototype.centerPt = function(hex) {
    var inPt = hex.centerPt();
    return this.transform.in2out(inPt);
};
HexGrid.prototype.cornerPts = function(hex) {
    var inPts = hex.cornerPts();
    var outPts = [];
    for (var i = 0; i < inPts.length; i++) {
        outPts[i] = this.transform.in2out(inPts[i]);
    }
    return outPts;
};

HexGrid.prototype.hexesIn = function(minX, minY, maxX, maxY) {
    var minObj = {};
    var maxObj = {};
    var grid = this;
    function xChecker(x, y) {
        var pt = new Point(x, y);
        var hex = grid.hexAt(pt);
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
};

muleObj.geometry.HexGrid = HexGrid;

})(muleObj);
