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
    return (hex) && (hex.x === this.x) && (hex.y === this.y);
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
Hex.prototype.ring = function(r) {
    var hex;
    if (r < 1) {
        hex = new Hex(this.x, this.y);
        return [hex];
    }
    hex = this.add(r, 0);
    var list = [];
    for (var i = 0; i<r; i+=1) {
        hex = hex.add(-1, 1);
        list.push(hex);
    }
    for (i = 0; i<r; i+=1) {
        hex = hex.add(-1, 0);
        list.push(hex);
    }
    for (i = 0; i<r; i+=1) {
        hex = hex.add(0, -1);
        list.push(hex);
    }
    for (i = 0; i<r; i+=1) {
        hex = hex.add(1, -1);
        list.push(hex);
    }
    for (i = 0; i<r; i+=1) {
        hex = hex.add(1, 0);
        list.push(hex);
    }
    for (i = 0; i<r; i+=1) {
        hex = hex.add(0, 1);
        list.push(hex);
    }
    return list;
};

muleObj.geometry.Hex = Hex;

var Point = muleObj.geometry.Point;
if (!Point) {
    console.log("NO POINT FUNCTION FOUND, NOT INCLUDING HEXAGON/POINT FEATURES");
    return;
}
Point.prototype.hexAt = function() {
	//       __
	//      | \|  width =  1.5
	//box = |_/|  height = sqrt(3) = 2*hexTriH;
    //
    var hexTriH = 0.86602540378;
    var x = this.x + 0.5;
    var y = this.y + hexTriH;
    var hx = Math.floor(x / 1.5);
    //var hy = Math.floor((y-hx*hexTriH)/ (2*hexTriH));
    var hy = Math.floor(0.5*((y/hexTriH)-hx));
    var boxX = x - 1.5*hx;
    var hex = new Hex(hx, hy);
    if (boxX < 1) {
        return hex;
    }
    boxX -= 1;
    var boxY = y - (2*hy+hx)*hexTriH;
    var slope = 2*hexTriH;
    if (boxY > hexTriH) {
        boxY -= 2*hexTriH;
        if (boxX*-slope < boxY) {
            return hex.add(1);
        }
        return hex;
    }
    if (boxX*slope > boxY) {
        return hex.add(1,-1);
    }
    return hex;
};

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

function HexMap() {
    this.xMap = new Map();
}
HexMap.prototype.setHex = function(hex, obj) {
    var yMap = this.xMap.get(hex.x);
    if (!yMap) {
        yMap = new Map();
        this.xMap.set(hex.x, yMap);
    }
    yMap.set(hex.y, obj);
};
HexMap.prototype.getHex = function(hex) {
    var yMap = this.xMap.get(hex.x);
    if (!yMap) {
        return;
    }
    return yMap.get(hex.y);
};
HexMap.prototype.clear = function() {
    this.xMap.clear();
};
HexMap.prototype.deleteHex = function(hex) {
    var yMap = this.xMap.get(hex.x);
    if (!yMap) {
        return false;
    }
    return yMap.delete(hex.y);
};
HexMap.prototype.forEach = function(callback, thisObj) {
    var map = this;
    this.xMap.forEach(function(yMap, x) {
        yMap.forEach(function(obj, y) {
            var hex = new Hex(x,y);
            callback.call(thisObj, obj, hex, map);
        });
    });
};

muleObj.geometry.HexMap = HexMap;

function HexGrid(transform) {
    this.transform = transform;
}

HexGrid.prototype.hexAt = function(pt) {
    var inPt = this.transform.out2in(pt);
    return inPt.hexAt();

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
    var transform = this.transform;
    var f = function(pt) {
        return transform.out2in(pt);
    };
    return hexesIn(minX, minY, maxX, maxY, f);
};

function hexesIn(minX, minY, maxX, maxY, transform) {
    var xMap = new Map();
    var grid = this;
    function xChecker(x, y) {
        var pt = new Point(x, y);
        if (transform) {
            pt = transform(pt);
        }
        var hex = pt.hexAt();
        var yDat = xMap.get(hex.x);
        if (!yDat) {
            xMap.set(hex.x, {x: hex.x,  min: hex.y, max: hex.y });
            return;
        }
        if (yDat.min > hex.y) {
            yDat.min = hex.y;
            return;
        }
        if (yDat.max < hex.y) {
            yDat.max = hex.y;
            return;
        }
    }
    var y = minY;
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
    xMap.forEach(function(yDat) {
        var hexSet = hexMap.get(yDat.x);
        if (!hexSet) {
            hexSet = new Set();
            hexMap.set(yDat.x, hexSet);
        }
        for (i = yDat.min ; i<=yDat.max ; i++) {
            var newHex = new Hex(yDat.x, i);
            hexList.push(newHex);
            hexSet.add(i);
        }
    });
    hexMap.list = hexList;
    return hexMap;
}

muleObj.geometry.HexGrid = HexGrid;
muleObj.geometry.hexesIn = hexesIn;

})(muleObj);
